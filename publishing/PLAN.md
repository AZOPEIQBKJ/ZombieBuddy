# Publishing ZombieBuddy Community

Current preparation: 2 October 2026. Identity: **ZombieBuddy Community**, AZOPEIQBKJ / Aftermath Systems, [GitHub fork](https://github.com/AZOPEIQBKJ/ZombieBuddy), branch `community-42.21`. Framework interface/publication text is English-only. Retain upstream history, MIT credit and the `ZombieBuddy` Mod ID. A future Workshop item must get a new ID; `3619862853` is credit/reference only.

## Current evidence

[Release status](../doc/CommunityReleaseStatus.md) records Community.4 normal Steam startup and simultaneous Java loading of LH 0.36.1, Loadout poc3 and Bulk Storage Native 0.1.8. The user observation and logs were reviewed by Testeur solo. These loading mechanisms are acquired and must not be replayed merely to fill an old request naming LH 0.36.0. Historical Community.3/installer preview.2 failures remain intact.

A one-time Community installation is accepted. Mods deliberately pinned to an excluded exact framework version may require author updates; minimum-only declarations are implemented in the owned modules. Keep framework fixes integrated, signature/approval checks intact, and no per-consumer adapter requirement. This is not an all-mod compatibility guarantee.

## Prepared artifacts

- Immutable runtime Community.4, source `f6cc18d3ff2706bb49d2cebace4bae05da7070af`.
- Immutable Windows installer preview.3, source `ad5143190fc8c84593b448ee6456e291d3c85264`, embedding the runtime and pinned MIT native bootstrap with source/notices. No game/JRE binary is distributed.
- Current release notes, installer guide, original artwork, support/rollback instructions and private Workshop preparation. No Workshop item ID or upload.
- Existing Community.3 draft remains blocked historical material. Prepare a distinct Community.4 draft, preserving its runtime source target and installer provenance. Do not publish automatically.

## Remaining procedure

1. Validate archive and remote-asset digests; attach the runtime ZIP/manifest/checksums, installer ZIP/separately named manifest/checksums and current INSTALLATION.md to the Community.4 draft. Never overwrite immutable dist archives or a published release/tag.
2. Complete the final installer executable interaction/transaction and representative approval persistence gap. The current installation was repaired during development: its observed startup and a zero-change installer plan do not establish fresh installation. Reuse local fixtures and observed loading. Since 2 October, the user performs game tests manually; the central tester, queue, reservations and READY/GO submissions are retired. After source/local/log checks, provide only short manual steps for the remaining mechanism and its expected result. Do not launch or control a game or recreate the former orchestration.
3. Retain only the genuinely missing gameplay/save-load mechanisms in the owning module requests. Host/dedicated and other-platform claims require separate evidence. Do not turn every module feature into a duplicate framework startup test.
4. Refresh release notes with the actual scope and obtain the release decision before making the draft a public prerelease. The old runtime archive guide predates the native-bootstrap repair; the release must prominently direct Windows players to the current INSTALLATION.md/installer route.
5. Preserve the recorded author outreach and review any reply/current Workshop requirements before uploading. Prepare a new private Workshop item using only `workshop-draft`; never upload the surrounding tools/evidence kit or reuse the author's Workshop ID. Account terms and creation of a Workshop item belong to the account holder.
6. Record the new ID, update cross-links, review the page and channel-specific installation route, then decide public visibility. Keep unchanged original Workshop dependencies handled centrally and disclose local-mod precedence.

## Preservation and maintenance

Do not restore historical profiles or old module versions over CURRENT. Keep saves, approvals, launch options and unrelated mod files intact. Automatic Community updates are not implemented. Binary signing identity, native bootstrap provenance and unsupported routes are documented in the installer guide.

`tools/prepare_publication.py` verifies the runtime's complete frozen manifest and prepares a new private local folder. Current installation/status documentation is overlaid as editorial material while the runtime JAR and signature remain byte-identical. `tools/ValidateWorkshopDraft.java` validates metadata/content using local PZ classes without a game loop or upload.
