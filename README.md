# Groupie-Tracker - Plateforme de jeux en ligne

Groupie Tracker est un projet mêlant golang, HTML et CSS, ayant pour objectif de développer un site internet permettant de jouer à trois jeux en ligne : un blindtest, un guess the song et un petit bac.

## Description

### Utilisation

Une fois inscrit et connecté, l'utilisateur peut choisir le jeu auquel il souhaite participer. Pour chaque jeu, il peut jouer sur le serveur public, ou bien créer un serveur privé pour jouer avec ses amis. Dans le second cas, le créateur aura la possibilité de modifier certains paramètres de la partie.

### Jeux

#### Guess the song

Le but du guess the song est de trouver le titre d'une musique à partir de ses paroles. 
Les points sont attribués en fonction de la rapidité des joueurs : plus vous répondez rapidement par rapport aux autres joueurs, plus vous obtenez de points ! Les points attribués sont doublés à chaque tour de jeu.
Un tour de jeu prend fin lorsque tous les joueurs ont trouvé le titre de la chanson, ou bien lorsque le compte à rebours se termine. Si un joueur ne suggère aucun titre, aucun point ne lui est attribué.
Entre chaque tour, le classement apparaît pour permetttre aux joueurs de connaître leur position actuelle.
La partie prend fin lorsque le nombre de tour atteint le nombre définit à la création de la partie. Le classement final apparaît pour déterminer le vainqueur de la partie.

#### Blindtest

Le but du blindtest est similaire à celui du guess the song, cette fois cependant, il faut retrouver le titre de la musique à partir du son : les paroles ne sont pas affichées ! 
La musique choisie ainsi que l'extrait de celle-ci sont aléatoires. La durée minimale de cet extrait est configurée à la création de la partie. 
Les extraits audios peuvent également être configurés à la création de la partie : ralentis, accélérés, inversé ou normal, à vous de choisir !
Les points sont attribués en fonction de la rapidité des joueurs : plus vous répondez rapidement par rapport aux autres joueurs, plus vous obtenez de points. 
Un tour de jeu prend fin lorsque tous les joueurs ont trouvé le titre de la chanson, ou bien lorsque le compte à rebours se termine. Si un joueur ne suggère aucun titre, aucun point ne lui est attribué.
Entre chaque tour, le classement apparaît pour permetttre aux joueurs de connaître leur position actuelle.
La partie prend fin lorsque le nombre de tour atteint le nombre définit à la création de la partie. Le classement final apparaît pour déterminer le vainqueur de la partie.

#### Petit bac

Le but du petit bac est de trouver un mot pour chaque catégorie avec la lettre imposé !
À chaque tour, une lettre aléatoire est imposé par le jeu. Les joueurs doivent trouver un mot commençant par cette lettre pour chaque catégorie imposée. Le même lettre ne peut pas être imposée deux fois lors d'une même partie.
Les catégories de mots sont imposés et non modifiables, les voici : 
- Artiste
- Album
- Groupe de musique
- Instrument de musique
- Featuring

Des points sont attribués pour chaque catégorie :
- Si vous n’avez pas donné de réponse ou une mauvaise réponse, vous n'obtenez aucun point.
- Si vous avez donné une réponse valide et qu'un autre joueur a donné la même réponse, vous obtenez 1 point.
- Si vous avez donné une réponse valide et qu'aucun autre joueur n'a donné la même réponse, vous obtenez 2 points.

Un tour de jeu prend fin lorsqu'un joueur à trouver un mot pour chaque catégorie, ou bien lorsque le compte à rebours se termine. Le temps qui est donné aux joueurs pour répondre peut être configuré à la création de la partie.
Entre chaque tour, le classement apparaît pour permettre aux joueurs de connaître leur position actuelle.
La partie prend fin lorsque le nombre de tour atteint le nombre définit à la création de la partie. Le classement final apparaît pour déterminer le vainqueur de la partie.

### Fonctionnalités diverses

En cas d'oubli de mot de passe, l'utilisateur peut recevoir un mail pour le réinitialiser.

## Démarrage

### Prérequis

- Visual Studio Code 
- go 1.22.0

### Installation et Exécution

À partir d'un terminal de commande Linux :
- Éxécutez la commande `git clone https://github.com/amadoudiop04/Groupie-Tracker`
- Éxécutez la commande `cd Groupie-Tracker`
- Éxécutez la commande `go run server.go`
- Une fenêtre devrait apparaître en bas à droite du logiciel. Cliquez sur le bouton `Open in Browser`

## Auteurs

Projet réalisé par Flandrin Hugo, Diop Amadou et Sghaier Yassine dans le cadre du module Groupie Tracker à Ynov.