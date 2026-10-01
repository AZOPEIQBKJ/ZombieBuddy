# Community fork evidence

User approval: 2026-10-01. Scope: integrate compatibility fixes in ZombieBuddy itself; no Aftermath compatibility agent. Local branch: `codex/community-42.21`, based on upstream `v2.3.3` / `0ddf161c27848f12d09e74de7fadbea9d50e621d`. Retain upstream packages and the `ZombieBuddy` Mod ID for binary compatibility. This preview is a distinct distribution, not an official Zed release.

## Verified contracts before implementation

- PZ 42.21.0 `ZomboidFileSystem.loadMods(java.util.List)` verified in the installed game bytecode, SHA-256 `e1a69eb743ede60b213a0fe7f8b83d4fcab773036d256cc4543a336f3b058a33`.
- Upstream `Loader.loadMods(ArrayList)` and `Patch_ZomboidFileSystem` enter/exit advice verified in pinned sources. Keep the public ArrayList overload while adding List handling. Preserve in-place reordering and approvals.
- `Agent.premain` ignores duplicate agents before processing their options. Keep that guard; the integrated fix must work on the first invocation without `patches_jar`.
- Windows `patches_jar` paths split at the drive colon. Baseline run actually reports `JAR not found: D`: 28 of 30 patched tests fail. This affects existing API users as well as the test harness.
- Three eager `HttpClient` initializers can throw outside callers' existing IOException handling. Moving creation into those handlers must preserve unknown/unverifiable results, not authorize a mod.
- `JavaModInfo` checks parent existence but not the null filename of a drive root. Test exact path shapes before/after fixing.
- `ZombieBuddy_Options.applySettings` dereferences a missing global; main-menu installation notification already exists. Guarding options must preserve that notification and valid-agent behavior.
- `SelfUpdater` trusts upstream's signing certificate and can replace the fork. Preview releases use manual updates; all automatic file-replacement entry points must be disabled, while signature-verification helpers and mod approvals remain intact.

## Baseline results

The unmodified stable sources compile on JDK 25.0.1 against the audited PZ JAR. Gradle 9.3.1 distribution SHA-256: `b266d5ff6b90eada6dc3b20cb090e3731302e553a27c5d3e4df1f0d76beaff06`.

- Upstream vanilla tests: 15/15 pass.
- Upstream patched tests: 2/30 pass; 28 fail after the Windows path parsing error.
- Upstream unit tests: 96/97 pass. The remaining test assumes a workshop.txt on the author's personal machine; replace it with a temporary fixture.
- Upstream Go installer tests pass on Go 1.25.5.

No game process launched. Previous real 42.21 Living Hordes loading remains BLOCKED; JVM tests do not establish solo, Host or dedicated-server acceptance. Existing Living Hordes UI/preset observations must not be repeated just because the framework changes.

Upstream reports: [#53](https://github.com/zed-0xff/ZombieBuddy/issues/53), [#46](https://github.com/zed-0xff/ZombieBuddy/issues/46), [#42](https://github.com/zed-0xff/ZombieBuddy/issues/42), [#29](https://github.com/zed-0xff/ZombieBuddy/issues/29), [#18](https://github.com/zed-0xff/ZombieBuddy/issues/18). PRs [#56](https://github.com/zed-0xff/ZombieBuddy/pull/56) and [#58](https://github.com/zed-0xff/ZombieBuddy/pull/58) were consulted; their broader lifecycle/performance changes are outside this preview.

## Candidate results on 1 October 2026

Final agent SHA-256: `a0711c36597795348fa94d31d400e34c04eaa263a4e0a3bb3fb3d6898395eb6a`. Version: `2.3.3-community.1`. JDK and the game's JRE are 25.0.1. Dependency notices are included inside this JAR.

- Gradle: **103 unit + 30 patched + 15 vanilla tests pass**, 148 total. Existing valid/invalid/tampered mod-signature cases still pass. HTTP construction failures preserve unknown/unverifiable states. Root-path fixtures and manual-update guards pass.
- Go installer suite passes, including quoted-path VDF round trips, preserved user arguments, idempotence and refusal of ambiguous launch wrappers. No installer was run on Steam configuration.
- Real PZ 42.21 bytecode, isolated JVMs: baseline intercepts zero of three native List calls; the community Java agent intercepts all three. A caller compiled against unmodified 2.3.3 still links and executes. Reordering preserves the remaining IDs for ArrayList and LinkedList.
- The same calls pass through the installed upstream Windows native loader, SHA-256 `c2ae9335e717ee24b2f4a40d1a3bf77f1519762a72a0459e766a2bbafc077f6c`, matching the v2.3.3 GitHub release asset. Only a temporary JAR copy is used; the installed JAR is unchanged.
- Duplicate-agent regression: a second `policy=allow-all` declaration does not change the first agent's effective `deny-new` policy. No compatibility adapter is supplied in any community scenario.
- Real Kahlua: missing-agent option application does not crash; all four existing settings still reach their callbacks. Missing-agent notification appears once; a loaded agent does not show it. UI objects are fixtures, not a rendered game observation. FR/EN translation key sets match.
- Four read-only preflight regressions pass, including the actual JSON-plus-Steam double-agent shape.
- JAR signed with the existing Aftermath Systems Ed25519 identity, SteamID64 `76561198061324182`. Both the existing signing tool and a standalone JDK verifier accept it. Altered-JAR and wrong-public-key checks fail as required. This is a `.zbs` signature, not an X.509 JAR signature.

The runtime report and signature results travel with the preview. The test harness initially hit a sandbox restriction reading JDK security configuration; execution with the game JRE in an isolated directory resolved that environment issue. Unit tests also initially mixed the game's bundled Bouncy Castle classes with the standalone test dependency; moving the engine-dependent reorder check into the runtime harness resolved the fixture conflict without changing production signature enforcement.

Still **not observed**: fresh game launch, enabled-mod resolution between original and fork, actual Java mod loading through the menu, save/reload, Linux dedicated server, Host and MP. The framework changes do not add gameplay mutations or networking; each dependent mod retains its own authority/persistence acceptance requirements.
