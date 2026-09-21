# Comment fonctionne le code de DunGo

Ce document s'adresse à quelqu'un qui doit **lire, corriger ou enrichir** DunGo. C'est un résumé : le [guide du code](DunGo-guide-complet.md) ([PDF](DunGo-guide-complet.pdf)) explique tout pas à pas, avec des rappels de Go. Pour jouer, voir le [README](../README.md).

```bash
go run ./src     # lancer le jeu
go vet ./src     # vérifier le code
```

Tout le code est dans **un seul paquet Go** (`package main`, dossier `src/`) et n'utilise que la bibliothèque standard.

---

## 1. Les fichiers et les étiquettes

Un fichier par thème. Chaque fonction porte un commentaire **sur sa ligne**, qui commence par l'étiquette de son fichier : `// [COMBAT] sert à …`. Une recherche sur `[COMBAT]` dans tout le projet liste donc toutes les fonctions du combat.

| Fichier | Étiquette | Contenu |
|---|---|---|
| `main.go` | `[ACCUEIL]` | `main()` (écran titre), nouvelle partie, « Continuer », prologue, crédits |
| `camp.go` | `[CAMP]` | le menu du camp (boucle principale d'une partie), l'entraînement |
| `character.go` | `[HEROS]` | le héros (`Character`), les classes, les difficultés, la création, la fiche, l'XP, la mort |
| `items.go` | `[INVENTAIRE]` | les objets (`items`), les équipements (`gear`), le sac, utiliser et équiper |
| `merchant.go` | `[MARCHAND]` | acheter et vendre |
| `blacksmith.go` | `[FORGERON]` | les recettes et la fabrication |
| `quests.go` | `[MISSIONS]` | les missions |
| `dungeon.go` | `[DONJON]` | les étages, les salles, l'écran de victoire |
| `sphinx.go` | `[SPHINX]` | les énigmes |
| `monster.go` | `[MONSTRES]` | le bestiaire, `newMonster`, le tour des monstres et des boss |
| `combat.go` | `[COMBAT]` | la boucle de combat, le tour du héros, les dégâts, la victoire et la défaite |
| `spells.go` | `[SORTS]` | les sorts |
| `effects.go` | `[EFFETS]` | saignement, brûlure, étourdissement, affaiblissement, Peau de Pierre |
| `save.go` | `[SAUVEGARDE]` | la sauvegarde et le chargement en JSON |
| `display.go` | `[AFFICHAGE]` | couleurs, dessins, cadres, menus, messages, jauges |
| `input.go` | `[SAISIE]` | lire le clavier |
| `art.go` | `[DESSINS]` | uniquement les dessins, en braille |

---

## 2. Le principe : une pile d'écrans

Chaque écran est une fonction qui boucle : **efface → affiche → lit un numéro → agit**. Elle ne rend la main que lorsque le joueur tape `0` (Retour).

```
main()                        écran titre
 ├─ newGame / continueGame
 │   └─ campMenu()            LE CAMP  ← on y revient toujours
 │       ├─ dungeonMenu()           choix de l'étage
 │       │   └─ exploreFloor()      suite de salles
 │       │       ├─ sphinxRiddle()  le sphinx
 │       │       └─ roomFight()     un monstre → fight()
 │       ├─ merchantMenu()          acheter / vendre
 │       ├─ blacksmithMenu()        les étals de la forge
 │       ├─ questMenu()             les missions
 │       ├─ trainingFight()         combat sans risque
 │       ├─ inventoryMenu()         le sac
 │       ├─ showCharacterSheet()    la fiche du héros
 │       └─ saveGame()
 └─ credits()
```

Les menus sont écrits simplement : une suite d'`option(...)`, puis un `switch` sur le numéro lu par `readChoice`. Les fonctions dont le nom finit par `Menu` sont des écrans qui bouclent.

Le héros est passé partout par pointeur (`character *Character`) : toutes les fonctions modifient le même héros.

---

## 3. Les données dans des tables

Le contenu du jeu est rangé dans des **tables** ; le code ne fait que les parcourir. Ajouter du contenu, c'est ajouter une ligne.

| Table | Fichier | Une ligne = |
|---|---|---|
| `classes` | character.go | une classe : PV, mana, sort, gains par niveau |
| `difficulties` | character.go | une difficulté : `MonsterPercent` et `RewardPercent` |
| `bestiary` | monster.go | un monstre : PV, attaque, butin, effet, dessin, réplique |
| `items` | items.go | un objet : prix, utilisable en combat ou non |
| `gear` | items.go | un équipement : emplacement, bonus, classe de l'arme |
| `spells` | spells.go | un sort : coût, description, dégâts, dessin |
| `shopItems` | merchant.go | un objet vendu, dans l'ordre d'affichage |
| `shelves` | blacksmith.go | un étal de la forge et ses recettes |
| `quests` | quests.go | une mission : cible, nombre, récompense |
| `floors` | dungeon.go | un étage : salles, monstres, boss, couleurs |
| `riddles` | sphinx.go | une énigme et ses 3 réponses |

---

## 4. Le héros et la sauvegarde

`Character` est **la seule structure sauvegardée**. Une sauvegarde est ce `Character` converti en JSON (`json.MarshalIndent`) dans `dungo_save_N.json`. Au chargement, `json.Unmarshal` fait l'inverse.

- Un champ ajouté à `Character` est sauvegardé automatiquement.
- Un champ renommé ou retiré du code est ignoré au chargement (il reste à sa valeur zéro).
- Les champs doivent commencer par une majuscule pour que `encoding/json` les voie.

---

## 5. Le combat

```
fight(héros, monstre, entraînement)
 └─ pour chaque tour :
     ├─ showFightScreen    l'écran de combat
     ├─ characterTurn      Attaquer / Sorts / Objets / Fuir
     ├─ monsterTurn        si le monstre est encore debout      (monster.go)
     ├─ endOfTurnEffects   saignement, brûlure                   (effects.go)
     └─ un PV à 0 ?        loseFight d'abord, sinon winFight
```

- **Dégâts** : `characterHitsMonster` (le héros frappe : affaibli ÷ 2, puis 10 % de critique × 2), `removeMonsterHP` (dégâts bruts : poison, brûlure), `monsterHitsCharacter` (le monstre frappe : Peau de Pierre ÷ 2).
- **Attaques des monstres** : `monsterTurn` choisit selon `monster.Key`. `normalMonsterTurn` pour les monstres ordinaires (× 2 tous les 3 tours), `goblinKingTurn`, `lichTurn` et `dragonTurn` pour les boss.
- **Sorts** (spells.go) : `castSpell` retire le mana, applique les dégâts, puis un `switch` sur le nom du sort ajoute l'effet (brûlure, soin, Peau de Pierre).
- **Effets** (effects.go) : des compteurs dans `Character` (`Bleeding`, `Burning`…), remis à zéro au début de chaque combat.
- **Tour non consommé** : un « Retour » dans un menu, une fuite devant un boss ou un objet refusé (potion à PV max, livre déjà appris) font rechoisir le héros.
- **Entraînement** : aucun gain, pas d'objets, et `trainingFight` rend au héros ses PV et son mana d'avant le combat.

---

## 6. Ajouter du contenu

- **Un monstre** : une ligne dans `bestiary`, puis sa clé dans les `Monsters` d'un étage.
- **Un objet** : une constante et une ligne dans `items`, l'ajouter à `shopItems` pour le vendre, et un `case` dans `useItem` pour son effet.
- **Une arme ou une armure** : une ligne dans `gear`, puis une recette dans un étal de `shelves`.
- **Une mission** : une ligne dans `quests`. `Target` doit être le **nom exact** du monstre.
- **Un sort** : une constante et une ligne dans `spells`, et un `case` dans `castSpell` s'il a un effet en plus des dégâts.

---

## 7. L'affichage (`display.go`) et la saisie (`input.go`)

- Les couleurs sont des **codes ANSI** (`Red = "\033[38;5;196m"`). Toute couleur ouverte est refermée par `Reset`.
- Les dessins sont en **braille Unicode** (`art.go`) : `printArt(dessin, couleurs...)` leur applique un dégradé du haut vers le bas.
- Toute la saisie passe par `readChoice(low, high)`, qui redemande tant que la réponse n'est pas un nombre valide : le jeu ne peut pas planter sur une mauvaise saisie.
