# Community changes

## Windows installer 0.1.0-preview.2

- Read running processes directly through Windows Toolhelp APIs. The external `tasklist` call used by preview.1 was denied by the execution environment during CLI verification; that attempt stopped without writing fixture files.
- Keep the same runtime and migration behavior. Add real process enumeration and process-name guard tests. Preserve preview.1 unchanged as an unpublished historical artifact.

## Windows installer 0.1.0-preview.1

- Add a separate offline Community installer for the normal Windows Steam launcher. Embed and verify the immutable Community.4 runtime archive and Community Ed25519 signature; no new runtime build, native DLL, private key or game binary is included.
- Preserve dependent mods, saves, approvals, unrelated JVM arguments and Steam account data. Consolidate the effective ZombieBuddy agent and remove only the exact obsolete Aftermath 42.21 adapter entry. Select local mods before Workshop copies, with that global duplicate-ID effect displayed in the preview.
- Add preview, rechecked preconditions, a process-lifetime installer lock, atomic writes, hashed backups, interruption recovery and rollback that retains unrelated Steam edits. The Steam backup contains only the PZ option edit.
- Exercise fresh/migrated installations, repeated runs, failed writes, changed files, recovery, ambiguous launch settings and Windows junction rejection locally. Real B42 filesystem/metadata code selects Community with the installed order while retaining an unchanged consumer, using a simulated installed-item list.
- No actual installer-driven game launch or multiplayer acceptance is claimed. The existing runtime-only archives remain immutable.

## 2.3.3-community.4

- Fix the compatibility comparator so a dotted distribution suffix such as `-community.3` is not read as a fourth numeric API component. The actual solo test of community.3 rejected Living Hordes pinned to API 2.3.3 before its Java code could load.
- Preserve real core-version limits, displayed distribution identity, signature checks and approval policy. No consumer dependency relaxation or adapter is required.
- Restore both original public ArrayList advice descriptors as unannotated forwarding overloads. Only the List overloads are instrumented; existing compiled callers keep linking without a duplicate engine hook.
- Add regressions for the exact min/max pin through both metadata parsing routes and the packaged JAR with the frozen consumer metadata. See [community.4 evidence](doc/CommunityRelease4Evidence.md).
- Community.3 remains a recorded failed candidate; do not publish its draft. Community.4 requires a new authorized game session after the interrupted session is recovered.

## 2.3.3-community.3

- Make the distribution English-only as requested by the maintainer. Remove French Java/Lua catalogues and bilingual publication text. English approval text remains structured and available before game initialization.
- Remove locale switching and child-JVM language propagation; an old `zombiebuddy.language` property has no effect. Console approval retains the original `y/n` contract.
- Keep the integrated 42.21 fixes, original community artwork, manual updates and signature/approval policy.
- Rebuild, independently verify and sign a new candidate. Earlier local archives remain immutable; their unpublished draft is superseded.
- See [community.3 evidence](doc/CommunityRelease3Evidence.md) for local results and pending gameplay acceptance.

## 2.3.3-community.2

- Replace all distributed upstream icons with original community artwork. Remove unused installation screenshots and original graphics source files from the branch tip.
- Share a structured UTF-8 English/French catalogue across Swing, ImGui (embedded/standalone), TinyFD and console approval interfaces. Resolve language from the JVM locale or `-Dzombiebuddy.language=fr|en`; pass the resolved language to the Swing child JVM.
- Translate signature notices while retaining raw diagnostic fields, denial rules and the request/response protocol. French console accepts `o`/`oui` as well as existing `y` answers; English retains `y/n`. EOF/cancellation and invalid signatures do not grant approval.
- Preserve all community.1 loading fixes and manual updates. New signed candidate, distinct version and hashes; community.1 remains immutable.
- Local Java suite: 115 unit + 30 patched + 15 vanilla = 160 passing tests. Real game-loop/Host/dedicated acceptance remains pending. See [community.2 evidence](doc/CommunityRelease2Evidence.md).


## 2.3.3-community.1

- Integrate the B42.21 `loadMods(List)` entry/exit advice and in-place ordering. Preserve the public `loadMods(ArrayList)` binary descriptor.
- Parse Windows `patches_jar` entries at the last colon, preserving drive letters. An external patch is not required for 42.21 loading.
- Construct HTTP clients lazily inside existing IOException handling. Keep failed Workshop checks unknown and signatures unverifiable when their key cannot be obtained; no automatic approval fallback.
- Handle root-adjacent mod/cache paths without a null filename crash.
- Keep the main-menu installation notification when the agent is absent, while preventing options from dereferencing a nil `ZombieBuddy`. Community notification and options are translated into French and English.
- Disable automatic JAR replacement and deferred `.new` writes in both Java update entry points. Keep mod signatures and the signed upstream registry intact.
- Preserve user Steam options in the upstream installer source, including quoted paths. Reject ambiguous wrappers rather than overwrite them. That installer is not part of this preview package.
- Add the Gradle wrapper, remove the author's keychain defaults, make archives reproducible and retain dependency licences and notices inside the JAR.

The code is based on upstream commit `0ddf161c27848f12d09e74de7fadbea9d50e621d`. Existing mods still require their own 42.21 engine compatibility. No game, Host, Linux dedicated-server or real multiplayer acceptance is claimed. The original Java approval dialogs retain their upstream English text; later candidates supersede its publication preparation and language decisions. Workshop publication and subscriptions are unchanged.
