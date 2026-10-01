"""Verify the frozen preview and prepare an offline publication draft. Never uploads."""
import argparse
import hashlib
import json
from pathlib import Path, PurePosixPath
import shutil
import subprocess
from zipfile import ZipFile

ROOT = Path(__file__).resolve().parents[1]


def digest(data):
    return hashlib.sha256(data).hexdigest()


def verified_payload(archive, config):
    if digest(archive.read_bytes()) != config["candidateZipSha256"]:
        raise ValueError("Candidate ZIP differs from the frozen SHA-256")
    with ZipFile(archive) as package:
        names = package.namelist()
        if len(names) != len(set(names)):
            raise ValueError("Duplicate archive entry")
        for name in names:
            path = PurePosixPath(name)
            if path.is_absolute() or ".." in path.parts or "\\" in name or ":" in name:
                raise ValueError("Unsafe archive path")
        manifest = json.loads(package.read("manifest.json"))
        if set(names) != set(manifest["files"]) | {"manifest.json"}:
            raise ValueError("Manifest does not enumerate the complete archive")
        if manifest["sourceCommit"] != config["candidateSourceCommit"]:
            raise ValueError("Candidate source commit mismatch")
        if manifest["version"] != config["candidateVersion"]:
            raise ValueError("Candidate version mismatch")
        payload = {name: package.read(name) for name in names}
        for name, expected in manifest["files"].items():
            if digest(payload[name]) != expected:
                raise ValueError("Manifest hash mismatch: " + name)
        jar = payload["Contents/mods/ZombieBuddy/libs/ZombieBuddy.jar"]
        if digest(jar) != config["candidateJarSha256"]:
            raise ValueError("Candidate agent mismatch")
        return payload, manifest


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--candidate", type=Path, required=True)
    parser.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()
    config = json.loads((ROOT / "publishing/publication.json").read_text(encoding="utf-8"))
    if config["workshopItemId"] is not None or config["initialWorkshopVisibility"] != "private":
        raise SystemExit("This preparer only creates a new private draft without a Workshop ID")
    payload, manifest = verified_payload(args.candidate, config)
    artwork = ROOT / "publishing/assets/workshop-cover.png"
    png = artwork.read_bytes()
    if png[:8] != b"\x89PNG\r\n\x1a\n" or tuple(int.from_bytes(png[n:n+4], "big") for n in (16, 20)) != (256, 256):
        raise SystemExit("Expected an original 256x256 PNG workshop cover")
    output = args.output.resolve()
    if output.exists():
        raise SystemExit("Output already exists; choose a new directory, never overwrite a frozen draft")
    output.mkdir(parents=True)
    draft = output / "workshop-draft"
    mod = draft / "Contents/mods/ZombieBuddy"
    for name, data in payload.items():
        if name.startswith("Contents/"):
            target = draft / name
            target.parent.mkdir(parents=True, exist_ok=True)
            target.write_bytes(data)
    # Workshop uploads Contents, so include notices with the Lua and Java payload.
    for name, data in payload.items():
        if name in ("LICENSE.txt", "THIRD_PARTY_NOTICES.md", "README-COMMUNITY.md", "COMMUNITY_CHANGELOG.md", "SUPPORT.md") or name.startswith(("licenses/", "doc/")):
            target = mod / name
            target.parent.mkdir(parents=True, exist_ok=True)
            target.write_bytes(data)
    (draft / "preview.png").write_bytes(png)
    en = (ROOT / "publishing/WORKSHOP_EN.txt").read_text(encoding="utf-8")
    description = en
    # Format verified from PZ 42.21 Workshop/ModTemplate/workshop.txt and its reader.
    metadata = ["version=1", "title=ZombieBuddy Community — development preview"]
    metadata.extend("description=" + line for line in description.splitlines())
    metadata.extend(["tags=", "visibility=private"])
    (draft / "workshop.txt").write_text("\n".join(metadata) + "\n", encoding="utf-8")
    # Copy the preparation material and untouched candidate alongside, outside Contents.
    materials = output / "materials"
    shutil.copytree(ROOT / "publishing", materials / "publishing")
    for folder in ("doc", "licenses", ".github"):
        shutil.copytree(ROOT / folder, materials / folder)
    for name in ("README.md", "README-COMMUNITY.md", "SUPPORT.md", "LICENSE.txt",
                 "THIRD_PARTY_NOTICES.md", "COMMUNITY_CHANGELOG.md"):
        shutil.copyfile(ROOT / name, materials / name)
    (materials / "tools").mkdir()
    for name in ("prepare_publication.py", "ValidateWorkshopDraft.java", "preflight.py", "VerifyZbs.java"):
        shutil.copyfile(ROOT / "tools" / name, materials / "tools" / name)
    release = output / "reference-candidate"
    release.mkdir()
    shutil.copyfile(args.candidate, release / args.candidate.name)
    (release / "manifest.json").write_bytes(payload["manifest.json"])
    (release / "SHA256SUMS.txt").write_text(
        config["candidateZipSha256"] + "  " + args.candidate.name + "\n", encoding="ascii")
    status = dict(prepared=True, uploaded=False, publicReleaseReady=False,
        preparationSourceCommit=subprocess.check_output(["git", "rev-parse", "HEAD"], cwd=ROOT, text=True).strip(),
        preparationSourceDirty=bool(subprocess.check_output(["git", "status", "--porcelain"], cwd=ROOT, text=True).strip()),
        candidateVersion=manifest["version"], candidateSourceCommit=manifest["sourceCommit"],
        candidateZipSha256=config["candidateZipSha256"],
        candidateJarSha256=config["candidateJarSha256"],
        workshopVisibility="private", workshopItemId=None,
        retainsFrozenRuntimeAssets=True, originalCommunityArtworkInRuntime=True,
        outstanding=config["outstanding"])
    (output / "PREPARATION_STATUS.json").write_text(json.dumps(status, indent=2) + "\n", encoding="utf-8")
    (output / "READ_ME_FIRST.txt").write_text(
        "LOCAL WORKSHOP PREPARATION\n"
        "Not ready for public release. No upload, install or game launch was performed.\n"
        "workshop-draft uses the frozen community.3 runtime with original icons and English-only text.\n"
        "The GitHub source branch and draft release are tracked separately from Workshop upload.\n"
        "Read materials/publishing/PLAN.md and PREPARATION_STATUS.json before further work.\n"
        "materials is a documentation/tool kit, not a complete buildable source checkout.\n"
        "Do not upload this entire preparation directory or use the original Workshop ID.\n",
        encoding="utf-8")
    hashes = {p.relative_to(output).as_posix(): digest(p.read_bytes())
              for p in sorted(output.rglob("*")) if p.is_file()}
    (output / "FILES_SHA256.json").write_text(json.dumps(hashes, indent=2) + "\n", encoding="utf-8")
    print(json.dumps(dict(output=str(output), files=len(hashes), uploaded=False,
        publicReleaseReady=False, candidateVerified=True)))


if __name__ == "__main__":
    main()
