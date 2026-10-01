# ZombieBuddy Community installer

Offline Windows installer for the normal Steam launcher. Uses only Go's standard library. See [player instructions and verified scope](../doc/CommunityInstaller.md).

The frozen `2.3.3-community.4` runtime ZIP is embedded at build time; it is not committed here. Stage the verified ZIP with `tools/build_installer.py --prepare-only` before running `go test ./...` or `go vet ./...`. Use Go 1.25.5. Windows junction and process-lock tests require Windows; the non-Windows implementation refuses installation.

Keep `installer/` as the upstream reference, with its own module/tests. Do not build or ship that original installer under the Community name. No native DLL is needed by this installer.

`--plan` is read-only. `--install` or the interactive Install choice shows the plan and asks for confirmation; `--yes` supplies explicit command-line confirmation without bypassing integrity or process checks. `--rollback RECEIPT` reverses a recorded transaction with the same checks. Normal interactive launches offer Install, Roll back and Quit.

Runtime archives and installer artifacts are independently versioned. Do not overwrite either after packaging. Actual Steam launch and game acceptance remain separate from fixture tests.
