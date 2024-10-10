# tp-hangman-1- Gestion des Utilisateurs

Ce projet est une application web simple développée en Golang, qui met en pratique la gestion des routes, des templates HTML, des formulaires, et des fichiers statiques. 

## Objectif du TP

L'objectif de ce projet est de pratiquer la programmation web en Golang, notamment :
- La gestion des routes
- L'utilisation des templates pour générer des pages HTML dynamiques
- La manipulation de formulaires pour soumettre et traiter des informations utilisateur
- La gestion des fichiers statiques (CSS, images)

## Fonctionnalités du projet

1. **Affichage d'une promotion d'étudiants**  
   Accessible via la route `/promo`, cette page affiche les informations d'une promotion (nom, filière, niveau) ainsi que la liste des étudiants avec leurs informations personnelles (nom, prénom, âge et genre avec une image).

2. **Formulaire pour l'utilisateur**  
   Accessible via la route `/user/form`, cette page présente un formulaire permettant aux utilisateurs de soumettre leurs informations personnelles (nom, prénom, date de naissance, sexe).

3. **Traitement des données utilisateur**  
   Une fois le formulaire soumis, les informations sont validées et stockées dans une structure globale. L'utilisateur est ensuite redirigé vers la page `/user/display` où ses informations sont affichées.

4. **Affichage d'une liste statique d'utilisateurs**  
   Accessible via `/user/list`, cette page présente une liste statique d'utilisateurs (nom, prénom, date de naissance, sexe).

5. **Menu de navigation**  
   Un menu de navigation présent sur toutes les pages permet de naviguer facilement entre les différentes sections de l'application.
## Structure du projet

Voici la structure des fichiers de l'application :


## Installation

### Prérequis

Avant de lancer le projet, assurez-vous d'avoir installé :
- [Golang](https://golang.org/dl/) version 1.16 ou supérieure
- [Git](https://git-scm.com/) pour cloner le projet

### Instructions

1. **Cloner le projet**  
   Ouvre un terminal et exécute la commande suivante pour cloner le projet depuis GitHub :
   ```bash
   git clone 
   cd tp_golang_web

