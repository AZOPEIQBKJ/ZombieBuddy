# ZombieBuddy Community 2.3.3-community.4 — development preview

Independent maintenance of ZombieBuddy 2.3.3 by Andrey "Zed" Zaikin. This preview integrates the PZ 42.21 loading fix and corrects dotted Community-version comparisons while retaining original 2.x callable signatures. No per-mod Aftermath compatibility adapter is required. **English only** for the framework interface and publication materials.

**Draft: not yet published or a stable/all-mod compatibility claim.** On 1 October 2026, normal Windows Steam startup and Java loading of Living Hordes 0.36.1, Loadout poc3 and Bulk Storage Native 0.1.8 were observed together with Community.4. Their signatures and Java entry points succeeded. Broader gameplay, save/load, approval-decision persistence, Host/dedicated and other platforms remain separately unaccepted.

## Windows installation

Download `ZombieBuddyCommunityInstaller-0.1.0-preview.3-windows-x64.zip` and follow the accompanying **INSTALLATION.md**. It contains the offline Windows executable, pinned MIT native bootstrap, bootstrap source/notices and its own manifest. No separate Java/Python installation is required. Close PZ and Steam before applying. Keep the receipt for rollback. The executable is not Authenticode signed.

The corrected native route has been observed working after a reversible development repair. The final preview.3 installation transaction remains a distinct pending observation; local transaction/rollback fixtures do not establish a fresh player installation. The installer selects local mods before Workshop copies, also affecting other duplicate local Mod IDs. Subscribing alone does not install the agent; one Community installation is required. Updates are manual.

The separate runtime ZIP is an immutable reference artifact. Its bundled documentation predates the native-bootstrap repair: **use INSTALLATION.md and the current repository guide for the normal Windows Steam route**, not an older bare `-javaagent` recommendation. Runtime-only ZIPs do not contain the native bootstrap DLL.

## Compatibility and evidence

- Preserve `ZombieBuddy` Mod ID, original Java namespaces, mod signatures, approvals and saves. Declared version ranges remain enforced; mods pinned to an excluded exact version may require an author update.
- Owned consumers now use a verified minimum 2.3.3 without an unnecessary maximum. This is not proof of compatibility with all future versions or all third-party mods.
- Runtime: 160 Java tests, five isolated JVM/native/Kahlua checks, signature checks and zero missing declared members across 185 audited upstream classes.
- Installer: 39 passing leaf cases in 21 Go groups, vet/structure and native/engine-selection checks; one unavailable symlink case skipped, real directory-junction refusal passed.

## Artifact identity

Runtime source: `f6cc18d3ff2706bb49d2cebace4bae05da7070af`.
Runtime ZIP SHA-256: `f1b66cf9e11aac8663de92de382944a979cb4f09e651d320d8d3938b471a3c99`.
Installer source: `ad5143190fc8c84593b448ee6456e291d3c85264`.
Installer ZIP SHA-256: `1b7ab4eb6a996231c9018eb3b13d21077f3301a41910a466ed48dee9ae41c96d`.

The release tag identifies the runtime source. `installer-manifest.json` identifies the separately built installer source. Existing immutable archives are unchanged by these later release notes. Use the named download archives, not GitHub's generated source archive.

Community Ed25519 `.zbs` public key: `989ac279f40f1a35fa0616645e319fc44cde9a15a842b7c365870147bfff3ce0`. It is independent of Zed's X.509 signing identity. Signature checks and the signed upstream author registry remain active; do not bypass a failure with allow-all.

Original framework and MIT copyright: Andrey "Zed" Zaikin. Community maintenance: Aftermath Systems / AZOPEIQBKJ. Related proposals: benoitthore #56 and yuruichang #58; the full #58 rewrite is not included. Licences and dependency notices accompany the downloads. No Community Workshop item is published.
