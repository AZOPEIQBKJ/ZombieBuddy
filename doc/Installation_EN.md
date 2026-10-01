# Installing the ZombieBuddy Community preview

Version 2.3.3-community.4 is based on Zed's ZombieBuddy 2.3.3. The audited target is PZ 42.21.0. It fixes the dotted-version compatibility rejection observed in community.3. Local checks do not establish gameplay or dedicated-server acceptance. The loading fix is integrated; no Aftermath compatibility JAR is required.

## Preflight

Close the game normally, or stop the server. Back up files and launch settings before replacing anything. Check package hashes against `manifest.json`. The `.jar.zbs` signature identifies Aftermath Systems; it is not Zed's X.509 signature.

Run `python tools/preflight.py --game "GAME_PATH"`. It reads the launcher JSON without modifying it. Also inspect Steam and service/script options. Provide each additional agent argument using `--extra-jvm-arg="ARGUMENT"`; add `--extra-options-reviewed` only after checking every launch source. The JSON alone is insufficient.

Keep exactly one ZombieBuddy agent. Remove duplicate entries and references to `AftermathLHCompat4221.jar` when migrating, preserving user arguments, other Java agents and approval records. A native agent already in the JSON counts even when Steam options look empty. No script in this package rewrites Steam or vanilla files.

Stop if `ZombieBuddy.jar.new` exists: the old native loader can replace the JAR before Java starts. Resolve that pending update explicitly first.

## Files and launching

1. Back up the previous JAR, signatures, mod directory and launch settings outside their active locations. Record original paths and hashes. Retain `.zombie_buddy` approval records.
2. Copy `Contents/mods/ZombieBuddy` into the selected local mods directory. Enable only one distribution with Mod ID `ZombieBuddy`. This preview has no published Workshop item. If the original item remains subscribed, B42.21's default `workshop,steam,mods` search order can select it before the local fork. For the reviewed Steam route, add `-modfolders mods,workshop,steam` after the launcher `--` separator (for example `-- -debug -modfolders mods,workshop,steam` when retaining debug). This prioritizes all local duplicate IDs; review that effect and verify the selected Community path in the log. Subscriptions are not changed automatically.
3. Copy the JAR and sidecar from the mod's `libs/` into the location used by the single agent. A new standard Java-agent Steam launch uses `-javaagent:ZombieBuddy.jar --`, combined with existing options after review. The JAR must be in the launcher's working directory, or use an absolute path.
4. An installation already using one working `-agentlib:zbNative` can retain that entry while replacing its JAR, once duplicates and pending `.new` updates are ruled out. This package does not supply a new native DLL. That route still needs acceptance in the actual installation.
5. At the next authorized game test, verify the community version and an actually loaded Java mod. A version watermark alone is not sufficient. Keep normal approval/signature checks; do not bypass a signature problem with `allow-all`.

## Interface language

The community interface and documentation are English-only. No `zombiebuddy.language` option is required; an option retained from the unpublished community.2 candidate has no effect. Native OS controls and external technical diagnostics may use their own system language.

## Linux dedicated server

Use the same candidate and compatible 42.21 mods on clients and server. In a direct Java invocation, `-javaagent:/absolute/path/ZombieBuddy.jar` is a JVM argument before the main class. Steam's `--` separator must not be blindly inserted into a direct Java command. Inspect the host's actual service/script first. The process needs read access to the JAR and write access to its configuration directory.

Use explicit approvals and a suitable headless frontend. `policy=deny-new` blocks unapproved JARs; it does not approve them. Console availability, persistent decisions and restart behavior require dedicated-server acceptance. Linux, Host and real multiplayer have not been accepted for this preview.

## Updates and rollback

Automatic framework replacement is disabled in this preview. Replace versions manually with the process stopped, verify hashes and retain the previous version. Mod-signature checks and the signed upstream author registry remain active.

To roll back, stop the process and restore the JAR, mod directory and launch settings together. Keep approval records and campaign saves intact. Returning to original 2.3.3 restores its 42.21 limitations as well. No save conversion is part of this migration.
