# Community.2 assets and localization

Audit: 1 October 2026. The earlier community.1 archive is preserved separately.

| Resource | Community.2 provenance |
| --- | --- |
| `42/icon_128.png`, `42/icon_256.png`, root icons | Exports of the original geometric community SVG in `publishing/assets/`. |
| `java/src/main/resources/zb_icon.png` | Original community icon, integrated before compilation/signing. |
| `common/media/ui/zb_steam_options_*` | Four unused upstream installation screenshots removed after source-reference review. |
| `icon.xcf`, `cmdline.png`, `cmdline.xcf` | Removed from the branch tip; historical installation guide replaced with community redirects. |
| `publishing/assets/workshop-cover.svg` and PNG | Original geometric/typographic cover; no game or upstream artwork sampled. |

The new SVGs and PNG exports are supplied under this repository's MIT licence. Upstream code copyright and Git history remain intact. No endorsement is claimed.

Approval text is structured in `java/src/main/resources/me/zed_0xff/zombie_buddy/i18n/messages_{en,fr}.properties`, read as UTF-8. Swing, embedded/standalone ImGui, TinyFD and console use the shared catalogue. Signature notices have separate presentation strings; raw diagnostic fields and approval/signature decisions remain unchanged. JVM locale selects French or English; `-Dzombiebuddy.language=fr|en` overrides it and propagates to Swing's child JVM. Later game-language changes are outside this mechanism. Native OS button captions and external technical diagnostics retain their own language.

Local verification: 160 Java tests, including locale resolution, catalogue parity/formatting, console `oui`/deny/EOF and invalid signatures; final shaded-JAR UTF-8 and child-process checks pass. See [current evidence](../doc/CommunityRelease2Evidence.md). This is not a claim of observed layout in the actual game: representative approval rendering and persistent decisions still require acceptance.
