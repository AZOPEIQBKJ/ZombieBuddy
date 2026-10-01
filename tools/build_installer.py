"""Build the offline Windows installer from the frozen runtime, without installation."""
import argparse
import hashlib
import json
import os
from pathlib import Path
import shutil
import subprocess
import zipfile

ROOT = Path(__file__).resolve().parents[1]
ARCHIVE_SHA = "f1b66cf9e11aac8663de92de382944a979cb4f09e651d320d8d3938b471a3c99"
VERSION = "0.1.0-preview.2"


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--go", type=Path, required=True)
    parser.add_argument("--runtime", type=Path, required=True)
    parser.add_argument("--output", type=Path, default=ROOT / "artifacts" / "installer")
    parser.add_argument("--prepare-only", action="store_true", help="Stage the verified embedded payload for development tests")
    args = parser.parse_args()
    archive = args.runtime.read_bytes()
    assert hashlib.sha256(archive).hexdigest() == ARCHIVE_SHA, "Unexpected runtime archive"
    module = ROOT / "community-installer"
    (module / "payload.zip").write_bytes(archive)
    if args.prepare_only:
        print("Verified immutable runtime staged for installer tests.")
        return
    if subprocess.check_output(["git", "status", "--porcelain"], cwd=ROOT, text=True).strip():
        raise SystemExit("Commit reviewed source before producing the distributable installer")
    commit = subprocess.check_output(["git", "rev-parse", "HEAD"], cwd=ROOT, text=True).strip()
    selection = json.loads((ROOT / "artifacts/installer-engine/results.json").read_text())
    assert selection["probeSourceSha256"] == hashlib.sha256((ROOT / "tools/InstallerModSelectionSmoke.java").read_bytes()).hexdigest()
    assert all(check["passed"] for check in selection["checks"]) and len(selection["checks"]) == 2
    env = dict(os.environ, GOTOOLCHAIN="local", GOPROXY="off", GOOS="windows", GOARCH="amd64", CGO_ENABLED="0",
               GOCACHE=str(ROOT / "artifacts" / "installer-go-cache"))
    tests = subprocess.run([str(args.go.resolve()), "test", "-json", "./..."], cwd=module, env=env, capture_output=True, text=True)
    if tests.returncode:
        raise SystemExit(tests.stdout + tests.stderr)
    test_events = [json.loads(line) for line in tests.stdout.splitlines() if line.startswith("{")]
    assert any(event.get("Action") == "pass" and "Test" not in event for event in test_events)
    out = args.output.resolve() / VERSION
    out.mkdir(parents=True, exist_ok=True)
    executable = out / "ZombieBuddyCommunityInstaller.exe"
    if executable.exists():
        raise SystemExit("Installer version already packaged; do not overwrite an immutable build")
    subprocess.run([str(args.go.resolve()), "build", "-trimpath", "-ldflags", f"-s -w -X main.sourceCommit={commit}", "-o", str(executable), "."], cwd=module, env=env, check=True)
    for source, name in [(ROOT / "LICENSE.txt", "LICENSE.txt"), (ROOT / "doc/CommunityInstaller.md", "INSTALLER.md")]:
        shutil.copyfile(source, out / name)
    metadata = dict(installerVersion=VERSION, runtimeVersion="2.3.3-community.4", sourceCommit=commit,
                    embeddedRuntimeSha256=ARCHIVE_SHA, executableSha256=hashlib.sha256(executable.read_bytes()).hexdigest(),
                    platform="windows-amd64", goVersion=subprocess.check_output([str(args.go.resolve()), "version"], text=True).strip(),
                    localTestsPassed=True, gameplayAccepted=False, published=False, authenticodeSigned=False)
    metadata["skippedTests"] = [event["Test"] for event in test_events if event.get("Action") == "skip" and "Test" in event]
    (out / "installer-tests.jsonl").write_text(tests.stdout, encoding="utf-8")
    shutil.copyfile(ROOT / "artifacts/installer-engine/results.json", out / "engine-selection-check.json")
    (out / "manifest.json").write_text(json.dumps(metadata, indent=2) + "\n", encoding="utf-8")
    archive_path = out / f"ZombieBuddyCommunityInstaller-{VERSION}-windows-x64.zip"
    with zipfile.ZipFile(archive_path, "x", zipfile.ZIP_DEFLATED) as package:
        for name in [executable.name, "LICENSE.txt", "INSTALLER.md", "manifest.json", "engine-selection-check.json", "installer-tests.jsonl"]:
            package.write(out / name, name)
    hashes = [f"{hashlib.sha256(path.read_bytes()).hexdigest()}  {path.name}" for path in (executable, archive_path)]
    (out / "SHA256SUMS.txt").write_text("\n".join(hashes) + "\n", encoding="ascii")
    print(json.dumps(dict(output=str(out), **metadata), indent=2))


if __name__ == "__main__":
    main()
