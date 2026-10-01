# Priority one: transparent replacement for existing ZombieBuddy mods

The maintainer clarified the product requirement on 2026-10-01: players may install ZombieBuddy Community once, as they previously installed Zed's distribution, then enable it for their save and continue using existing ZombieBuddy-dependent mods without work by those mods' authors. A one-time Community installation or migration is explicitly accepted. Subscription-only migration is no longer a release requirement. Compatibility with unchanged consumers remains the primary release gate; a Living Hordes-only fix does not satisfy it.

## Compatibility contract

- Preserve the technical `ZombieBuddy` Mod ID, upstream Java packages/classes/method descriptors, annotations and Lua-facing contracts used by existing 2.x consumers.
- Preserve existing `require=\\ZombieBuddy`, version limits, JARs, signatures, approval records and saves. Do not edit/rebuild consumer mods or require individual adapters to make the fork work.
- Handle framework compatibility regressions inside Community. Living Hordes is one unchanged representative consumer, not the product boundary.
- Treat existing Workshop dependencies on the original item as a migration requirement to solve centrally. Do not make third-party maintainers retarget their pages or republish their mods as the accepted solution.
- Establish a deterministic single active framework when Steam still has the original item installed. A shared Mod ID alone does not prove which physical distribution is selected.
- Validate the documented Community installation/migration and activation path from a real original installation. Back up and replace the framework, preserve approvals, saves and unrelated launch settings, and provide rollback. Existing machine-specific development helpers are evidence for preparation, not a finished user installer.
- Do not promise every mod works merely because its framework dependency loads. Define representative coverage of the public framework mechanisms and record concrete remaining incompatibilities; unrelated changes in the game or consumer code must remain visible rather than be silently waived.

## Source findings that determine installation requirements

1. The installed Windows native agent loads `ZombieBuddy.jar` from the game working directory before normal mod loading. Its updater applies only an already staged adjacent `.new` file; it does not discover a newly subscribed Community item. Source: `c/windows/zbNative.c`, `Agent_OnLoad`, `check_and_apply_update`.
2. Original 2.3.3's `SelfUpdater.checkAndUpdateIfNewer` accepts a replacement only if its JAR entries have an X.509 signature matching Zed's hard-coded certificate fingerprint. Our independent `.zbs`/Ed25519 signature is not that certificate. We do not possess or impersonate Zed's signing identity. Source: upstream commit `0ddf161c27848f12d09e74de7fadbea9d50e621d`, `java/.../SelfUpdater.java`.
3. On audited B42.21, original 2.3.3's `loadMods(ArrayList)` advice misses the game's `loadMods(List)` method. This can prevent reaching normal Java-mod discovery and its updater at all. Existing isolated native tests reproduce this boundary.
4. Default mod search order is `workshop,steam,mods`. The installer selects local `mods` first, retaining the relative order of other sources. A real-engine probe of `getAllModFolders` and `ChooseGameInfo.getModDetails`, using a simulated installed Workshop folder list, selects the original with the default order and Community with the installed order; the unchanged consumer remains discoverable. Other local duplicate Mod IDs also become preferred, which is disclosed in the preview. Actual Steam UI/game acceptance remains pending. Source: installed B42.21 `MainScreenState.main`, `ZomboidFileSystem.getAllModFolders`, `setModIdToDir` and `ChooseGameInfo` bytecode; `tools/InstallerModSelectionSmoke.java`.
5. The reviewed vanilla Lua file-writing routes target the Lua cache or a mod's common folder and reject relative traversal; no supported Workshop-only replacement of the preloaded root agent has been established. Do not invent an engine API, exploit a path restriction or disable signature checks to claim this gate passed.

These findings prevent using the original automatic updater for an independently signed Community JAR. They do not prevent a user-run Community installer from replacing the root agent while the game is stopped. The maintainer explicitly accepted that route; Zed's signing identity or cooperation is not required for this replacement. Preserve signature and approval checks for dependent mods. A future automatic Community updater can use a Community trust identity, but it is not implemented by the current manual-update preview.

## Ordered next work and acceptance

1. Completed: recover the interrupted QA session, retaining its blocked result and unchanged consumer evidence.
2. Completed locally: package and deploy Community.4 with the version-filter correction and original callable signatures. Actual game acceptance remains pending.
3. Prepared locally: the independent [Windows installer](CommunityInstaller.md) provides preview, backup hashes, one effective agent and rollback. Fresh installation, original replacement, interruptions, stale plans, later Steam activity and rollback are exercised in fixtures. Engine-level distribution selection passes with a simulated Steam folder list. Validate the complete normal Steam/game journey next; no consumer metadata changes or automatic subscription changes.
4. Build a small coverage matrix of unchanged consumers using distinct framework mechanisms: dependency/version parsing, legacy compiled loader calls, annotated patches, Lua exposure, approvals/preloads and save/reload. Reuse valid existing evidence. More mods are added only to cover a distinct mechanism or concrete reported failure.
5. Under a new game-session GO after the observed blocker/interruption, test the revised framework with the same unchanged LH archive, then the necessary uncovered mechanisms and actual player migration journey. Solo does not establish host/dedicated or multiplayer behavior.

Release remains blocked while the primary journey is unproven. No Workshop item or public binary release should claim drop-in, all-mod compatibility from the current local checks.
