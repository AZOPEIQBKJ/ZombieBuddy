# ZombieBuddy Community 2.3.3-community.3 — development preview

Draft prerelease; actual gameplay acceptance is pending.

Independent maintenance of ZombieBuddy 2.3.3 by Andrey "Zed" Zaikin. The PZ 42.21 List-based loading fix is integrated into the framework; the public 2.x ArrayList entry point is retained. No Aftermath compatibility adapter is required.

Community.3 uses **English only** for the interface, documentation and publication materials, with original community artwork. It supersedes the unpublished bilingual candidate. Existing loading, Windows path, HTTP initialization, root-path and missing-agent fixes remain integrated. Updates are manual; signature checks, the signed upstream author registry and approval rules remain active.

158 Java tests, isolated Java/native-loader and Kahlua checks against the audited PZ 42.21 build, preflight tests, packaged-catalogue verification and positive/negative Ed25519 checks passed for this candidate. The final JAR stays English under a non-English JVM locale and contains no French catalogue. Actual loading, approval rendering/persistence, save/reload, Host, dedicated/client multiplayer, Linux and macOS remain unaccepted. The upstream installer code is unchanged and previously tested; no installer or native DLL is shipped.

Read `doc/Installation_EN.md` in the archive for migration and rollback. Keep one distribution and one effective agent across all launch settings. Use the named release ZIP, not GitHub's generated source archive. Other mods still require their own 42.21 port and acceptance.

## Frozen assets

Source commit and archive hashes are filled from the frozen manifest before upload. Attach the named ZIP, `manifest.json` and `SHA256SUMS.txt`.

Signing identity: Aftermath Systems, Ed25519 `.zbs`, not Zed's X.509 signature. Public key: `989ac279f40f1a35fa0616645e319fc44cde9a15a842b7c365870147bfff3ce0`. Establish this key independently from the maintainer's verified identity.

Original framework and MIT copyright: Andrey "Zed" Zaikin. Community maintenance: Aftermath Systems. Related proposals: benoitthore #56 and yuruichang #58; the full #58 rewrite is not included. Licence and library notices accompany the package.
