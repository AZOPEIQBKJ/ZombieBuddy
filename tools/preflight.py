"""Read-only inspection of an installation before migrating to the community preview."""
import argparse
import hashlib
import json
from pathlib import Path
import sys

PZ_SHA256 = "e1a69eb743ede60b213a0fe7f8b83d4fcab773036d256cc4543a336f3b058a33"


def inspect_options(json_args, extra_args, reviewed):
    combined = [*json_args, *extra_args]
    agents = [arg for arg in combined if "zombiebuddy.jar" in arg.lower()
              and "-javaagent:" in arg.lower() or "-agentlib:zbnative" in arg.lower()
              or "-agentpath:" in arg.lower() and "zbnative" in arg.lower()]
    issues = []
    if not reviewed:
        issues.append("Steam/service/script options have not been reviewed; this is not a complete preflight.")
    if len(agents) != 1:
        issues.append(f"Expected exactly one ZombieBuddy agent across all launch sources; found {len(agents)}.")
    if any("aftermathlhcompat4221" in arg.lower() for arg in combined):
        issues.append("Remove the obsolete Aftermath 42.21 adapter entry when migrating to the integrated fork.")
    return agents, issues


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--game", type=Path, required=True)
    parser.add_argument("--launcher-json", type=Path)
    parser.add_argument("--extra-jvm-arg", action="append", default=[])
    parser.add_argument("--extra-options-reviewed", action="store_true",
                        help="Only use after reviewing Steam/service/script options; provide all extra agent arguments")
    args = parser.parse_args()
    game = args.game.resolve()
    launcher = args.launcher_json or game / "ProjectZomboid64.json"
    data = json.loads(launcher.read_text(encoding="utf-8-sig"))
    agents, issues = inspect_options(data.get("vmArgs", []), args.extra_jvm_arg, args.extra_options_reviewed)
    jar = game / "projectzomboid.jar"
    sha = hashlib.sha256(jar.read_bytes()).hexdigest()
    if sha != PZ_SHA256:
        issues.append("The game JAR is not the audited PZ 42.21 build.")
    if (game / "ZombieBuddy.jar.new").exists():
        issues.append("A pending ZombieBuddy.jar.new exists: the native launcher could overwrite the fork before Java starts.")
    print(json.dumps(dict(gameSha256=sha, zombieBuddyAgents=agents, issues=issues,
        readyForReviewedMigration=not issues, modifiesFiles=False,
        limits="Does not verify enabled mods, Workshop subscriptions, running processes or gameplay."), indent=2))
    return 1 if issues else 0


if __name__ == "__main__":
    sys.exit(main())
