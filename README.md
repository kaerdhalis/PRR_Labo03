# PRR - Laboratoire 3 : Election

## Lancement et parametrage de l'application

Le parametreage des applications se fait dans le fichier json **config.json** du package config,
il permet de determiner le nombre de processus, leur adresses Ip ainsi que les delais artificiel.
Le champs ArtificialDelay determine le temps en seconde que chaque processus va passer avant de redemander une election.
Le lancement de l'application se fait au moyen de ligne de commande en passant le numero du processus en parametre
ex: **go run processus.go id** (avec les id allant de 0 a n-1)
Nous avons pris la decision de limiter le parametrage au fichier json uniquement et de ne pas accepter les parametres en ligne de commande car nous estimons que passer des adresses en parametre est tres fastidieux et peut engendrer rapidement des errreurs.

## Implementation

###Processus
Module s'occupant de la partie tache applicative du processus, elle se charge de lancer l'administrateur et va periodiquement lui demander la valeur de l'elu pour relancer une election si  l'elu ne repond pas

###Mutex
Implemente l'algorithme de Chang et Roberts pour gerer l'election du processus ayant la meilleure aptitude
La boucle principale est basée sur un select qui attend les signaux de la tache applicative ou les messages des autres process.


###Connection
Package gerant la communication UDp(envoie et reception des messages) et transmet les messages reçus au gestionnaire,s'occuppe egalement de l'envoie et la reception des ack pour savoir quel processus est en panne.

###Config
Package gerant la configuration generale de l'application(nombre de process, leurs adresses et le temps avant de relancer une election) en lisant le fichier config.json

## Ce qui reste a faire
L'application fonctionne,on peut lancer plusieurs processus et ceux-ci vont elire celui ayant la meilleure aptitude de plus si le processus elu crashe une nouvelle election est bien lancée et lorsque qu'un processus etant tombe en panne se reconnecte il reintegre bien l'anneau et peu etre réélu.
Cependant nous avons un problème au niveau du reseau, en effet les processus recoivent regulierement des messages vides alors que aucun autre processus n'a envoyé de message(meme dans le cas d'un processus unique) ce qui peut entrainer des elections perpetuelles,notre solution est de verifier que la liste contenau dans les messages n'est pas vide mais nous n'avons pas trouvé la cause de ces messages phantomes.
