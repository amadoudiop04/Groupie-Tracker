
Pour valider le module groupie tracker, vous devez valider les 3 projets dans un groupe de 
3 ou 4 personnes.
L’objectif de ce projet est de concevoir une plateforme web proposant un accès à vos 3 
projets en lien avec la musique. Les 3 projets seront jouables en multijoueur.
La plateforme devra fonctionner de la manière suivante : 
Les utilisateurs auront la possibilité de créer un compte. Lors de l'inscription, ils devront 
fournir les informations suivantes :
- Pseudo
- Adresse e-mail
- Mot de passe
Le mot de passe devra être confirmé dans un champ dédié. Les adresses e-mail et les 
pseudonymes devront être uniques.
Une fois inscrit, l'utilisateur pourra se connecter en utilisant soit son pseudo, soit son 
adresse e-mail, accompagné de son mot de passe.
Les mots de passe devront suivre les recommandations de la CNIL.
Un scoreboard permettra d’afficher les pseudos et les scores cumulés des utilisateurs
dans la partie. 
Les morceaux et autres informations concernant la musique ne sont pas stockés par le 
site directement, mais utilisent des APIs externes.
Une landing page présentera les différents projets. 
Le site doit aussi être esthétique et consistant dans sa présentation.
Une fois connecté, l'utilisateur pourra choisir le jeu auquel il souhaite jouer. 
Les utilisateurs devront pouvoir créer des salles de jeu et y inviter leurs amis.
L'entièreté de la plateforme devra être codée en golang et HTML/CSS
JavaScript autorisé uniquement pour les animations et la gestion du temps réel
Les données devront être stockées dans une base de données SQLITE en suivant le 
modèle fourni
Les jeux sont les suivants :

#Guess the song

Ce jeu vous affiche un extrait des paroles d’une chanson, vous devez ensuite retrouver le 
titre de la musique originelle. 
La répartition des points à la fin de chaque manche doit se faire de manière décroissante 
en fonction de la rapidité de la chaque réponse donnée.
Les points seront doublés à chaque tour.
Le tour se termine quand tous les participants ont trouvé le bon titre ou quand le temps de 
réponse tombe à 0.
Une fois que le nombre de tours atteint le nombre de tour max définit lors de la création de 
la partie, la partie se termine et le jeu donne un classement des joueurs en fonction des 
points qu’ils ont accumulé durant la partie


#Blind test

Dans ce jeu, une musique est jouée au début de la manche, et il faut être le plus rapide 
possible pour trouver la chanson correspondante.
La musique joué doit être aléatoire tout comme le moment de ma musique.
Lors de la création de la partie, le créateur doit pouvoir configurer la durée minimale des 
extraits audios.
Les extraits audios doivent être jouable avec des modes tel que :
- Audios ralenti
- Audio accéléré
- Audio inversé
- Audios normal.
La répartition des points à la fin de chaque manche doit se faire de manière décroissante 
en fonction de la rapidité de la chaque réponse donnée.
Une fois que le nombre de tours atteint le nombre de tour max définit lors de la création de 
la partie, la partie se termine et le jeu donne un classement des joueurs en fonction des 
points qu’ils ont accumulé durant la partie.


#Petit bac
 
Ce jeu est une adaptation du jeu du petit bac et reprend donc ses règles.
Le jeu donne une lettre aléatoire au début de chaque manche et les joueurs doivent 
trouver des mots commençant avec la lettre donnée par le jeu.
Les lettres données en début de manche doivent être unique dans chaque partie.
Les catégories de jeu sont les suivantes : 
- Artiste
- Album
- Groupe de musique
- Instrument de musique
- Featuring
Le tour s’arrête dès lors qu’un joueur à donner un mot pour toutes les catégories ou que le 
temps de réponse est à 0.
Le temps de réponse doit être paramétrable lors de la création de la partie.
Les points donner pour chaque mot à la fin de chaque tour se répartisse de la manière 
suivante :
- Si vous n’avez pas donner de réponse ou une mauvaise réponse -> 0 pts
- Si vous avez donné une réponse valide non- unique -> 1 pts
- Si vous avez donné une réponse valide et unique -> 2 pts
Une fois que le nombre de tours atteint le nombre de tour max définit lors de la création de 
la partie, la partie se termine et le jeu donne un classement des joueurs en fonction des 
points qu’ils ont accumulé durant la partie