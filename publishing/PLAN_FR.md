# Publier ZombieBuddy Community

Préparé le **1er octobre 2026**, à la demande de l'utilisateur. Nom confirmé : **ZombieBuddy Community**. Compte personnel confirmé : **AZOPEIQBKJ**. Destination prévue : `AZOPEIQBKJ/ZombieBuddyCommunity` ; le dépôt distant et la page Workshop ne sont pas créés par ces préparatifs.

Le correctif reste intégré dans ZombieBuddy. Le projet conserve l'historique upstream, l'API 2.x et le Mod ID `ZombieBuddy`. Il reçoit un **nouvel ID Workshop** : `3619862853` est celui de Zed, utilisé uniquement pour les crédits. Ne pas déclarer l'original comme dépendance Workshop du fork, car cela réintroduirait les deux distributions.

## Ce qui est prêt

| Livrable | Fichier / état |
| --- | --- |
| Page GitHub du fork | [README](../README.md), corrigé pour ne plus envoyer vers l'installateur original |
| Publication GitHub | [Notes FR/EN](RELEASE_NOTES.md), procédure ci-dessous, version et empreintes figées |
| Description Workshop | [Anglais](WORKSHOP_EN.txt) et [français](WORKSHOP_FR.txt), textes BBCode à actualiser avec la réception finale |
| Installation et restauration | [FR](../doc/Installation_FR.md), [EN](../doc/Installation_EN.md), prévol en lecture seule |
| Messages externes | [Auteur et modération](CONTACTS.md), rédigés, non envoyés |
| Crédits et droits | MIT originale, [notices](../THIRD_PARTY_NOTICES.md), [inventaire des visuels et traductions](ASSETS_AND_LOCALIZATION.md) |
| Support | [Guide FR/EN](../SUPPORT.md), modèle d'issue et de PR dans `.github/` |
| Visuel propre au fork | [Couverture SVG](assets/workshop-cover.svg), PNG 256 × 256 et icônes préparés ; pas encore substitués aux visuels du binaire figé |
| Assemblage hors ligne | [Script](../tools/prepare_publication.py), qui vérifie le ZIP et prépare un brouillon Workshop privé sans ID |
| État machine | [publication.json](publication.json), `uploaded=false`, `publicReleaseReady=false` |

Les dons automatiques hérités de `.github/FUNDING.yml` ont été retirés du fork pour éviter une attribution ambiguë du bouton Sponsor. Les crédits et la licence de Zed restent intacts. L'ancienne description `steam.txt` est remplacée par un renvoi vers les textes communautaires. Les tâches Rake et l'installateur historiques ne sont pas la chaîne de publication du fork.

## Ordre recommandé

### 1. Préparer les échanges et ouvrir le dépôt source

Le code couvert est sous [MIT dans la version 2.3.3](https://raw.githubusercontent.com/zed-0xff/ZombieBuddy/v2.3.3/LICENSE.txt), avec conservation du copyright et de la licence. Le silence de l'auteur n'est pas une déclaration d'abandon. Les visuels tiers et les règles Workshop se traitent séparément ; les [messages](CONTACTS.md) présentent un fork indépendant sans présumer l'arrêt du projet.

Une mise en ligne des **sources clairement expérimentales** peut précéder la réception des binaires. Après instruction de publication, créer un fork GitHub de `zed-0xff/ZombieBuddy`, propriétaire `AZOPEIQBKJ`, nom `ZombieBuddyCommunity`. Conserver la filiation ; ne pas réinitialiser l'historique. Vérifier que le nom est disponible au moment de créer. La branche de maintenance recommandée est `community-42.21`, issue de notre branche locale, pas `master` upstream qui correspond à une autre évolution.

Commandes préparées, **non exécutées** :

```powershell
# Depuis le dépôt indépendant, après création du fork distant et revue du commit.
git remote add origin https://github.com/AZOPEIQBKJ/ZombieBuddyCommunity.git
git remote -v
git push origin codex/community-42.21:community-42.21
```

Si `origin` existe déjà, examiner sa destination au lieu de le remplacer. Ne pas pousser vers `upstream`, ni utiliser `--mirror`, `--force` ou `--all`. Sur GitHub, choisir `community-42.21` comme branche par défaut, activer Issues et le signalement privé des vulnérabilités, puis vérifier README, licences et crédits. Mettre à jour les mentions « dépôt prévu » seulement après vérification des liens réels. Le dépôt ne contient ni clés privées, ni JAR du jeu, ni configuration utilisateur. Ne pas envoyer le répertoire `artifacts/` ou le monorepo Aftermath.

Références GitHub : [fork](https://docs.github.com/en/pull-requests/how-tos/work-with-forks/fork-a-repo), [releases et brouillons](https://docs.github.com/en/repositories/releasing-projects-on-github/managing-releases-in-a-repository). La compilation complète exige une installation légitime de PZ ; ne pas prétendre qu'une CI publique valide le moteur si elle n'a fait que des contrôles de texte, ni y téléverser le jeu.

### 2. Fermer le lot technique avant diffusion de binaires

La candidate `2.3.3-community.1` reste inchangée et sert de référence. Le prochain lot est concret : appliquer les nouveaux visuels, retirer/remplacer les captures de l'ancienne installation et externaliser les textes des dialogues Java en FR/EN. [Inventaire et critères](ASSETS_AND_LOCALIZATION.md). Il produira une **nouvelle version**, au minimum `2.3.3-community.2`, avec ses propres signatures et empreintes.

Conserver les mises à jour manuelles et les contrôles de confiance actuels pour la première publication. Ne pas ajouter un auto-updater ou redistribuer l'installateur original pour accélérer la sortie. Le registre d'auteurs signé upstream reste une dépendance réseau : déclarer ce fait, sans le remplacer par une autorisation automatique lors d'une panne.

Recompiler et rejouer les tests touchés, puis la suite requise pour le paquet final. Vérifier l'absence de classes du jeu, clés privées, DLL et exécutable upstream ; vérifier positivement et négativement la signature. Les seules modifications de documentation n'imposent pas de refaire 148 tests. Toute modification du JAR impose de rattacher les preuves au **nouveau JAR**, pas à celui de community.1.

### 3. Recevoir l'installation et le fonctionnement promis

Les 148 tests et JVM isolées de community.1 sont acquis à leur portée. Aucun succès en partie/MP n'est déclaré. La demande solo locale `ZBC-SOLO-233C1-r1` reste DRAFT, dans la file de tests du workspace Aftermath.

Avant session, figer une migration réversible : jeu arrêté, copies des fichiers et paramètres, un seul agent dans JSON + Steam/service/script, absence de `.jar.new`, choix explicite de distribution, conservation des autres options. Ne pas modifier un abonnement Workshop ou les réglages globaux à partir d'un simple GO de test. Le propriétaire LH doit actualiser sa propre demande avec le nouveau framework.

| Observation manquante | Essai utile et limite |
| --- | --- |
| Chargement réel 42.21 | Un seul démarrage partagé avec S01 Living Hordes ; version Community + Java LH réellement chargé + console sans nouvelle erreur pertinente. Un watermark seul ne suffit pas. |
| Décisions d'approbation | Sur une fixture sûre, vérifier l'affichage FR/EN modifié, un refus et une autorisation, puis la persistance au redémarrage. Réutiliser ce redémarrage pour les autres observations compatibles. Ne pas contourner avec `allow-all`. |
| Sauvegarde/reprise | Mutualiser avec S02 LH si sélectionné. Prouver la reprise du consommateur choisi ; ne pas en déduire la compatibilité de tous les mods/sauvegardes. |
| Installation/retour arrière | Vérifier les chemins et paramètres réellement retenus ; constater que restaurer l'ensemble remet l'état précédent. Le retour à 2.3.3 rétablit aussi ses limites 42.21. |
| Host | Seulement avant de revendiquer Host : hôte et client, approbations, Java chargé des deux côtés, reconnexion/reprise pertinente. |
| Linux dédié et client | Seulement avant de revendiquer ce contexte : démarrage sans dialogue bloquant, droits de lecture/écriture, politiques persistantes, connexion/reconnexion, redémarrage. Ni une JVM Windows ni un solo ne couvre cela. |

Réutiliser les aperçus et clics LH observés le 30 septembre. Ne pas refaire les quatre presets, le cas 101 %, une promenade ou un lancement par variante. La nouvelle incertitude concerne le chargement du framework et les écrans modifiés. Le Testeur solo agit après son GO explicite ; MP/dédié se planifient séparément. Une première preview peut être limitée à Windows solo si seule cette portée est reçue ; afficher les autres environnements comme non pris en charge/testés au lieu d'attendre ou de promettre tous les contextes.

### 4. Préparer puis publier la release GitHub

Depuis le commit final propre : construire, signer, vérifier et emballer. Garder versions et SHA cohérents entre JAR, `mod.info`, notes, ZIP, manifeste et demande d'acceptation. Créer un tag annoté `vVERSION` sur **ce commit exact**, puis le pousser explicitement vers `origin`. Ne pas déplacer le tag d'une archive déjà publiée.

Dans GitHub : créer une release **draft**, cocher **pre-release**, choisir le tag vérifié, coller les notes FR/EN actualisées et joindre le ZIP nommé, `manifest.json`, `SHA256SUMS.txt`. Vérifier les liens et le téléchargement, puis publier seulement après la décision de mise en ligne. Ne pas publier community.1 par défaut si community.2 la remplace ; les notes actuelles documentent la référence locale.

### 5. Publier un nouvel item Workshop

Avant cet envoi, consigner la démarche auteur/modération et résoudre les droits des éléments réellement redistribués. La [politique PZ](https://projectzomboid.com/blog/modding-policy/) consultée le 1er octobre 2026 demande notamment crédits et droits (§1/4) ainsi qu'une prise de contact avec la modération pour une reprise de mod cassé/absent (§8). MIT autorise le code couvert ; ce n'est pas une garantie d'acceptation par la plateforme. L'absence de réponse ne devient pas une permission supposée.

Le dossier local est un **brouillon privé**, sans ID, hors du répertoire Workshop actif. Il n'est pas le paquet public final : il reprend community.1 et ses limites. Le futur paquet doit être régénéré depuis la candidate reçue et inclure les licences dans `Contents/mods/ZombieBuddy/`, car Steam envoie `Contents`, pas les notices placées à côté.

Préparation hors ligne, depuis le dépôt du fork, vers un dossier qui n'existe pas encore :

```powershell
python tools/prepare_publication.py --candidate CHEMIN_DU_ZIP_FIGE --output NOUVEAU_DOSSIER_LOCAL
```

Le script exige le ZIP/SHA/version/commit déclarés dans `publication.json`, contrôle tous les fichiers contre le manifeste, refuse de réécrire une préparation et n'effectue aucun envoi. Le premier contrôle de format a été effectué avec les classes réelles PZ 42.21 : `readWorkshopTxt`, visibilité privée (2), ID absent, texte FR conservé, `validateContents` et PNG 256 × 256 acceptés. Ce contrôle n'est ni un upload Steam ni une réception en jeu. Le constructeur PZ exige un chemin dans ses préfixes autorisés : utiliser une copie sous `Workshop/` d'un cache **isolé**, pas le profil utilisateur actif. `tools/ValidateWorkshopDraft.java` ne lance aucune boucle de jeu et n'appelle aucune API d'envoi.

Préparation hors ligne, depuis le dépôt du fork, vers un dossier qui n'existe pas encore :

```powershell
python tools/prepare_publication.py --candidate CHEMIN_DU_ZIP_FIGE --output NOUVEAU_DOSSIER_LOCAL
```

Le script exige le ZIP/SHA/version/commit déclarés dans `publication.json`, contrôle tous les fichiers contre le manifeste, refuse de réécrire une préparation et n'effectue aucun envoi. Le premier contrôle de format a été effectué avec les classes réelles PZ 42.21 : `readWorkshopTxt`, visibilité privée (2), ID absent, texte FR conservé, `validateContents` et PNG 256 × 256 acceptés. Ce contrôle n'est ni un upload Steam ni une réception en jeu. Le constructeur PZ exige un chemin dans ses préfixes autorisés : utiliser une copie sous `Workshop/` d'un cache **isolé**, pas le profil utilisateur actif. `tools/ValidateWorkshopDraft.java` ne lance aucune boucle de jeu et n'appelle aucune API d'envoi.

Après autorisation de l'envoi, utiliser le publieur PZ avec le paquet final, titre distinct, visuel communautaire et description vérifiée. **Créer un nouvel item**, conserver privé pendant la revue de la page, puis noter l'ID renvoyé par Steam et contrôler que ce n'est pas `3619862853`. Ajouter les liens croisés GitHub/Workshop avec cet ID réel et recevoir une installation depuis ce canal avant de passer public. Les conditions Steam applicables sont acceptées par le titulaire du compte. Un upload privé est déjà un envoi externe ; il n'est pas effectué ici.

Le Mod ID reste `ZombieBuddy` ; le nouvel ID Workshop n'est pas interchangeable. Les profils serveurs utilisent le nouvel item Workshop avec ce Mod ID. Les mods qui imposent l'ancien ID Workshop doivent actualiser leur métadonnée/dépendance : pas de désabonnement ni de migration silencieuse des joueurs. Après mise en ligne, ne pas annoncer que « s'abonner suffit » : l'agent requiert toujours son installation.

## Maintenance après publication

Trier les issues par panne de chargement/installation, confiance/signature, régression puis confort. Pour chaque correctif : reproduction, preuve moteur, changement intégré, test ciblé et nouvelle version. Examiner les évolutions upstream sans fusion automatique de son alpha 3.x. Fournir des correctifs upstream quand pertinent. En cas de régression publiée, annoncer la version affectée et le retour arrière ; ne pas remplacer silencieusement un ZIP ou une signature sous le même numéro.

## Prochaine action concrète

Terminer le lot **visuels + dialogues FR/EN** dans une nouvelle candidate, en parallèle de l'envoi manuel des messages préparés. Ensuite, préparer la migration et mutualiser l'acceptation avec LH. Les préparatifs actuels ne lancent ni le jeu, ni un serveur, ni une publication, ni un message externe.
