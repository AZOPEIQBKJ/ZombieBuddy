# Third party notices

ZombieBuddy was created by Andrey "Zed" Zaikin and is distributed under the MIT licence in [LICENSE.txt](LICENSE.txt). The community fork preserves that notice and the upstream Git history. Upstream code retains its original attribution. Community.2 replaces the shipped runtime icons with original community artwork and removes unused upstream installation screenshots; the audit is in publishing/ASSETS_AND_LOCALIZATION.md. This distribution is independently maintained by Aftermath Systems.

The Java agent embeds these pinned dependencies. Exact dependency JAR hashes are in [licenses/dependencies.json](licenses/dependencies.json); their licence and notice texts are included in `licenses/` and in the agent under `META-INF/community/licenses/`.

| Component | Version | Licence and notices |
| --- | --- | --- |
| Byte Buddy | 1.18.8 | Apache 2.0; its ASM notice is also included |
| Byte Buddy Agent | 1.18.8 | Apache 2.0 |
| ClassGraph | 4.8.184 | MIT |
| Bouncy Castle bcprov-jdk18on | 1.78.1 | Bouncy Castle permissive licence |
| Gson | 2.12.1 | Apache 2.0 |

Additional licence sources: [Bouncy Castle](https://github.com/bcgit/bc-java/blob/r1rv78/LICENSE.html), [Gson 2.12.1](https://github.com/google/gson/blob/gson-parent-2.12.1/LICENSE). Gradle wrapper code retains its Apache licence header. Build and test dependencies are downloaded by Gradle/Go; they are not separately shipped in the preview package.

Compatibility work was informed by upstream issue reports and the public proposals of benoitthore (PR #56) and yuruichang (PR #58). The preview contains focused changes; it does not incorporate the full PR #58 rewrite. Original reports and the local before/after evidence are linked from [CommunityEvidence](doc/CommunityEvidence.md).

No Project Zomboid code, game JAR, extracted engine class, private key or user configuration is included in the release archive. Game code is only read locally for compilation and contract tests. The original native loader and Windows installer are not included in this Java-agent preview package.
