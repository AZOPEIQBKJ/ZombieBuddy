# Community.4 release status — 2 October 2026

The current development candidate is **ZombieBuddy Community 2.3.3-community.4**, with the separate Windows installer **0.1.0-preview.3**. The normal Steam startup and simultaneous Java loading of three updated consumers were observed manually on 1 October. The binary release is still a draft; no Community Workshop item has been published.

Workflow update, 2 October: the user now performs game tests manually. The former tester task and queue are retired; historical reports remain evidence. Complete source/local/log checks first, then propose only short manual steps for distinct remaining uncertainties. Do not repeat accepted loading solely for a new build, operator or report, or recreate session reservations and READY/GO submissions.

## Observed on the installed PZ 42.21.0 build

The user's screenshot shows Community.4 and three active Java mods. The final console reports valid signatures and successful Java entry points for:

| Consumer | Observed version | Declared framework range |
| --- | --- | --- |
| Aftermath Systems: Living Hordes | 0.36.1-poc | 2.3.3 minimum, no maximum |
| Aftermath Systems: Loadout | 0.1.0-poc3 | 2.3.3 minimum, no maximum |
| Aftermath Systems: Bulk Storage Native | 0.1.8-poc | 2.3.3 minimum, no maximum |

Steam records the normal launcher at 23:50:23 Europe/Brussels, PID 28896, with exit code 0 at 23:51:43. The console reaches `onGameInitComplete`. Its lines 130/131, 240/241 and 316/317 identify the consumers and successful entry points. The designated solo tester inspected the screenshot, verified the seven archived evidence hashes and accepted reuse of these loading observations. No repeated menu/bootstrap test is required solely for documentation.

These are updated consumers, not proof that every unchanged mod works. Third-party mods that explicitly exclude the current framework version may require an author update; their declared bounds remain enforced. The owned modules now use a reviewed minimum without an unnecessary exact maximum. Game/API compatibility remains a separate requirement.

## Local checks and immutable artifacts

- Runtime: 160 Java tests, five isolated JVM/native/Kahlua checks, signature checks and a declared API audit with zero missing members across 185 upstream classes.
- Installer: 21 Go groups, 39 passing leaf cases, vet and structure checks. One unavailable symbolic-link test was skipped; the directory-junction refusal case passed. Native bootstrap and engine-selection checks passed locally.
- Runtime source: `f6cc18d3ff2706bb49d2cebace4bae05da7070af`.
- Runtime ZIP SHA-256: `f1b66cf9e11aac8663de92de382944a979cb4f09e651d320d8d3938b471a3c99`.
- Runtime JAR SHA-256: `14050862a63b55252816aba580dc4206ef540e83b494b16eae15112f62ec7445`.
- Installer source: `ad5143190fc8c84593b448ee6456e291d3c85264`.
- Installer ZIP SHA-256: `1b7ab4eb6a996231c9018eb3b13d21077f3301a41910a466ed48dee9ae41c96d`.
- Installer EXE SHA-256: `484276a68193ef8fc83f563cac4ca875312265bfe1d253f170a8b7485bdf4d2c`.

The installer includes the pinned, unmodified MIT native bootstrap and its source/notice. This bootstrap is required by the normal Windows launcher; it is not a per-consumer compatibility adapter. Updates remain manual.

## Remaining observations

The installed native route was repaired reversibly during development and then observed working. That establishes the runtime path, not an actual preview.3 installation transaction. The final EXE has a read-only, zero-change plan on this installation. An unchanged-installation check must not be described as a fresh migration.

Still distinct: the final installer interaction/transaction, representative approval UI and decision persistence, relevant consumer gameplay/save-load mechanisms, Host/dedicated multiplayer, and non-Windows platforms. Earlier Shift-triggered approval and log writes are accepted at their recorded scope, but do not prove a controlled cross-restart persistence scenario. Local fixtures are not gameplay acceptance.

Historical Community.3 and installer preview.2 failures remain recorded. The old C4/r3 request freezes LH 0.36.0 and was not executed with that artifact; do not substitute 0.36.1 hashes or replay the old build just to fill a report. Future testing should cover only the remaining distinct mechanisms and preserve current saves, approvals, profiles and module versions.

Internal evidence references: `docs/research/ZOMBIEBUDDY_MANUAL_THREE_CONSUMERS_2026-10-01.md` and `docs/testing/manual-evidence/zbc-three-consumers-20261001/RESULT.md` in the development workspace. Raw user logs and profile data are not release assets.
