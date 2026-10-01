# Artwork and English interface

The maintainer requested English-only distribution on 1 October 2026. Community.3 replaces the bilingual candidate; previous archives remain historical.

All shipped mod icons and the JAR approval icon use original community SVG exports. Unused upstream installation screenshots and original graphics source files were removed from the branch tip. The Workshop cover is original geometric/typographic artwork. No upstream or game artwork was sampled. New SVGs and PNG exports use the repository's MIT licence; upstream code attribution and history remain intact.

Approval text is structured in `java/src/main/resources/me/zed_0xff/zombie_buddy/i18n/messages_en.properties`. Swing, embedded/standalone ImGui, TinyFD and console share it. Only English Java and Lua catalogues are distributed. The Java interface stays English regardless of JVM locale or a stale community.2 language override. Native OS button captions and external diagnostics may retain their system language.

No trust decision, signed-author registry or signature policy is relaxed. See [current evidence](../doc/CommunityRelease3Evidence.md). Representative actual-game rendering and persistent decisions remain pending; local catalogue tests are not a claim of observed game layout.
