# Comment fonctionne le code de DunGo

Ce document s'adresse à quelqu'un qui doit **lire, corriger ou enrichir** DunGo.
Pour jouer, voir le [README](../README.md).

```bash
go run ./src     # lancer le jeu
go vet ./src     # vérifier le code
```

Tout le code est dans **un seul paquet Go** (`package main`, dossier `src/`) et
n'utilise que la bibliothèque standard.

---

## 1. Les fichiers

| Fichier | Contenu |
|---|---|
| `main.go` | `main()` (écran titre), nouvelle partie, prologue, menu du camp, crédits |
| `character.go` | Le héros (`Character`), les classes, les difficultés, la création, la fiche, les niveaux, la mort |
| `items.go` | Les objets (`items`) et équipements (`gear`), le sac, utiliser et équiper un objet |
| `camp.go` | Le marchand, le forgeron et les missions |
| `dungeon.go` | Les étages, les salles, le sphinx, l'écran de victoire |
| `monster.go` | Le bestiaire et `newMonster` |
| `combat.go` | La boucle de combat, les sorts, les dégâts, les attaques des boss, les effets |
| `save.go` | Sauvegarde et chargement en JSON |
| `utils.go` | Couleurs, saisie, menus, cadres, jauges |
| `ascii.go` | Uniquement les dessins, en braille |

---

## 2. Le principe : une pile d'écrans

Chaque écran est une fonction qui boucle : **efface → affiche → lit un numéro →
agit**. Elle ne rend la main que lorsque le joueur tape `0` (Retour).

```
main()                      écran titre
 ├─ newGame / continueGame
 │   └─ campMenu()          LE CAMP  ← on y revient toujours
 │       ├─ exploreDungeon()      choix de l'étage
 │       │   └─ exploreFloor()    suite de salles
 │       │       ├─ riddleEvent() le sphinx
 │       │       └─ roomFight()   un monstre → fight()
 │       ├─ merchant()            acheter / vendre
 │       ├─ blacksmith()          les étals de la forge
 │       ├─ missionBoard()        les missions
 │       ├─ trainingFight()       combat sans risque
 │       ├─ accessInventory()     le sac
 │       ├─ displayInfo()         la fiche du héros
 │       └─ saveGame()
 └─ credits()
```

Les menus sont écrits simplement : une suite d'`option(...)`, puis un `switch`
sur le numéro lu par `readChoice`.

Le héros est passé partout par pointeur (`c *Character`) : toutes les fonctions
modifient le même héros.

---

## 3. Les données dans des tables

Le contenu du jeu est rangé dans des **tables** ; le code ne fait que les
parcourir. Ajouter du contenu, c'est ajouter une ligne.

| Table | Fichier | Une ligne = |
|---|---|---|
| `classes` | character.go | une classe : PV, mana, sort, gains par niveau |
| `difficulties` | character.go | une difficulté : `Power` et `Reward` en % |
| `bestiary` | monster.go | un monstre : PV, attaque, butin, effet, dessin, réplique |
| `items` | items.go | un objet : prix, utilisable en combat ou non |
| `gear` | items.go | un équipement : emplacement, bonus, classe de l'arme |
| `spells` | combat.go | un sort : coût, description, dégâts, dessin |
| `shelves` | camp.go | un étal de la forge et ses recettes |
| `quests` | camp.go | une mission : cible, nombre, récompense |
| `floors` | dungeon.go | un étage : salles, monstres, boss, couleur |
| `riddles` | dungeon.go | une énigme et ses 3 réponses |

---

## 4. Le héros et la sauvegarde

`Character` est **la seule structure sauvegardée**. Une sauvegarde est ce
`Character` converti en JSON (`json.MarshalIndent`) dans `dungo_save_N.json`.
Au chargement, `json.Unmarshal` fait l'inverse.

- Un champ ajouté à `Character` est sauvegardé automatiquement.
- Un champ retiré du code est ignoré au chargement.
- Les champs doivent commencer par une majuscule pour que `encoding/json` les voie.

---

## 5. Le combat (`combat.go`)

```
fight(héros, monstre, entraînement)
 └─ pour chaque tour :
     ├─ showFight        l'écran de combat
     ├─ characterTurn    Attaquer / Sorts / Objets / Fuir
     ├─ monsterTurn      si le monstre est encore debout
     ├─ updateEffects    saignement, brûlure
     └─ un PV à 0 ?      winFight ou loseFight
```

- **Coup critique** : 10 % de chance, dégâts × 2 (`damageMonster`).
- **Attaques des monstres** : `monsterTurn` choisit selon `m.Kind`.
  `basicTurn` pour les monstres ordinaires (× 2 tous les 3 tours),
  `goblinKingTurn`, `lichTurn` et `dragonTurn` pour les boss.
- **Sorts** : `castSpell` retire le mana, applique les dégâts, puis un
  `switch` sur le nom du sort ajoute l'effet (brûlure, soin, Peau de Pierre).
- **Effets** : des compteurs dans `Character` (`Bleeding`, `Burning`…),
  remis à zéro au début de chaque combat.
- **Entraînement** : aucun gain, et `trainingFight` rend au héros ses PV et
  son mana d'avant le combat.

---

## 6. Ajouter du contenu

- **Un monstre** : une ligne dans `bestiary`, puis sa clé dans les `Monsters`
  d'un étage.
- **Un objet** : une constante et une ligne dans `items`, l'ajouter à
  `shopItems` pour le vendre, et un `case` dans `useItem` pour son effet.
- **Une arme ou une armure** : une ligne dans `gear`, puis une recette dans un
  étal de `shelves`.
- **Une mission** : une ligne dans `quests`. `Target` doit être le **nom exact**
  du monstre.
- **Un sort** : une constante dans `character.go`, une ligne dans `spells`, et
  un `case` dans `castSpell` s'il a un effet en plus des dégâts.

---

## 7. L'affichage (`utils.go`)

- Les couleurs sont des **codes ANSI** (`Red = "\033[38;5;196m"`). Toute couleur
  ouverte est refermée par `Reset`.
- Les dessins sont en **braille Unicode** : `printArt(dessin, couleurs...)`
  leur applique un dégradé du haut vers le bas.
- Toute la saisie passe par `readChoice(min, max)`, qui redemande tant que la
  réponse n'est pas un nombre valide : le jeu ne peut pas planter sur une
  mauvaise saisie.
