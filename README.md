# Objectifs

L'objectif de cet exercice est de coder un serveur d'un petit jeu.
Le but du jeu est de trouver une cible et de l'éliminer.

## Règle du jeu

Le jeu se déroule sur une carte carré de 21 cases de coté.

La cible est placé aléatoirement en début de jeu sur une position.

Le joueur est placé au milieu de la carte en début de partie.

Le joueur peut voir la cible si elle se situe à 2 cases de lui.

Le joueur ne peut pas sortir de la carte.

Le joueur peut effectuer les actions suivantes :

* se déplacer vers le haut
* se déplacer vers le bas
* se déplacer vers la gauche
* se déplacer vers la droite
* tirer sur la cible

La cible doit être touchée trois fois pour être éliminer

## Routes du serveurs

**Request**

```json
  POST /move
  {
    "direction": "up|down|left|right"
  }
```

**Response**

```json
  {
    "position": {
      "x" => 1,
      "y" => 7
    },
    "target": {
      "x" => 2,
      "y" => 8
    } || null
  }
```

**Request**

```json
  POST /shoot
  {
    "x": 2,
    "y": 4
  }
```

**Response**

```json
  {
    "result": "touch|miss|kill"
  }
```

## Bonus

**Request**

```json
  GET /map
```

**Response**

Une représentation graphique de la carte à l'instant en asci art :-)

## Livrable attendu

Une PR vers le depot actuel contenant le code et les tests associés.

Il est important que le site fonctionne comme indiqué dans cette documentation mais le respect des bonnes pratiques est aussi important.

De même que la documentation et le nom des commits.

N'oubliez pas que vous devrez expliquer et justifier vos choix techniques lors de votre:w
 entretien.
