# ZombieBuddy Community

Independent maintenance preview by Aftermath Systems, based on ZombieBuddy 2.3.3 by Andrey "Zed" Zaikin. **The PZ 42.21 loading fix is integrated in this agent. No Aftermath compatibility adapter is required.** The community name was confirmed on 1 October 2026. This is not an official upstream release.

Fork de maintenance indépendant par Aftermath Systems, basé sur ZombieBuddy 2.3.3 de Zed. **Le correctif PZ 42.21 est intégré dans cet agent. Aucun adaptateur de compatibilité Aftermath n'est requis.** Nom confirmé le 1er octobre 2026. Les réceptions en partie et en multijoueur restent à effectuer.

- [Installation en français](doc/Installation_FR.md)
- [Installation in English](doc/Installation_EN.md)
- [Changes and acceptance limits](COMMUNITY_CHANGELOG.md)
- [Current evidence](doc/CommunityRelease2Evidence.md) and [original regression baseline](doc/CommunityEvidence.md)
- [MIT licence](LICENSE.txt) and [third party notices](THIRD_PARTY_NOTICES.md)

The preview preserves the public 2.x packages, annotations, `ZombieBuddy` Mod ID and existing mod-approval policy. Enable only one ZombieBuddy distribution. Automatic replacement of the framework JAR is disabled for previews; updates are manual. `.zbs` signatures and the verified author registry remain enforced as before.

## Build and verification

Use JDK 25 and the Gradle wrapper under `java/`. The wrapper pins Gradle 9.3.1 and its official SHA-256. Set `PZ_CLASSPATH` to the legally installed, audited PZ 42.21 `projectzomboid.jar` (comma-separated paths if needed), then run `gradlew.bat test shadowJar` on Windows or `./gradlew test shadowJar` on Unix. The game JAR is a compile/test input and must never be packaged. Windows tests are locally executed; running this command on Linux is not yet proof of Linux compatibility.

There is no automatic use of Zed's private signing key or macOS keychain. Sign the final JAR's `.zbs` with the maintainer's existing Ed25519 author identity outside the repository, independently verify it, and keep private keys out of source and logs.

Run `python tools/test_runtime.py --game GAME_DIRECTORY --jdk JDK_DIRECTORY --baseline UNMODIFIED_233_JAR` for isolated before/after tests using real B42 bytecode and Kahlua. It never starts a game loop. `python tools/test_preflight.py` checks the combined-launch regression. The upstream Go installer fixes can be tested with `go test ./...` under `installer/`; that installer is not shipped by this preview because its distribution identity still targets upstream.

After reviewing and committing source, `python tools/package_preview.py --java JAVA_EXECUTABLE --public-key TRUSTED_PUBLIC_KEY_HEX` requires passing JVM reports, matching runtime/localization-test JAR hashes and a valid signature sidecar. It independently verifies the signature using the JDK, then creates a versioned ZIP, per-file manifest and SHA256SUMS without installing or uploading them. A public key included in a download is informational; establish trust from the maintainer's previously verified identity.

The upstream Git history retains the original README and documentation, whose download links describe the original project. Use these community instructions for this distribution. The preview archive contains the guides under `doc/` and does not include a Gradle build workspace; build from the source checkout. The publication manifest identifies the frozen runtime source commit separately from later editorial preparation commits.
