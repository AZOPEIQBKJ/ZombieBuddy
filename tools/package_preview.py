"""Assemble a signed local preview without installing it or publishing anything."""
import argparse
import hashlib
import json
from pathlib import Path
import subprocess
import xml.etree.ElementTree as ET
from zipfile import ZipFile, ZipInfo, ZIP_DEFLATED

ROOT = Path(__file__).resolve().parents[1]


def sha(data):
    return hashlib.sha256(data).hexdigest()


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--output", type=Path, default=ROOT / "artifacts")
    parser.add_argument("--java", type=Path, required=True)
    parser.add_argument("--public-key", required=True, help="Trusted maintainer Ed25519 public key, 64 hex characters")
    args = parser.parse_args()
    version = (ROOT / "java/VERSION").read_text().strip()
    jar = ROOT / "java/build/libs/ZombieBuddy.jar"
    signature = jar.with_suffix(".jar.zbs")
    if not signature.is_file():
        raise SystemExit("Sign and independently verify ZombieBuddy.jar.zbs before packaging")
    subprocess.run([str(args.java), str(ROOT / "tools/VerifyZbs.java"), args.public_key, str(jar)], check=True, timeout=45)
    runtime = json.loads((ROOT / "artifacts/runtime-tests/results.json").read_text())
    assert runtime["candidateSha256"] == sha(jar.read_bytes()), "Runtime proof belongs to a different JAR"
    assert all(item["passed"] for item in runtime["results"])
    tests = {}
    for suite in ("unitTest", "test_vanilla", "test_patched"):
        reports = list((ROOT / "java/build/test-results" / suite).glob("TEST-*.xml"))
        assert reports, f"Missing {suite} reports"
        total = 0
        for report in reports:
            element = ET.parse(report).getroot()
            assert int(element.get("failures", "0")) == int(element.get("errors", "0")) == 0, report
            total += int(element.get("tests"))
        tests[suite] = total
    commit = subprocess.check_output(["git", "rev-parse", "HEAD"], cwd=ROOT, text=True).strip()
    dirty = bool(subprocess.check_output(["git", "status", "--porcelain"], cwd=ROOT, text=True).strip())
    if dirty:
        raise SystemExit("Commit the reviewed source before producing a traceable preview")
    files = {}
    mod = "Contents/mods/ZombieBuddy/"
    for directory in ("42", "common"):
        for source in sorted((ROOT / directory).rglob("*")):
            if source.is_file():
                files[mod + source.relative_to(ROOT).as_posix()] = source.read_bytes()
    files[mod + "libs/ZombieBuddy.jar"] = jar.read_bytes()
    files[mod + "libs/ZombieBuddy.jar.zbs"] = signature.read_bytes()
    for name in ("LICENSE.txt", "README-COMMUNITY.md", "COMMUNITY_CHANGELOG.md", "THIRD_PARTY_NOTICES.md", "SUPPORT.md"):
        files[name] = (ROOT / name).read_bytes()
    for language in ("FR", "EN"):
        files["doc/Installation_" + language + ".md"] = (ROOT / "doc" / ("Installation_" + language + ".md")).read_bytes()
    files["doc/CommunityEvidence.md"] = (ROOT / "doc/CommunityEvidence.md").read_bytes()
    files["doc/CommunityRelease2Evidence.md"] = (ROOT / "doc/CommunityRelease2Evidence.md").read_bytes()
    for source in sorted((ROOT / "licenses").iterdir()):
        if source.is_file(): files["licenses/" + source.name] = source.read_bytes()
    files["tools/preflight.py"] = (ROOT / "tools/preflight.py").read_bytes()
    files["tools/VerifyZbs.java"] = (ROOT / "tools/VerifyZbs.java").read_bytes()
    files["signing-public-key.txt"] = (args.public_key.lower() + "\n").encode("ascii")
    files["evidence/runtime-tests.json"] = (ROOT / "artifacts/runtime-tests/results.json").read_bytes()
    signature_evidence = json.loads((ROOT / "artifacts/signature-check.json").read_text())
    assert signature_evidence["jarSha256"] == sha(jar.read_bytes())
    assert all(signature_evidence[name] for name in ("validSignature", "tamperedJarRejected", "wrongKeyRejected"))
    files["evidence/signature-check.json"] = (ROOT / "artifacts/signature-check.json").read_bytes()
    localization = json.loads((ROOT / "artifacts/localization-check.json").read_text())
    assert localization["jarSha256"] == sha(jar.read_bytes())
    assert localization["shadedUtf8Catalogue"] and localization["childLocalePropagation"]
    files["evidence/localization-check.json"] = (ROOT / "artifacts/localization-check.json").read_bytes()
    manifest = dict(distribution="ZombieBuddy Community", version=version, upstream="2.3.3",
        upstreamCommit="0ddf161c27848f12d09e74de7fadbea9d50e621d", sourceCommit=commit,
        targetGameSha256=runtime["gameSha256"], tests=tests, gameplayAccepted=False,
        multiplayerAccepted=False, autoUpdate=False, signatureFormat="ZBS Ed25519; not X.509",
        signingPublicKey=args.public_key.lower(),
        files={name: sha(data) for name, data in sorted(files.items())})
    files["manifest.json"] = (json.dumps(manifest, indent=2) + "\n").encode()
    output = args.output.resolve() / version
    output.mkdir(parents=True, exist_ok=True)
    archive = output / ("ZombieBuddyCommunity-" + version + ".zip")
    if archive.exists():
        raise SystemExit("Version already packaged; use a new version rather than overwrite it")
    with ZipFile(archive, "x", ZIP_DEFLATED) as package:
        for name, data in sorted(files.items()):
            info = ZipInfo(name, (2026, 10, 1, 0, 0, 0))
            info.compress_type = ZIP_DEFLATED
            info.external_attr = 0o644 << 16
            package.writestr(info, data)
    with ZipFile(archive) as package:
        assert package.testzip() is None
        for name, expected in manifest["files"].items():
            assert sha(package.read(name)) == expected
        assert not any(name.lower().endswith((".dll", ".exe", ".private.der")) for name in package.namelist())
    (output / "manifest.json").write_bytes(files["manifest.json"])
    (output / "SHA256SUMS.txt").write_text(sha(archive.read_bytes()) + "  " + archive.name + "\n", encoding="ascii")
    print(json.dumps(dict(archive=str(archive),sha256=sha(archive.read_bytes()),files=len(files),tests=tests)))


if __name__ == "__main__":
    main()
