# Community.3 English-only candidate evidence

On 1 October 2026 the maintainer explicitly requested English only. This supersedes the earlier bilingual choice for this fork. No game session or deployment is authorized by this change.

Verified source contracts: approval frontends share the structured English catalogue; Swing no longer needs a language argument because the catalogue is invariant. Console uses its original y/n input contract. Signature presentation remains separate from raw diagnostic fields and trust decisions. Remove French Java and Lua resources and bilingual publication text, while retaining upstream history and earlier frozen archives.

Clean compilation is required to exclude stale French resources. Check the final shaded JAR under a non-English JVM locale and an obsolete language override; the displayed text must remain English and no French catalogue may be present. Existing approval denial/persistence and signature cases remain in the regression suite. No new gameplay API, networking hook or save mutation is introduced.

The audited PZ 42.21.0 game SHA-256 remains `e1a69eb743ede60b213a0fe7f8b83d4fcab773036d256cc4543a336f3b058a33`. Actual game/Host/dedicated acceptance remains pending.

## Local results

- Clean Gradle build: **158 Java tests passed** (113 unit, 30 patched, 15 vanilla). The smaller total than community.2 removes two obsolete locale-switching cases; signature and approval regression coverage remains.
- Five isolated checks passed: original baseline, integrated List handling, duplicate-agent first-policy retention, original native loader, and Kahlua options/missing-agent notice. The Kahlua harness now checks required English keys and absence of French resources.
- `CatalogueSmoke.java` passed against the final shaded JAR with a French JVM default and obsolete French override: displayed text remains English; French catalogue absent.
- Ed25519 valid-signature, tampered-JAR and wrong-key checks passed for JAR SHA-256 `40c73114a14614ef6288f57657d77aafbe959f5ea550c3c6730ea0d5c85f1702`. Signing identity unchanged.
- Source/bytecode proof for Lua fallback: PZ 42.21 `Translator.forLanguageStack` adds the default language and iterates the stack; `lambda$loadFiles$1` calls `tryFillMapFromMods` for each language. `Languages` initializes its default as EN/English. This supports shipping only EN resources; it is not a rendered-game observation.
- Four preflight tests and unchanged Go installer checks were already passed; no code in those paths changed in this English-only slice. No public game JAR, native DLL, installer or private key is included in the runtime ZIP.

No game loop or new console.txt was produced. Actual approval layout/persistence, solo loading and save/reload remain pending, as do Host, dedicated/client, Linux/macOS and real multiplayer. No installation, global launch setting, subscription or campaign save changed. See the release manifest for the exact source commit and archive hash.
