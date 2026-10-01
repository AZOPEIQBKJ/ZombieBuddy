# Priority one: transparent replacement for existing ZombieBuddy mods

The maintainer explicitly confirmed this product requirement on 2026-10-01: players should subscribe to the community distribution, enable it for their save and continue using existing ZombieBuddy-dependent mods, without work by those mods' authors. This is the primary release gate. A successful local manual migration or a Living Hordes-only fix does not satisfy it.

## Compatibility contract

- Preserve the technical `ZombieBuddy` Mod ID, upstream Java packages/classes/method descriptors, annotations and Lua-facing contracts used by existing 2.x consumers.
- Preserve existing `require=\\ZombieBuddy`, version limits, JARs, signatures, approval records and saves. Do not edit/rebuild consumer mods or require individual adapters to make the fork work.
- Handle framework compatibility regressions inside Community. Living Hordes is one unchanged representative consumer, not the product boundary.
- Treat existing Workshop dependencies on the original item as a migration requirement to solve centrally. Do not make third-party maintainers retarget their pages or republish their mods as the accepted solution.
- Establish a deterministic single active framework when Steam still has the original item installed. A shared Mod ID alone does not prove which physical distribution is selected.
- Validate the intended subscription/activation path from a real original installation. A manually copied agent, custom launch parameter or private QA configuration is laboratory preparation, not proof of that player journey.
- Do not promise every mod works merely because its framework dependency loads. Define representative coverage of the public framework mechanisms and record concrete remaining incompatibilities; unrelated changes in the game or consumer code must remain visible rather than be silently waived.

## Source findings that currently block the desired journey

1. The installed Windows native agent loads `ZombieBuddy.jar` from the game working directory before normal mod loading. Its updater applies only an already staged adjacent `.new` file; it does not discover a newly subscribed Community item. Source: `c/windows/zbNative.c`, `Agent_OnLoad`, `check_and_apply_update`.
2. Original 2.3.3's `SelfUpdater.checkAndUpdateIfNewer` accepts a replacement only if its JAR entries have an X.509 signature matching Zed's hard-coded certificate fingerprint. Our independent `.zbs`/Ed25519 signature is not that certificate. We do not possess or impersonate Zed's signing identity. Source: upstream commit `0ddf161c27848f12d09e74de7fadbea9d50e621d`, `java/.../SelfUpdater.java`.
3. On audited B42.21, original 2.3.3's `loadMods(ArrayList)` advice misses the game's `loadMods(List)` method. This can prevent reaching normal Java-mod discovery and its updater at all. Existing isolated native tests reproduce this boundary.
4. Default mod search order is `workshop,steam,mods`. Local Community selection in the current lab required explicit folder priority; two subscribed same-ID distributions have not been accepted as an automatic migration path. Source: installed B42.21 `MainScreenState.main`, `ZomboidFileSystem.getAllModFolders` and `setModIdToDir` bytecode.
5. The reviewed vanilla Lua file-writing routes target the Lua cache or a mod's common folder and reject relative traversal; no supported Workshop-only replacement of the preloaded root agent has been established. Do not invent an engine API, exploit a path restriction or disable signature checks to claim this gate passed.

These findings block the existing update/bootstrap route; they are not a proof about every possible future design. Automatic installation/update is an unresolved engineering requirement. A one-time external migration tool or cooperation from the original maintainer would change the player journey and must be explicitly accepted before replacing the stated objective. The current manual-update preview remains a development artifact.

## Ordered next work and acceptance

1. Restore the interrupted QA session after the player closes PZ. Keep its blocked result and intact consumer evidence.
2. Finish the framework-only version-filter correction and package a new immutable candidate. Reproduce old rejection/new acceptance against unchanged consumer metadata; compare upstream callable signatures and retain existing binary caller tests.
3. Resolve and document a verified bootstrap/distribution selection route for the subscription-only goal. Record any irreducible user action before implementation; do not hide it in installation notes.
4. Build a small coverage matrix of unchanged consumers using distinct framework mechanisms: dependency/version parsing, legacy compiled loader calls, annotated patches, Lua exposure, approvals/preloads and save/reload. Reuse valid existing evidence. More mods are added only to cover a distinct mechanism or concrete reported failure.
5. Under a new game-session GO after the observed blocker/interruption, test the revised framework with the same unchanged LH archive, then the necessary uncovered mechanisms and actual player migration journey. Solo does not establish host/dedicated or multiplayer behavior.

Release remains blocked while the primary journey is unproven. No Workshop item or public binary release should claim drop-in, all-mod compatibility from the current local checks.
