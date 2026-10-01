"""Exercise the real B42 mod-selection path with a simulated installed-item list; no game loop."""
import argparse
import hashlib
import json
import os
from pathlib import Path
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[1]
GAME_SHA = "e1a69eb743ede60b213a0fe7f8b83d4fcab773036d256cc4543a336f3b058a33"


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--game", type=Path, required=True)
    parser.add_argument("--jdk", type=Path, required=True)
    args = parser.parse_args()
    game = args.game.resolve()
    gamejar = game / "projectzomboid.jar"
    assert hashlib.sha256(gamejar.read_bytes()).hexdigest() == GAME_SHA
    source = ROOT / "tools/InstallerModSelectionSmoke.java"
    report = ROOT / "artifacts/installer-engine"
    report.mkdir(parents=True, exist_ok=True)
    checks = []
    with tempfile.TemporaryDirectory(prefix="selection-", dir=report) as tmp:
        work = Path(tmp)
        classes = work / "classes"
        classes.mkdir()
        subprocess.run([str(args.jdk.resolve() / "bin/javac.exe"), "-encoding", "UTF-8", "-proc:none", "-cp", str(gamejar),
                        "-d", str(classes), str(source)], check=True, timeout=60)
        for case in ("default", "community"):
            cwd = work / case
            cwd.mkdir()
            result = subprocess.run([str(game / "jre64/bin/java.exe"), "-Duser.home=" + str(cwd), "-cp",
                                     os.pathsep.join(map(str, (classes, gamejar))), "InstallerModSelectionSmoke", case],
                                    cwd=cwd, text=True, capture_output=True, encoding="utf-8", errors="replace", timeout=45)
            log = result.stdout + result.stderr
            (report / (case + ".log")).write_text(log, encoding="utf-8")
            assert result.returncode == 0 and "PASS " + case + ":" in log, log
            checks.append(dict(case=case, passed=True))
    assert hashlib.sha256(gamejar.read_bytes()).hexdigest() == GAME_SHA
    data = dict(gameSha256=GAME_SHA, probeSourceSha256=hashlib.sha256(source.read_bytes()).hexdigest(),
                checks=checks, realEngineFilesystem=True, simulatedSteamFolderList=True, gameLaunched=False,
                multiplayerAccepted=False, limits="Synthetic metadata; no native Steam client, mod UI or gameplay acceptance")
    (report / "results.json").write_text(json.dumps(data, indent=2) + "\n", encoding="utf-8")
    print(json.dumps(data, indent=2))


if __name__ == "__main__":
    main()
