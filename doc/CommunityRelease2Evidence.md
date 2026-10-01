# Community.2 release preparation — source contracts

Scope authorized on 1 October 2026: complete independently actionable preparation, then give the user the remaining procedure in chat. Integrate original community artwork and structured FR/EN approval text; preserve signing/approval decisions, API 2.x and Mod ID. No deployment, game session or external message follows from these changes.

Verified before implementation:

- `SwingModApprovalFrontend.runSwingSubprocessBatch` launches `SwingApprovalMain` with the agent JAR, serialized request/response, and validates row ID/hash on return. Propagate one normalized JVM language property before `-cp`; do not alter that protocol.
- ImGui renders through `ImguiApprovalDialog`; the standalone `ImguiApprovalMain` uses the same dialog. TinyFD and console are independent retained frontends. All cancellation/invalid-signature branches must remain deny paths.
- Language resolution will use the standard JVM locale at dialog time, with an explicit `-Dzombiebuddy.language=fr` or `en` override. This avoids invoking game translation classes during early agent startup or from the standalone Swing JVM. It does not claim to follow a later in-game language change. UTF-8 properties are read using `InputStreamReader`, not platform-default encoding. Unsupported explicit languages fall back to English.
- `ZBSVerifier` raw diagnostic strings and constructors have existing unit contracts. Preserve those strings/fields and expose translated UI details separately; upstream/JVM exception details and Steam-supplied ban reasons remain diagnostic data, not translated decisions.
- `zb_steam_options_*` appears only in publication inventory, with no runtime source references after the community.1 notification rewrite. Remove these unused screenshots. `cmdline.png` is referenced only by the historical installation guide; replace that guide with a short redirect to community guides before removing it.
- `zb_icon.png` is loaded by ImGui, while `42/mod.info` references the two versioned PNG icons. Replace with original geometric community artwork and retain paths. No engine API change or world/save mutation is introduced.

Game source baseline remains PZ 42.21.0, SHA-256 `e1a69eb743ede60b213a0fe7f8b83d4fcab773036d256cc4543a336f3b058a33`. Local tests, package/signature evidence and actual gameplay acceptance will be reported separately below when available.

## Local results

- JDK 25 / Gradle 9.3.1: **160 Java tests passed** (115 unit, 30 patched, 15 vanilla), then shaded JAR produced. The 12 added tests cover locale/catalogue/console/signature presentation; existing approval and signature suites continue to pass.
- Five isolated runtime checks passed with actual PZ 42.21 classes: original baseline, integrated candidate, double-agent first-policy retention, existing native loader, and Kahlua. No game loop was started. Public ArrayList compatibility and List handling remain checked.
- Four preflight tests passed. Go installer code is unchanged from the community.1 passing suite and is not shipped.
- The final shaded JAR loaded accented UTF-8 resources and passed normalized French selection to a new English-default JVM. This checks packaged resources and child propagation, not actual Swing/ImGui layout.
- Ed25519 `.zbs` verified independently through the JDK; a tampered JAR and wrong public key were refused. Signing identity: Aftermath Systems, not an X.509 signature from Zed.
- JAR SHA-256: `365744cb6037be4ed5575201fdc398401e3805f1c1cbdb3ad3a585844165b63a`. Final ZIP/source commit are in the release manifest; no game class, native DLL, installer or private key is included.

No new game console exists because no game session was launched. Actual approval rendering, persistence, solo loading and save/reload remain pending; Host, dedicated, Linux/macOS and real multiplayer are not accepted. No game installation, Steam setting or subscription was changed. No save migration or gameplay authority change is introduced. The upstream signed author registry remains a network dependency; all dependent mods still require their own compatibility checks.
