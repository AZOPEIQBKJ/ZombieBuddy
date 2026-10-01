# Assets, identity and localization before binary publication

Audit: 1 October 2026. The frozen community.1 ZIP is not rewritten by publication preparation.

## Visual resources

| Resource | Provenance / next action |
| --- | --- |
| `42/icon_128.png`, `42/icon_256.png` | Upstream graphics retained in community.1. Replace in the next candidate with the prepared community icons. |
| `java/src/main/resources/zb_icon.png` | Upstream graphic inside the JAR, used by the approval UI. Replace with community artwork before the next build/signature. |
| `common/media/ui/zb_steam_options_*` (four PNGs) | Upstream installation screenshots; contain original instructions. Audit actual Lua references and replace with newly captured/redacted community instructions or remove only if unused. Do not infer reuse rights for third-party UI merely from the root MIT code licence. |
| Root `icon_*.png`, `icon.xcf`, `cmdline.*` | Historical source assets; provenance beyond the upstream repository has not been independently established. Replace/remove from the public branch tip where unnecessary; preserve attribution and history. Keeping history is not a rights clearance. |
| `publishing/assets/community-mark.svg` and generated PNGs | New geometric vector design prepared for this fork; no upstream/game artwork sampled. Use as the new icon. |
| `publishing/assets/workshop-cover.svg` and PNG | New typographic/geometric cover; explicitly identifies the community distribution. Original SVG retained for resizing. |

New SVG designs and their generated PNG exports are supplied under the repository's MIT licence; original copyright notices remain intact. The wordmark identifies an independent community fork, not an official edition. Final brand review should consider confusion with upstream; do not claim endorsement.

Runtime replacements require a new candidate and rechecking asset references and UI layout. A new Workshop cover beside unchanged Contents does not prove that the JAR's icon or old installation screenshots were replaced.

## FR/EN gaps found in Java source

Lua options and the modified missing-agent notice have FR/EN resources. The following inherited paths still contain literal player-facing English:

| Path under `java/src/main/java/me/zed_0xff/zombie_buddy/` | Texts to extract and translate |
| --- | --- |
| `frontend/SwingApprovalMain.java` | Window title, introduction, columns, unsigned/unknown status, trust explanations, persistent decisions, Cancel/Yes/No and tooltips |
| `frontend/ImguiApprovalMain.java` | Standalone approval title and rendered approval controls/status strings; inspect helper renderers too |
| `frontend/TinyfdModApprovalFrontend.java` | Prompt, invalid-signature denial, trust warning, dates/statuses |
| `frontend/ConsoleModApprovalFrontend.java` | Allow/deny prompts, persistence prompt, unknown fields and accepted answer guidance |
| `ZBSVerifier.java` and approval helpers | User-visible signature/trust notices propagated into those frontends; distinguish stable diagnostic logs from UI copy |

Implement a structured FR/EN catalogue used by all retained frontends. Establish locale resolution from actual engine/JVM and subprocess behavior before coding: approval may occur before the game translations exist, and Swing/ImGui are separate entry points. Do not guess that a Lua translation function is available from early Java bootstrap. Carry the chosen locale to child processes explicitly if required; retain an English fallback for unavailable keys/locales.

Acceptance criteria: matching FR/EN keys and formatting arguments; localized safety/consent semantics unchanged; console EOF still denies; invalid signature still denies; cancellation never grants trust; persisted flags and response protocol unchanged. Test all variants locally through their shared catalogue/protocol; real UI review should cover representative rendered mechanisms without requiring a game launch per string.

No new gameplay API, trust policy or updater is required for this work. There is no complete FR/EN claim for community.1.
