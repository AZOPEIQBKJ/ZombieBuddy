# Publier ZombieBuddy Community

Préparation du 1er octobre 2026. Nom confirmé : **ZombieBuddy Community**, compte **AZOPEIQBKJ**. Le fork existant [AZOPEIQBKJ/ZombieBuddy](https://github.com/AZOPEIQBKJ/ZombieBuddy) est réutilisé ; branche de maintenance `community-42.21`. Les branches historiques ne sont ni remplacées ni fusionnées. Le Mod ID reste `ZombieBuddy` ; un futur item Workshop recevra un nouvel ID, différent de `3619862853`.

## Livrables

- Candidate **2.3.3-community.2**, correctif 42.21 intégré, visuels originaux, dialogues FR/EN structurés, mise à jour manuelle, signatures et approbations conservées.
- **160 tests Java**, cinq contrôles runtime/Kahlua, quatre prévols, signature positive/négative et catalogue du JAR/sous-processus vérifiés. Les tests Go de community.1 restent acquis pour l'installateur inchangé et non distribué.
- [Installation FR](../doc/Installation_FR.md), [EN](../doc/Installation_EN.md), [crédits/licences](../THIRD_PARTY_NOTICES.md), [support](../SUPPORT.md), modèles d'issue/PR, [notes de release](RELEASE_NOTES.md).
- [Messages auteur/modération](CONTACTS.md) rédigés, non envoyés ; [Workshop FR](WORKSHOP_FR.txt) et [EN](WORKSHOP_EN.txt), couverture et dossier privé hors profil actif.
- [publication.json](publication.json) fixe la candidate et ses empreintes ; `publicReleaseReady=false` tant que l'acceptation manque.

## Sources et release GitHub

L'authentification Git existante permet la préparation du dépôt sans créer de jeton ou nouveau compte. Pousser uniquement notre branche vers `origin` sous `community-42.21`, sans force. L'historique upstream, `master`, `patch-1` et `v2x` restent conservés. La branche communautaire doit être la branche par défaut, Issues activé et le signalement privé de vulnérabilité disponible avant diffusion.

La release **draft + pre-release** doit viser le commit exact contenu dans le manifeste du ZIP, avec le tag `v2.3.3-community.2`. Joindre uniquement `ZombieBuddyCommunity-2.3.3-community.2.zip`, `manifest.json` et `SHA256SUMS.txt`. Les archives de sources automatiques ne sont pas le paquet installable. Ne pas diffuser le JAR du jeu, une clé privée, les dossiers de travail ou le kit préparatoire entier. La CI publique ne peut pas être présentée comme ayant testé le moteur sans installation légitime.

Le code source peut être visible avec son statut expérimental. La release reste en brouillon jusqu'à la décision de diffusion et l'acceptation de la portée annoncée. Ne pas déplacer un tag publié ou écraser une archive existante ; une modification du JAR impose nouvelle signature, version et preuves.

## Réception avant publication des binaires

Le banc n'a lancé aucune partie. La demande solo de community.2 reste DRAFT dans la file du workspace Aftermath. Préparer une migration réversible : jeu arrêté, sauvegardes des fichiers et paramètres, un seul agent effectif JSON + Steam/script/service, absence de `.jar.new`, distribution unique, conservation des autres paramètres et approbations. L'autorisation d'une session ne modifie pas automatiquement les abonnements ou réglages globaux.

Mutualiser S01 avec Living Hordes : version Community visible **et** Java LH réellement chargé (`observer Java loaded version=0.36.0-poc`), console sans nouvelle erreur pertinente. Le watermark seul est insuffisant. Réutiliser les aperçus et clics LH déjà observés le 30 septembre ; aucun lancement par preset ou recette. La demande LH et son profil doivent être repris par son propriétaire avant sélection.

Les dialogues modifiés demandent une observation représentative FR/EN, refus/autorisation et persistance au redémarrage, sur fixture sûre. Mutualiser la sauvegarde/reprise avec le S02 LH lorsqu'il est sélectionné ; ne pas transformer cela en campagne de tests répétitive. Le Testeur solo agit après son **GO explicite**. Host, Linux dédié et MP demandent leurs propres preuves avant annonce ; une première preview peut être limitée à Windows solo.

## Auteur, modération et Workshop

La [MIT upstream](https://raw.githubusercontent.com/zed-0xff/ZombieBuddy/v2.3.3/LICENSE.txt) autorise la réutilisation du code couvert avec sa notice. La [politique PZ](https://projectzomboid.com/blog/modding-policy/) traite séparément crédits, droits des éléments redistribués et reprise d'un mod cassé/absent (§8 : contact de la modération). Le silence ne prouve pas l'abandon. Envoyer les messages préparés, conserver date/lien/réponse ; ne pas prétendre avoir épuisé des canaux non essayés.

Après ces démarches et la réception de la candidate :

1. Publier la release GitHub en gardant **pre-release** et ses limites explicites.
2. Copier le seul dossier `workshop-draft` préparé dans le dossier `Zomboid/Workshop` du profil retenu, sous un nom distinct. Il contient `Contents`, `preview.png` et `workshop.txt` privé sans ID. Ne pas copier tout le kit.
3. Dans le publieur Workshop du jeu, sélectionner ce dossier et créer **un nouvel item privé**. Le titulaire du compte accepte les éventuelles conditions Steam. Noter l'ID retourné ; ne jamais utiliser l'ID de Zed.
4. Mettre à jour les liens croisés GitHub/Workshop avec l'ID réel, vérifier la page et une installation depuis ce canal. Le Mod ID reste `ZombieBuddy` ; conserver une seule distribution et un seul agent.
5. Passer l'item public quand cette revue est terminée. Les mods qui imposent l'ancien item doivent changer leur propre dépendance Workshop. L'abonnement seul n'installe toujours pas l'agent Java.

`tools/prepare_publication.py` vérifie ZIP/manifeste/SHA, refuse un dossier existant et n'envoie rien. `tools/ValidateWorkshopDraft.java` contrôle le format avec les classes réelles du jeu, dans un cache isolé, sans boucle de jeu ni upload. Le JAR, les licences et les guides accompagnent les contenus Workshop.

## Maintenance

Trier les issues : chargement/installation, confiance/signature, régressions, confort. Chaque correction reste intégrée au framework avec reproduction, preuve, test ciblé et version nouvelle. Examiner upstream sans fusion automatique de l'alpha 3.x. En cas de régression, identifier la version et le retour arrière ; ne pas remplacer silencieusement un fichier publié.
