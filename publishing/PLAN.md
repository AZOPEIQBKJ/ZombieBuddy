# Publishing ZombieBuddy Community

Confirmed identity: **ZombieBuddy Community**, maintained by **AZOPEIQBKJ / Aftermath Systems**. Use the existing [GitHub fork](https://github.com/AZOPEIQBKJ/ZombieBuddy), branch `community-42.21`. English-only interface and publication materials. Keep upstream history, the 2.x API and Mod ID `ZombieBuddy`. Any community Workshop item must receive a new ID; `3619862853` is the original item for credit only.

## Primary release gate

[Transparent replacement for unchanged existing mods](../doc/DropInCompatibility.md) is priority one. The player target is subscription plus activation, with no per-mod adaptation. Manual migration is development setup; the original signed bootstrap and distribution selection are unresolved product gates. Community.3 failed the actual version filter. Community.4 fixes it locally and is installed for development after verified session recovery; it has not been accepted in a new game session.

## Prepared

The source branch is public and the repository's default. Issues and private vulnerability reporting are enabled. The signed candidate, installation/rollback guide, licence notices, original artwork, release notes, contact messages and private Workshop folder are prepared. The release remains a **draft prerelease** until actual acceptance. `publication.json` pins the final candidate and remote state.

Publication preparation itself did not change the installation. On 2026-10-01 the maintainer reported all contact messages sent and separately authorized reversible local migration and the designated solo session with Living Hordes. The local installation was subsequently upgraded to Community.4 with verified backups, one native agent and reviewed Steam options. No Workshop subscription was changed. The first solo attempt was blocked and then interrupted before cleanup; isolated JVM checks do not establish gameplay, save/reload or multiplayer support.

## Remaining procedure

1. Retain dates, links and replies for the [reported outreach](CONTACTS.md); do not resend automatically or assume abandonment. The [MIT licence](../LICENSE.txt) covers the licensed code with its notice retained; [PZ policy](https://projectzomboid.com/blog/modding-policy/) separately addresses credits and Workshop maintenance/reupload procedure.
2. Recovery is complete and Community.4 is installed reversibly for development. A new session GO is required for its actual game acceptance. Migration is performed; keep backups and preserve campaign approval records. The isolated LH profile retains the existing trusted author, one agent and no adapter. Do not use `allow-all` to bypass an error. The separate approval fixture is prepared but inactive; manual approval rendering/persistence remains outstanding if the operator cannot interact with security prompts.
3. Share one actual loading observation with Living Hordes: community version plus `observer Java loaded version=0.36.0-poc`, and no relevant new loading/signature errors. A version watermark alone is insufficient. Reuse previously accepted LH UI/preset observations. Check representative modified approval rendering and persistence; combine the required restart with the relevant LH save/reload scenario where possible. Host and dedicated/client support require separate evidence.
4. Update release notes with the observed scope, then publish the prepared GitHub draft as a **pre-release**. The draft must target the source commit inside the manifest. Its three assets are the named runtime ZIP, `manifest.json` and `SHA256SUMS.txt`. Do not substitute GitHub's source archive or upload the entire preparation kit. Never overwrite a published version or move its tag.
5. After the Workshop contact process and candidate acceptance, copy only `workshop-draft` into the selected profile's `Zomboid/Workshop` directory under a distinct folder name. In the game's Workshop publisher, select it and create a **new private item**. The account holder accepts any Steam terms. Record the returned ID and verify that it is not the original author's ID.
6. Update GitHub/Workshop cross-links with that real ID, review the page and an installation through this channel, then make the item public. Mods that require the old Workshop ID need their own dependency update. Subscribing alone still does not install the Java agent.

## Verification and maintenance

`tools/package_preview.py` packages reviewed, committed source only after matching test/signature/catalogue proofs. `tools/prepare_publication.py` verifies the frozen ZIP/manifest, refuses to overwrite a preparation, and creates a private folder without an item ID. `tools/ValidateWorkshopDraft.java` checks the actual PZ file format in an isolated cache; it never uploads or starts a game loop.

Updates are manual. No upstream Windows installer or native DLL is distributed. The signed upstream author registry remains a network dependency. Keep each correction integrated, reproduced and tested; use a new version when the JAR changes. Review upstream changes without automatically merging the 3.x alpha.
