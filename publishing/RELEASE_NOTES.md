# ZombieBuddy Community 2.3.3-community.2 — development preview

Draft prerelease. Actual gameplay acceptance is pending.

## English

Independent community maintenance of ZombieBuddy 2.3.3 by Andrey "Zed" Zaikin. The PZ 42.21 `List` loading correction is built into the framework, with the public 2.x `ArrayList` entry point retained. No Aftermath compatibility adapter is required.

Community.2 adds original community artwork and structured French/English approval text. JVM locale selects the language; `-Dzombiebuddy.language=fr|en` overrides it and propagates to the approval subprocess.

This preview also fixes Windows paths in the existing `patches_jar` API, HTTP-client initialization failures, root-path handling and Lua options when the agent is absent. Framework updates are manual. Mod signatures, the signed upstream author registry and approval rules remain in place.

**Evidence:** 160 Java tests, Go installer regressions, four preflight checks, isolated Java/native-loader and Kahlua checks using the audited game, plus positive and negative Ed25519 checks. **Not accepted:** actual game loading, save/reload, Host, dedicated/client multiplayer, Linux or macOS. The installer code has tests but no installer or native DLL is shipped. Other mods still need their own PZ 42.21 port and acceptance.

Read the community installation and rollback guide in the archive. Keep one distribution/agent and review all launch sources. Use the named release ZIP for the prepared runtime package, not GitHub's automatically generated source archive. Do not use the original Windows installer as a community updater.

## Français

Maintenance communautaire indépendante de ZombieBuddy 2.3.3 de Zed. Le correctif de chargement `List` sous PZ 42.21 est intégré ; l'entrée publique `ArrayList` de l'API 2.x est conservée. Aucun adaptateur Aftermath n'est requis.

Community.2 intègre les visuels communautaires originaux et les dialogues FR/EN structurés. La langue JVM est utilisée, avec surcharge `-Dzombiebuddy.language=fr|en` transmise au sous-processus.

Cette preview corrige aussi les chemins Windows de l'API `patches_jar`, l'initialisation HTTP, les chemins racine et les options Lua sans agent. Les mises à jour du framework sont manuelles. Signatures des mods, registre d'auteurs signé et règles d'approbation restent actifs.

**Contrôles locaux :** 160 tests Java, régressions Go, quatre contrôles de prévol, JVM isolées Java/chargeur natif et Kahlua sur le jeu audité, signatures Ed25519 valides/altérées/mauvaise clé. **Non réceptionné :** chargement en partie, sauvegarde/reprise, Host, dédié et clients MP, Linux ou macOS. Aucun installateur ni DLL native n'est livré. Les autres mods demandent leur propre portage et réception 42.21.

Lire le guide d'installation et de retour arrière du paquet, conserver un seul agent et une seule distribution. Utiliser le ZIP nommé de la release, pas l'archive de sources générée automatiquement par GitHub. Ne pas utiliser l'installateur original pour mettre le fork à jour.

## Frozen artefacts / Artefacts figés

The source commit and ZIP SHA-256 will be filled from the frozen manifest before upload. Attach the named ZIP, its `manifest.json` and `SHA256SUMS.txt`.

Signing identity: Aftermath Systems, Ed25519 `.zbs`, not Zed's X.509 signature. Public key: `989ac279f40f1a35fa0616645e319fc44cde9a15a842b7c365870147bfff3ce0`. Establish this key independently from the maintainer's verified identity, not solely from the download carrying it.

Original work and MIT copyright: Andrey "Zed" Zaikin. Community maintenance: Aftermath Systems. Related proposals: benoitthore #56, yuruichang #58. Full notices and licence texts are included in the ZIP. Keep this release in draft until acceptance, then publish as pre-release with the actual supported scope, not as latest stable.
