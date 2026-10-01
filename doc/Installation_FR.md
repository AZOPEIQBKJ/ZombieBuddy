# Installation de la preview ZombieBuddy Community

Version 2.3.3-community.1, basée sur ZombieBuddy 2.3.3 de Zed. Cible auditée : PZ 42.21.0. Cette candidate a des contrôles locaux ; la réception en partie et sur dédié reste à effectuer. Elle contient le correctif de chargement : aucun JAR de compatibilité Aftermath n'est nécessaire.

## Vérifier avant toute copie

Fermer normalement le jeu et, pour un serveur, arrêter celui-ci. Conserver une sauvegarde des fichiers que l'on remplace et des options de lancement. Vérifier les empreintes du paquet avec `manifest.json` ; la signature `.jar.zbs` identifie Aftermath Systems. Ce fichier n'est pas une signature X.509 de Zed.

Lancer `python tools/preflight.py --game "CHEMIN_DU_JEU"` pour lire le JSON de lancement. Ce contrôle ne modifie aucun fichier et reste incomplet tant que les options Steam ou du script/service serveur n'ont pas été examinées. Après examen, fournir chaque argument d'agent supplémentaire avec `--extra-jvm-arg="ARGUMENT"` et ajouter `--extra-options-reviewed`. Ne pas utiliser ce drapeau pour masquer une source de paramètres inconnue.

Il doit rester exactement un agent ZombieBuddy. Retirer de la configuration de lancement choisie les anciens doublons et les références à `AftermathLHCompat4221.jar` lors de la migration. Conserver les paramètres personnels (`-cachedir`, mémoire, debug, autres agents). Une entrée native déjà présente dans le JSON compte même si Steam n'affiche rien. Aucun script de ce paquet ne réécrit Steam ou les fichiers vanilla.

S'il existe `ZombieBuddy.jar.new`, interrompre la migration et examiner cette mise à jour en attente : l'ancien lanceur natif peut l'appliquer avant le démarrage de Java. Ne pas la laisser remplacer silencieusement la candidate.

## Installer après revue du prévol

1. Sauvegarder hors du répertoire actif l'ancien `ZombieBuddy.jar`, ses éventuelles signatures et la copie de mod que l'on va remplacer. Noter leurs emplacements et empreintes pour le retour arrière. Conserver les approbations dans `.zombie_buddy`.
2. Copier le dossier `Contents/mods/ZombieBuddy` du paquet dans le dossier local de mods sélectionné. Activer une seule distribution portant le Mod ID `ZombieBuddy`. L'original Workshop et le fork ne doivent pas être activés ensemble ; ce choix n'est pas automatisé. Aucune nouvelle page Workshop n'est publiée avec cette preview.
3. Copier le JAR et sa signature depuis `Contents/mods/ZombieBuddy/libs/` vers le répertoire où l'unique agent les charge. Pour une nouvelle installation utilisant l'agent Java standard, les options Steam sont `-javaagent:ZombieBuddy.jar --`, à combiner avec les options conservées après revue. Le JAR doit se trouver dans le répertoire de travail du lanceur, ou son chemin doit être absolu.
4. Si l'installation possède déjà **un seul** `-agentlib:zbNative` opérationnel, son remplacement de JAR peut conserver cette entrée, après vérification de l'absence de `.new` et de doublons. Le paquet ne fournit pas de nouvelle DLL. Cette route exige sa propre réception ; les essais JVM seuls ne valident pas l'installation réelle.
5. Au prochain essai autorisé, relever `ZombieBuddy Community v2.3.3-community.1` et confirmer qu'un mod Java attendu est effectivement chargé. Un affichage de version seul ne suffit pas. Conserver les décisions d'approbation normales ; ne pas utiliser `allow-all` pour contourner un problème de signature.

## Serveur dédié Linux

Le serveur et ses clients doivent avoir la même candidate et des mods compatibles 42.21. Pour une JVM lancée directement, `-javaagent:/chemin/absolu/ZombieBuddy.jar` est un argument JVM placé avant la classe principale ; le séparateur `--` des options Steam n'est pas un argument universel à ajouter à Java. Examiner le vrai script ou service de l'hébergeur avant modification. Le processus doit pouvoir lire le JAR et écrire son répertoire de configuration.

En environnement sans interface, utiliser une politique connue et des approbations préparées explicitement. `policy=deny-new` bloque les JAR non approuvés ; il ne les rend pas fonctionnels. La disponibilité d'une entrée console et la persistance des approbations doivent être réceptionnées sur le dédié. Aucun succès Linux, Host ou MP réel n'est annoncé pour cette candidate.

## Mise à jour et retour arrière

Les mises à jour automatiques du JAR sont désactivées pour cette preview. Une nouvelle version se remplace explicitement, jeu/serveur arrêtés, avec vérification de ses empreintes et conservation de la précédente. Les contrôles de signature des mods et le registre d'auteurs upstream restent actifs.

Pour revenir en arrière, arrêter le processus, restaurer ensemble le JAR, la copie de mod et les paramètres sauvegardés. Ne pas supprimer les approbations ni modifier les sauvegardes de campagne. Le retour à l'original 2.3.3 remet aussi ses limites de chargement 42.21 : il ne constitue pas une réparation du jeu. La migration n'autorise aucune conversion de sauvegarde.
