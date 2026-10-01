# Community changes

## 2.3.3-community.1

- Integrate the B42.21 `loadMods(List)` entry/exit advice and in-place ordering. Preserve the public `loadMods(ArrayList)` binary descriptor.
- Parse Windows `patches_jar` entries at the last colon, preserving drive letters. An external patch is not required for 42.21 loading.
- Construct HTTP clients lazily inside existing IOException handling. Keep failed Workshop checks unknown and signatures unverifiable when their key cannot be obtained; no automatic approval fallback.
- Handle root-adjacent mod/cache paths without a null filename crash.
- Keep the main-menu installation notification when the agent is absent, while preventing options from dereferencing a nil `ZombieBuddy`. Community notification and options are translated into French and English.
- Disable automatic JAR replacement and deferred `.new` writes in both Java update entry points. Keep mod signatures and the signed upstream registry intact.
- Preserve user Steam options in the upstream installer source, including quoted paths. Reject ambiguous wrappers rather than overwrite them. That installer is not part of this preview package.
- Add the Gradle wrapper, remove the author's keychain defaults, make archives reproducible and retain dependency licences and notices inside the JAR.

The code is based on upstream commit `0ddf161c27848f12d09e74de7fadbea9d50e621d`. Existing mods still require their own 42.21 engine compatibility. No game, Host, Linux dedicated-server or real multiplayer acceptance is claimed. The original Java approval dialogs retain their upstream English text; a complete FR/EN interface audit and independent public branding remain prerequisites for public distribution. Workshop publication and subscriptions are unchanged.

Cette preview intègre les corrections dans ZombieBuddy. La validation locale ne remplace pas les essais en partie et sur serveur. Les dialogues Java hérités restent en anglais ; leur traduction complète et la revue de publication restent à traiter avant diffusion publique.
