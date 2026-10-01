# Publishing ZombieBuddy Community

Confirmed identity: **ZombieBuddy Community**, maintained by **AZOPEIQBKJ / Aftermath Systems**. Use the existing [GitHub fork](https://github.com/AZOPEIQBKJ/ZombieBuddy), branch `community-42.21`. English-only interface and publication materials. Keep upstream history, the 2.x API and Mod ID `ZombieBuddy`. Any community Workshop item must receive a new ID; `3619862853` is the original item for credit only.

## Prepared

The source branch is public and the repository's default. Issues and private vulnerability reporting are enabled. The signed candidate, installation/rollback guide, licence notices, original artwork, release notes, contact messages and private Workshop folder are prepared. The release remains a **draft prerelease** until actual acceptance. `publication.json` pins the final candidate and remote state.

No game loop, campaign, server, Steam launch setting or Workshop subscription is changed by publication preparation. Tests and isolated JVM checks do not establish gameplay, save/reload or multiplayer support.

## Remaining procedure

1. Send the [prepared messages](CONTACTS.md) to Zed and PZ moderation. Retain dates, links and replies. Do not assume abandonment. The [MIT licence](../LICENSE.txt) covers the licensed code with its notice retained; [PZ policy](https://projectzomboid.com/blog/modding-policy/) separately addresses credits and Workshop maintenance/reupload procedure.
2. Authorize the reversible installation migration and the designated solo tester's session. Before READY, freeze the selected profile, effective JSON + Steam/service arguments, backup/rollback paths, one active distribution/agent, no pending `.jar.new`, and the safe approval fixture. Preserve unrelated settings and approval records. Do not use `allow-all` to bypass an error.
3. Share one actual loading observation with Living Hordes: community version plus `observer Java loaded version=0.36.0-poc`, and no relevant new loading/signature errors. A version watermark alone is insufficient. Reuse previously accepted LH UI/preset observations. Check representative modified approval rendering and persistence; combine the required restart with the relevant LH save/reload scenario where possible. Host and dedicated/client support require separate evidence.
4. Update release notes with the observed scope, then publish the prepared GitHub draft as a **pre-release**. The draft must target the source commit inside the manifest. Its three assets are the named runtime ZIP, `manifest.json` and `SHA256SUMS.txt`. Do not substitute GitHub's source archive or upload the entire preparation kit. Never overwrite a published version or move its tag.
5. After the Workshop contact process and candidate acceptance, copy only `workshop-draft` into the selected profile's `Zomboid/Workshop` directory under a distinct folder name. In the game's Workshop publisher, select it and create a **new private item**. The account holder accepts any Steam terms. Record the returned ID and verify that it is not the original author's ID.
6. Update GitHub/Workshop cross-links with that real ID, review the page and an installation through this channel, then make the item public. Mods that require the old Workshop ID need their own dependency update. Subscribing alone still does not install the Java agent.

## Verification and maintenance

`tools/package_preview.py` packages reviewed, committed source only after matching test/signature/catalogue proofs. `tools/prepare_publication.py` verifies the frozen ZIP/manifest, refuses to overwrite a preparation, and creates a private folder without an item ID. `tools/ValidateWorkshopDraft.java` checks the actual PZ file format in an isolated cache; it never uploads or starts a game loop.

Updates are manual. No upstream Windows installer or native DLL is distributed. The signed upstream author registry remains a network dependency. Keep each correction integrated, reproduced and tested; use a new version when the JAR changes. Review upstream changes without automatically merging the 3.x alpha.
