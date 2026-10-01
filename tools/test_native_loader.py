"""Check the Windows DLL dependency gap without creating a JVM or starting PZ."""
import argparse
import ctypes
import hashlib
import json
from pathlib import Path
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
NATIVE_SHA = "c2ae9335e717ee24b2f4a40d1a3bf77f1519762a72a0459e766a2bbafc077f6c"


def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def child(game, mode):
    binary = game / "jre64/bin"
    api = ctypes.WinDLL("kernel32", use_last_error=True)
    api.LoadLibraryW.argtypes = [ctypes.c_wchar_p]
    api.LoadLibraryW.restype = ctypes.c_void_p
    api.SetDllDirectoryW.argtypes = [ctypes.c_wchar_p]
    api.SetDllDirectoryW.restype = ctypes.c_int
    api.GetProcAddress.argtypes = [ctypes.c_void_p, ctypes.c_char_p]
    api.GetProcAddress.restype = ctypes.c_void_p
    jvm = api.LoadLibraryW(str(binary / "server/jvm.dll"))
    assert jvm, ctypes.get_last_error()
    if mode == "native-loader-search-path":
        assert api.SetDllDirectoryW(str(binary))
    instrument = api.LoadLibraryW(str(binary / "instrument.dll"))
    error = ctypes.get_last_error() if not instrument else 0
    if mode == "native-loader-search-path":
        assert api.SetDllDirectoryW(None)
    print(json.dumps(dict(mode=mode, jvmLoaded=bool(jvm), instrumentLoaded=bool(instrument),
                         windowsError=error, agentOnLoadExport=bool(instrument and api.GetProcAddress(instrument, b"Agent_OnLoad")),
                         vmCreated=False, gameLaunched=False)))


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--game", type=Path, required=True)
    parser.add_argument("--mode", choices=["standard-search-path", "native-loader-search-path"])
    parser.add_argument("--output", type=Path, default=ROOT / "artifacts/installer-native/results.json")
    args = parser.parse_args()
    game = args.game.resolve()
    if args.mode:
        child(game, args.mode)
        return
    native = game / "zbNative.dll"
    assert digest(native) == NATIVE_SHA, "Unaudited native bootstrap"
    assert "java.instrument" in (game / "jre64/release").read_text()
    source = ROOT / "c/windows/zbNative.c"
    assert 'SetDllDirectoryA(".\\\\jre64\\\\bin")' in source.read_text()
    inputs = [native, game / "jre64/bin/instrument.dll", game / "jre64/bin/server/jvm.dll"]
    before = {str(path): digest(path) for path in inputs}
    results = []
    for mode in ("standard-search-path", "native-loader-search-path"):
        p = subprocess.run([sys.executable, __file__, "--game", str(game), "--mode", mode],
                           cwd=game, capture_output=True, text=True, timeout=15)
        assert p.returncode == 0, p.stderr
        results.append(json.loads(p.stdout))
    assert not results[0]["instrumentLoaded"] and results[0]["windowsError"] == 126, "Expected audited baseline gap not reproduced"
    assert results[1]["instrumentLoaded"] and results[1]["agentOnLoadExport"], "Native search-path fix failed"
    assert before == {str(path): digest(path) for path in inputs}, "Installed inputs changed"
    report = dict(passed=True, nativeLoaderSha256=NATIVE_SHA, nativeSourceSha256=digest(source),
                  probeSourceSha256=digest(Path(__file__)), checks=results, inputHashes=before,
                  scope="LoadLibraryW/SetDllDirectoryW only; no JVM creation, game loop, installation or gameplay acceptance")
    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text(json.dumps(report, indent=2) + "\n", encoding="utf-8")
    print(json.dumps(report, indent=2))


if __name__ == "__main__":
    main()
