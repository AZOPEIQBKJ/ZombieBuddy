# ZombieBuddy Community

Independent maintenance of [ZombieBuddy by Andrey "Zed" Zaikin](https://github.com/zed-0xff/ZombieBuddy), based on version 2.3.3. Maintained by Aftermath Systems under the original MIT licence, with the upstream history and credits preserved.

**The PZ 42.21 loading fix is integrated into ZombieBuddy itself. No separate Aftermath compatibility adapter is needed.**

**Current status: local development preview `2.3.3-community.1`.** 148 Java tests and isolated checks against the actual 42.21 game bytecode pass. Real gameplay, save/reload, Host and dedicated-server acceptance remain pending. There is no public community download or Workshop item yet. Linux and macOS support are not established by the Windows checks.

## Players / Joueurs

| English | Français |
| --- | --- |
| [Installation, migration and rollback](doc/Installation_EN.md) | [Installation, migration et retour arrière](doc/Installation_FR.md) |
| [Changes and current limits](COMMUNITY_CHANGELOG.md) | [Préparatifs de publication](publishing/PLAN_FR.md) |
| [Support and bug reports](SUPPORT.md) | [Support et signalement des bugs](SUPPORT.md) |

Ce fork vise à maintenir le chargement des mods Java sous PZ 42.21. Le correctif est intégré au framework. La candidate locale est vérifiée au banc ; elle n'est pas encore réceptionnée en partie ou en multijoueur.

Keep **one** active distribution with Mod ID `ZombieBuddy`, and **one** effective ZombieBuddy agent across all launch settings. Existing 2.x namespaces and the Mod ID are preserved for compatibility. A mod that hard-codes the original Workshop item may still need a dependency update by its maintainer. Installing the framework does not automatically port every dependent mod to 42.21.

Updates are manual. The original Windows installer is not a community installer. The candidate includes a Java agent and its `.zbs` signature; no installer, native DLL or game files are distributed. Mod approvals and signature checks remain active. Only run Java mods from sources you trust: they execute with the permissions of the game process.

## Developers

- [Build and local verification](README-COMMUNITY.md)
- [Evidence, baseline and acceptance limits](doc/CommunityEvidence.md)
- [Upstream API guide](doc/ModdingGuide.md) — retain 2.x API compatibility; examples are not a blanket B42.21 guarantee.
- [Publication materials and release gates](publishing/PLAN_FR.md)

Build with JDK 25 and the pinned Gradle wrapper, using your legally installed game as a compile/test input. Never include the game JAR or private signing keys in a public repository or release. Legacy Rake/installer workflows retain upstream assumptions; use the community build and packaging instructions.

## Credits and licence

Original framework: **Andrey "Zed" Zaikin**. Community maintenance: **Aftermath Systems**. Related reports and proposals include benoitthore's PR #56 and yuruichang's PR #58; this fork does not incorporate the full PR #58 rewrite. See [MIT licence](LICENSE.txt), [dependency notices](THIRD_PARTY_NOTICES.md) and [visual-asset audit](publishing/ASSETS_AND_LOCALIZATION.md).

The original README and distribution instructions remain in the [upstream 2.3.3 history](https://github.com/zed-0xff/ZombieBuddy/blob/0ddf161c27848f12d09e74de7fadbea9d50e621d/README.md). They describe the original distribution.
