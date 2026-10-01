# Community Windows installer

## Player procedure

1. Extract the installer ZIP and run `ZombieBuddyCommunityInstaller.exe`. Choose **Install**. No Python, Java installation or runtime download is needed by the installer.
2. Confirm the detected PZ installation and Steam account. If several accounts or game libraries are found, select the intended one. The normal profile is detected, including an existing Steam `-cachedir` option.
3. Read the preview, close the game and exit Steam normally, then confirm installation. Community's local copy will take priority over Workshop copies. This also applies to other duplicate local Mod IDs.
4. Keep the receipt path displayed on completion. Start Steam normally and enable ZombieBuddy Community for the intended save. Dependent mods and their Workshop dependencies stay unchanged.

To undo the installation, close the game and Steam, run the same installer, choose **Roll back**, and provide its `receipt.json`. If an installation was interrupted, use that receipt to recover before retrying. New edits to managed files or launch options are preserved by refusing an unsafe rollback. Backups stay in `<game>/.zbc-installations/`; do not delete them while rollback may be needed.

Read-only inspection is available with `--plan`. Explicit paths can be supplied with `--game`, `--profile` and `--steam-config`; `--help` lists the options. The installer never closes applications for the player and never starts a game.

This is a development preview, pending actual game acceptance. Use the same installer again after an update only when its documented engine build matches. Automatic updates are not enabled.

## Scope and verified contracts

The maintainer accepted a one-time Community installation/migration on 2026-10-01. This installer targets the normal Windows Steam launcher for the audited PZ 42.21.0 build. It installs the frozen Community.4 runtime; it does not change dependent mods, subscriptions, approval records or saves. The initial installer is separate from the runtime archive, with its own version and source provenance.

- VERIFIED in the installed `ProjectZomboid64.json`: `vmArgs` and per-Windows-version `windows.*.vmArgs` supply JVM arguments. Unknown JSON fields and unrelated arguments must be preserved.
- VERIFIED in upstream installer source and the local Steam configuration: app 108600 uses `UserLocalConfigStore/Software/Valve/Steam/Apps/108600/LaunchOptions`. Only that value may change; Steam must be stopped before writing. Rollback must preserve other Steam changes made since installation.
- VERIFIED in PZ 42.21 `MainScreenState.main`: `-modfolders` consumes the following comma-separated token after the launcher `--` separator. `ZomboidFileSystem.getAllModFolders` uses that order; `setModIdToDir` uses `putIfAbsent`.
- VERIFIED from the manual preview.2 launch: the normal Windows launcher fails before Java starts with `Could not find agent library instrument ... Can't find dependent libraries`. The bundled `java.instrument` module and DLL are present. Isolated `java.exe` checks had hidden this native dependency-search difference.
- VERIFIED in `c/windows/zbNative.c` and `tools/test_native_loader.py`: the native bootstrap sets the bundled JRE DLL search directory before loading `instrument.dll`. A local DLL-only probe reproduces error 126 without that path and succeeds with it, without creating a JVM or launching PZ. Preview.3 uses exactly one `-agentlib:zbNative`, retaining the Community.4 JAR and all consumer contracts.
- VERIFIED in the frozen release manifest: the package targets game SHA-256 `e1a69eb743ede60b213a0fe7f8b83d4fcab773036d256cc4543a336f3b058a33`. Other engine builds require a new audit before enabling installation.
- VERIFIED in the Go 1.25.5 Windows syscall wrappers and [Microsoft's process-snapshot documentation](https://learn.microsoft.com/en-us/windows/win32/toolhelp/taking-a-snapshot-and-viewing-processes): Toolhelp enumeration provides executable names without opening processes or changing their state. The installer only reads this snapshot and closes its handle.
- PENDING: actual normal Steam launch after this installer, Steam UI selection with original dependencies downloaded, save/reload and real multiplayer. Local fixtures and engine-bytecode checks do not establish these observations.

## Installation design

Ship an offline, self-contained Windows executable with the exact reviewed runtime ZIP embedded and its SHA-256 pinned at build time. Verify the manifest and the JAR's Ed25519 signature using the pinned Community public key before planning any writes. No private key or game binary is embedded.

The tool detects Steam libraries and asks for paths/account selection when ambiguous. It previews the chosen game, profile and Steam account. It replaces the framework JAR and signature in the game directory, installs the framework's local mod files and the pinned `zbNative.dll`, configures one native bootstrap agent and removes duplicate ZombieBuddy launch entries. Existing identical native DLLs are retained byte-for-byte; an unknown DLL is refused. Retain existing agent options except the exact obsolete Aftermath 42.21 adapter entry; conflicting options require review. Repairing preview.2 preserves its options while replacing the failed bare Java-agent route.

The Windows bootstrap is unchanged MIT-licensed upstream code, credited to Andrey "Zed" Zaikin. The installer embeds the audited upstream binary SHA-256 `c2ae9335e717ee24b2f4a40d1a3bf77f1519762a72a0459e766a2bbafc077f6c`; this installer build does not rebuild it or claim a reproducible source-to-binary match. Its source, export definition and Makefile accompany the installer ZIP. `c/windows/zbNative.c` is unchanged from upstream `0ddf161c27848f12d09e74de7fadbea9d50e621d` (file SHA-256 `de63e8ce5c211c563a50147886d8559c5685e91ef9b65cdc8aaef55cbf26260b`). The framework and gameplay fixes remain inside Community; this bootstrap is not a consumer-mod adapter. No game/JRE DLL is distributed.

Local `mods` takes priority over Workshop sources so the installed Community mod wins over the still-downloaded original. **This also gives other local duplicate Mod IDs priority.** Display that effect in the installation preview and retain the relative order of the other sources. Do not claim that a same-ID Workshop subscription alone selects Community. Refuse ambiguous additional local ZombieBuddy copies or unexpected files in the managed mod directory.

Before applying, require Steam/PZ/Java to be stopped, reject pending `.new` replacements, unexpected paths/reparse points and changed files, and save the managed before-state plus per-file hashes in a game-local transaction directory. A Windows process-lifetime lock prevents concurrent installers on the same game directory. Write files atomically. A failed/interrupted installation leaves a journal and verified backups; rollback restores only managed targets and refuses to overwrite unexpected edits. The Steam backup contains only the PZ edit, never a copy of the account configuration; rollback restores only the selected app's LaunchOptions if unrelated Steam state has changed.

Installation never starts the game or enables mods in an existing save. The player enables `ZombieBuddy` through the normal UI. The installer can be rerun to verify an unchanged installation; updates remain manual. Alternate launcher batch files, custom Steam command wrappers, non-Windows platforms and hosted-server services need separate support and acceptance.

## Local evidence

- Go 1.25.5 on Windows: the fixture suite covers installation and recovery, plus `go vet`. One symbolic-link creation test is skipped because this account cannot create symbolic links; the real Windows directory-junction refusal test passes. The packaged test log records the exact counts.
- Fixtures cover fresh installation, original migration, the complete frozen runtime payload, repeated installation, every planned write interrupted, rollback interrupted and resumed, changed files/backups, later unrelated Steam changes, missing app/options blocks, ambiguous launch settings, a changed engine after preview and concurrent installers.
- The engine-selection probe uses the audited 42.21 classes, real mod-directory discovery and `ChooseGameInfo.getModDetails`. The installed-item list and mod metadata are fixtures. Default order selects the original framework; local-first order selects Community and retains the consumer. This is not a native Steam, UI, game-loop or server acceptance test.
- Preview.2 was actually applied through its CLI and then failed the normal Steam launch before the menu. A subsequent user-run diagnostic captured the instrumentation DLL error. Those failures remain historical evidence; preview.3 is not thereby game-accepted.
- Preview.3 fixtures cover fresh native-loader deployment/removal on rollback, retaining a matching existing loader, refusing an unknown or changed loader, and preview.2 repair with exact rollback. The DLL-only probe reproduces the dependency-search gap and its correction without a game session. Existing Community.4 native-agent runtime probes remain applicable to the unchanged JAR/bootstrap.
- Installer `0.1.0-preview.1` failed closed when this execution environment denied its external `tasklist` process query; the write-capable CLI changed no fixture files. Preview.2 uses the standard Windows Toolhelp process snapshot API directly, with an actual enumeration test and checks for running Steam/PZ/Java. Preview.1 remains an immutable local artifact and must not be published.

The packaged manifest records the source commit, runtime/native/executable checksums, engine and native-loader reports, and complete installer test log. Build with `tools/build_installer.py --go PATH --runtime ZIP --native-loader PATH`; only the pinned bootstrap is accepted. Actual normal Steam launch of the corrected route, approval behavior, save/reload and multiplayer remain pending. The unchanged runtime retains its separate Community.4 evidence; installer tests do not upgrade that evidence to gameplay acceptance. Earlier runtime-only installation notes suggesting a bare Java agent are superseded for this normal Windows launcher.
