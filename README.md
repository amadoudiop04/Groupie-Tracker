# Groupie-Tracker - Plateforme de jeux en ligne

Groupie Tracker est un projet mêlant golang, HTML et CSS, ayant pour objectif de développer un site internet permettant de jouer à trois jeux en ligne, tous liés à l'univers musical : un blindtest, un guess the song et un petit bac.

## Description

### Utilisation

Une fois inscrit et connecté, l'utilisateur peut choisir le jeu auquel il souhaite participer. Pour chaque jeu, il peut jouer sur le serveur public, ou bien créer un serveur privé pour jouer avec ses amis. Dans le second cas, le créateur aura la possibilité de modifier certains paramètres de la partie.

### Jeux

#### Guess the song 📜

Le but du guess the song est de trouver le titre d'une musique à partir de ses paroles. 
La partie prend fin lorsque le nombre de tour atteint le nombre définit à la création de la partie. Le classement final apparaît pour déterminer le vainqueur de la partie. 
Attention ! Plus les joueurs répondent rapidement et plus cela handicape les autres. En effet, lorsqu'un joueur répond avant les autres, il consomme un tour de jeu que les autres ne pourront donc pas avoir. C'est l'opportunité d'obtenir plus de points pour les joueurs les plus rapides !

#### Blindtest 🎧 

Le but du blindtest est similaire à celui du guess the song, cette fois cependant, il faut retrouver le titre de la musique à partir du son : les paroles ne sont pas affichées !  
La musique choisie est aléatoire. La durée minimale de cet extrait est configurée à la création de la partie. 
Un tour de jeu prend fin lorsque le compte à rebours se termine, ou bien lorsque le joueur envoie sa réponse. Si un joueur ne suggère aucun titre, aucun point ne lui est attribué.
La partie prend fin lorsque le nombre de tour atteint le nombre définit à la création de la partie. Le classement final apparaît pour déterminer le vainqueur de la partie.

#### Petit bac ✏️

Le but du petit bac est de trouver un mot pour chaque catégorie avec la lettre imposé !
À chaque tour, une lettre aléatoire est imposé par le jeu. Les joueurs doivent trouver un mot commençant par cette lettre pour chaque catégorie imposée.
Les catégories de mots sont imposés et non modifiables, les voici : 
- Artiste
- Album
- Groupe de musique
- Instrument de musique
- Featuring

Un tour de jeu prend fin lorsqu'un joueur a trouvé un mot pour chaque catégorie, ou bien lorsque le compte à rebours se termine. Le temps qui est donné aux joueurs pour répondre peut être configuré à la création de la partie.
La partie prend fin lorsque le nombre de tour atteint le nombre définit à la création de la partie. Le classement final apparaît pour déterminer le vainqueur de la partie.

### Fonctionnalités diverses

En cas d'oubli de mot de passe, l'utilisateur peut recevoir un mail pour le réinitialiser.  
Les joueurs peuvent personnalisés leur partie lors de la création de leur serveur.  
Les utilisateurs peuvent modifier les informations de leur compte à tout moment depuis l'espace profil, accessible depuis la page d'accueil.

## Démarrage

### Prérequis

- Visual Studio Code 
- go 1.22.0
- ngrok

### Installation et Exécution

Pour jouer seul : 
-----------------
    À partir d'un terminal de commande Linux :
    - Exécutez la commande `git clone https://github.com/amadoudiop04/Groupie-Tracker`
    - Exécutez la commande `cd Groupie-Tracker`
    - Exécutez la commande `go run server.go`
    - Une fenêtre devrait apparaître en bas à droite du logiciel. Cliquez sur le bouton `Open in Browser`


Pour jouer à plusieurs :
------------------------
Commencez par installer ngrok (`https://ngrok.com/download`). Le dossier téléchargé est compressé, vous devez le décompresser.
Créez un compte sur le site de ngrok (`https://dashboard.ngrok.com/signup`). Gardez cette page ouverte.
À partir d'un terminal de commande Linux, vous pouvez maintenant lancer les commandes :
    - Exécutez la commande `git clone https://github.com/amadoudiop04/Groupie-Tracker`
    - Exécutez la commande `cd Groupie-Tracker`
    - Exécutez la commande `go run server.go`
Enfin, ouvrez l'application "ngrok.exe" depuis le dossier décompressé. Sur la page d'accueil de ngrok, vous pourrez trouver une commande commençant par `ngrok config add-authtoken` (il faut être connecté pour y avoir accès). Copiez cette commande et collez la dans le terminal de l'application ngrok.
Vous pouvez maintenant lancer la commande `ngrok.exe http 8080` depuis le terminal ngrok. Le site vient d'être mis à disposition de tous sur internet. Pour y accéder, un lien est mis à votre disposition par ngrok. Suivez ce lien sur votre navigateur, cliquez sur `Visit Site`, et voilà ! Vous pouvez maintenant jouer avec n'importe qui !
Le site peut être supprimé à tout instant par la personne l'ayant lancé. Au prochain lancement, l'hôte aura seulement besoin de lancer le serveur golang avec la commande `go run server.go`, puis dans lancer la commande `ngrok.exe http 8080` dans le terminal ngrok.

## Auteurs

Projet réalisé par Flandrin Hugo, Diop Amadou et Sghaier Yassine dans le cadre du module Groupie Tracker à Ynov Lyon.