# Community.4 — dotted version suffix compatibility

**Later acceptance, 1 October:** normal Steam startup and three updated Java consumers were observed manually. See [current status](CommunityReleaseStatus.md). The freeze-time descriptions below retain their historical scope.

On 2026-10-01, a real Steam launch loaded and displayed Community.3 from the local distribution and verified its signature. Living Hordes Java was rejected before loading: its unchanged `zbVersionMin=2.3.3` / `zbVersionMax=2.3.3` did not contain the incorrectly interpreted framework version. No world was created. UI control was interrupted with Escape while PZ was still open; the operator stopped and marked recovery required. No subsequent game acceptance is claimed.

## Root cause and correction

`Loader.loadMods` calls `JavaModInfo.parse`/`parseMerged`, which calls `isVersionInRange`. Its numeric comparator split the complete version at every dot, then trimmed nonnumeric suffixes per component. Thus `2.3.3-community.3` became numeric `2.3.3.3` and exceeded the consumer's maximum. Existing `-beta` coverage had no extra dot, and the prior isolated LH load called the loader after metadata selection; neither exercised this failure.

The comparator now removes the whole prerelease/build suffix before splitting the numeric core. `2.3.3-community.4` compares as API 2.3.3; true core changes remain bounded. Distribution text and manifest retain Community.4. `getVersion` is unchanged. Update ordering retains its separate prerelease comparator; automatic updates remain disabled. No PZ API, consumer metadata, approval rule or persistence schema changes.

A declared public/protected surface audit of 185 upstream classes also found the two original `Patch_loadMods2` ArrayList entry/exit descriptors missing after the List port. Community.4 restores unannotated forwarding overloads. The actual engine hook retains only the List advice; a binary probe compiled against original 2.3.3 now calls the legacy descriptors as well.

Two new regression methods reproduce the comparator and actual metadata-parser failure before the correction. They also cover merged metadata, reverse comparison, build metadata and out-of-range core versions. `ConsumerMetadataSmoke` checks the manifest-derived version against actual frozen LH files in an isolated JVM. Source proof and these local checks do not replace an authorized game loop, save/reload or multiplayer observation.

Local results: **160 Java tests passed** (115 unit, 30 patched, 15 vanilla); all five isolated native/JVM/Kahlua checks passed. The final packaged JAR rejects the unchanged LH metadata with Community.3 and accepts it with Community.4. The English catalogue passes under a French JVM default. Valid-signature, tampered-JAR and wrong-key checks passed. Final JAR SHA-256: `14050862a63b55252816aba580dc4206ef540e83b494b16eae15112f62ec7445`.

At artifact freeze, Community.4 was prepared locally. After the Testeur completed restoration and closed the session, the owner installed the same 11 framework files under the authorized reversible migration, preserving normal JSON/profile and previous backups. No new game launch was made. Actual new game acceptance remains pending, as does the [primary transparent-replacement requirement](DropInCompatibility.md). Never overwrite the frozen Community.3 archive.

The final declared public/protected signature audit reports **zero missing members across 185 upstream classes**. This is a static surface check, not proof of every consumer behavior. The original binary caller probe also runs through the restored advice descriptors, with one explicit legacy call plus three engine List calls in Community; the original baseline still misses those three engine calls.
