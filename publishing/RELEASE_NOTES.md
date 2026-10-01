# ZombieBuddy Community 2.3.3-community.1 — development preview

**Draft for the frozen local candidate. No release has been published.** If a new JAR, translation or runtime asset is produced, assign a new version and replace these details with its own evidence; never reuse this candidate's hashes or acceptance.

## English

Independent community maintenance of ZombieBuddy 2.3.3 by Andrey "Zed" Zaikin. The PZ 42.21 `List` loading correction is built into the framework, with the public 2.x `ArrayList` entry point retained. No Aftermath compatibility adapter is required.

This preview also fixes Windows paths in the existing `patches_jar` API, HTTP-client initialization failures, root-path handling and Lua options when the agent is absent. Framework updates are manual. Mod signatures, the signed upstream author registry and approval rules remain in place.

**Evidence:** 148 Java tests, Go installer regressions, four preflight checks, isolated Java/native-loader and Kahlua checks using the audited game, plus positive and negative Ed25519 checks. **Not accepted:** actual game loading, save/reload, Host, dedicated/client multiplayer, Linux or macOS. The installer code has tests but no installer or native DLL is shipped. Other mods still need their own PZ 42.21 port and acceptance.

Read the community installation and rollback guide in the archive. Keep one distribution/agent and review all launch sources. Use the named release ZIP for the prepared runtime package, not GitHub's automatically generated source archive. Do not use the original Windows installer as a community updater.

## Français

Maintenance communautaire indépendante de ZombieBuddy 2.3.3 de Zed. Le correctif de chargement `List` sous PZ 42.21 est intégré ; l'entrée publique `ArrayList` de l'API 2.x est conservée. Aucun adaptateur Aftermath n'est requis.

Cette preview corrige aussi les chemins Windows de l'API `patches_jar`, l'initialisation HTTP, les chemins racine et les options Lua sans agent. Les mises à jour du framework sont manuelles. Signatures des mods, registre d'auteurs signé et règles d'approbation restent actifs.

**Contrôles locaux :** 148 tests Java, régressions Go, quatre contrôles de prévol, JVM isolées Java/chargeur natif et Kahlua sur le jeu audité, signatures Ed25519 valides/altérées/mauvaise clé. **Non réceptionné :** chargement en partie, sauvegarde/reprise, Host, dédié et clients MP, Linux ou macOS. Aucun installateur ni DLL native n'est livré. Les autres mods demandent leur propre portage et réception 42.21.

Lire le guide d'installation et de retour arrière du paquet, conserver un seul agent et une seule distribution. Utiliser le ZIP nommé de la release, pas l'archive de sources générée automatiquement par GitHub. Ne pas utiliser l'installateur original pour mettre le fork à jour.

## Frozen artefacts / Artefacts figés

- Source commit: `af69c7fe04e4c82bac42055c2f5072e50e643106`.
- ZIP: `ZombieBuddyCommunity-2.3.3-community.1.zip`.
- ZIP SHA-256: `2007d4db768970f549367cd9b15d0c52131564670c17e6c173c172f43ac83430`.
- Agent SHA-256: `a0711c36597795348fa94d31d400e34c04eaa263a4e0a3bb3fb3d6898395eb6a`.
- Attach the ZIP, its matching `manifest.json` and `SHA256SUMS.txt` as release assets.
- `.zbs` identity: Aftermath Systems, Ed25519. This is not an X.509 signature by Zed. Establish the public key from the maintainer's independently verified identity, not solely from the download carrying it.

Original work and MIT copyright: Andrey "Zed" Zaikin. Community maintenance: Aftermath Systems. Related proposals: benoitthore #56, yuruichang #58. Full notices and licence texts are included in the ZIP.

Before publishing, link the guides to the exact accepted release tag and set **pre-release**, not **latest stable**. If community.1 remains internal, do not create a public tag/release for it; publish the later accepted candidate instead.
