# Comment fonctionne le code de DunGo

Ce document explique l'organisation du code, les structures de données et les
conventions du projet. Il s'adresse à quelqu'un qui doit **lire, corriger ou
enrichir** DunGo — pas au joueur (pour jouer, voir le [README](../README.md)).

---

## 1. Vue d'ensemble

DunGo est un RPG au tour par tour qui tourne entièrement dans le terminal.
Il n'utilise que la bibliothèque standard de Go : aucune dépendance à installer.

```bash
go run ./src        # lancer une partie
go build ./src      # produire un exécutable
go vet ./src        # vérifier le code
```

Tout le programme est dans **un seul paquet** (`package main`, dossier `src/`).
Les fichiers ne sont donc pas des modules étanches : ils servent uniquement à
ranger le code par thème. N'importe quelle fonction peut appeler n'importe
quelle autre.

---

## 2. Le jeu est une pile d'écrans

Il n'y a ni moteur de jeu ni boucle temps réel. Le programme est une suite de
fonctions qui affichent un écran et attendent un chiffre. Chaque écran suit
toujours le même rythme :

```
effacer l'écran  →  afficher  →  lire un choix  →  agir  →  recommencer
```

Le parcours complet :

```
main()
 └─ titleScreen()                    écran titre
     ├─ newGame() / continueGame()
     │   └─ campMenu()               LE CAMP  ← on y revient toujours
     │       ├─ displayInfo()        fiche du héros
     │       ├─ accessInventory()    sac, potions, équipement
     │       ├─ merchant()           acheter / vendre
     │       ├─ blacksmith()         forger
     │       ├─ missionBoard()       contrats de chasse
     │       ├─ trainingFight()      arène (combat sans enjeu)
     │       ├─ exploreDungeon()     choix de l'étage
     │       │   └─ exploreFloor()   suite de salles
     │       │       ├─ roomEvent()  fontaine / colporteur / sphinx
     │       │       └─ fight()      combat tour par tour
     │       └─ saveGame()
     ├─ howToPlay()
     └─ whoAreThey()
```

Une fonction d'écran ne rend la main que lorsque le joueur choisit « Retour ».
C'est ce qui donne la pile : `campMenu` appelle `merchant`, qui appelle
`buyMenu`, qui revient à `merchant`, qui revient à `campMenu`.

---

## 3. Les fichiers

| Fichier | Rôle |
|---|---|
| `main.go` | Écran titre, menu du camp, aide, prologue. Contient les **tables de menus**. |
| `character.go` | Le héros : structure `Character`, lignées, création, niveaux, mort, fiche. |
| `monster.go` | Le bestiaire : structure `Monster` et les statistiques de chaque créature. |
| `combat.go` | La boucle de combat, les sorts, les dégâts, les schémas d'attaque des boss. |
| `status.go` | Les effets temporaires : saignement, brûlure, étourdissement, affaiblissement. |
| `inventory.go` | Le sac, la table `items`, l'usage des objets. |
| `equipment.go` | Les armures et les armes (table `gear`), les emplacements portés. |
| `merchant.go` | Mordecai : achat et revente. |
| `forge.go` | Borin : les recettes et la fabrication. |
| `missions.go` | Les contrats de chasse et leur progression. |
| `dungeon.go` | Les étages, les salles, l'entrée en scène des boss, l'écran de victoire. |
| `events.go` | Fontaine, marchand ambulant et énigmes du sphinx. |
| `save.go` | Sauvegarde et chargement JSON, les 3 emplacements. |
| `utils.go` | **La boîte à outils d'affichage** : couleurs, dessins, cadres, jauges, saisie. |
| `ascii.go` | Uniquement des constantes : tous les dessins, en braille. |

---

## 4. Les structures de données

Le principe du projet : **les données du jeu vivent dans des tables**, le code
ne fait que les parcourir. Pour ajouter du contenu, on ajoute une ligne à une
table ; on ne touche presque jamais à la logique.

### `Character` — le héros (`character.go`)

C'est **la seule structure sauvegardée sur disque**. Tout ce qui doit survivre
à un « Continuer » doit y être ajouté. Elle regroupe :

- les caractéristiques (`HP`, `MaxHP`, `Mana`, `Attack`, `Initiative`, `XP`…) ;
- ce que le héros possède (`Inventory`, `Equipment`, `Gold`, `Skills`) ;
- la progression (`FloorsCleared`, `Victories`, `Deaths`, `QuestIndex`…) ;
- les effets temporaires (`Bleeding`, `Burning`, `Stunned`, `Weakened`,
  `StoneSkin`), remis à zéro au début et à la fin de chaque combat.

Deux petites méthodes évitent des répétitions partout :
`c.knows(sort)` et `c.canCast(sort)`.

### `Class` — les trois lignées (`character.go`)

```go
var classes = []Class{
    {"Humain", 100, 40, 10, SpellSecondWind, "Polyvalent, se soigne", Sky, artHuman, 10, 5},
    ...
}
```

Une ligne = une lignée jouable : PV, mana, initiative, sort de départ, couleur,
portrait et gains par niveau.

### `Difficulty` — Facile / Normal / Difficile (`character.go`)

```go
var difficulties = []Difficulty{
    {"Facile", "Monstres affaiblis (-25 % de PV et d'attaque)", Green, 75, 100},
    {"Normal", "L'aventure telle qu'elle a été pensée", Yellow, 100, 100},
    {"Difficile", "Monstres renforcés (+30 %), mais +50 % d'XP et de Y-Coins", Red, 130, 150},
}
```

`Power` et `Reward` sont des **pourcentages**, utilisés à la fois pour écrire le
menu et pour calculer les statistiques dans `newMonster`. Les deux ne peuvent
donc pas se contredire.

### `Monster` — le bestiaire (`monster.go`)

```go
"troll": {Name: "Troll des cavernes", MaxHP: 85, Attack: 12, Initiative: 4,
          XP: 55, Gold: 14, Drop: ItemTrollSkin, DropChance: 75,
          Effect: EffectStun, Art: artTroll, Color: Green, Cry: "..."},
```

`newMonster(kind, difficulty)` renvoie une **copie** de l'entrée du bestiaire,
mise à l'échelle de la difficulté. Le bestiaire lui-même n'est jamais modifié :
c'est pour cela qu'on peut affronter dix trolls d'affilée.

Le champ `Pattern` est une **fonction** : c'est la façon d'attaquer du monstre.
S'il vaut `nil`, le monstre utilise `basicPattern` (attaque doublée tous les
3 tours). Les boss ont chacun le leur : `goblinKingPattern`, `lichPattern`,
`dragonPattern`.

### `Item` et `Gear` — les objets (`inventory.go`, `equipment.go`)

Deux tables, deux natures d'objets :

- `items` décrit ce qui se consomme ou se revend : catégorie, couleur, prix
  chez Mordecai, et `InFight` (utilisable en plein combat) ;
- `gear` décrit ce qui se porte : emplacement, PV, attaque, lignée de l'arme.

Le sac lui-même est une simple `map[string]int` : nom de l'objet → quantité.

### `Spell` — les sorts (`combat.go`)

```go
SpellFireball: {Cost: 15, Description: "18 dégâts + brûlure", Damage: 18,
                Art: artFireball, Colors: fireColors, Extra: castBurn},
```

`Damage` suffit pour un sort offensif simple. `Extra` est une fonction
optionnelle pour tout le reste (soigner, enflammer, protéger).

### `Floor`, `Quest`, `Recipe`, `Riddle`

Mêmes principes : une table, une ligne par contenu.

---

## 5. La boîte à outils d'affichage (`utils.go`)

C'est le fichier le plus utilisé du projet. Tout l'affichage passe par lui,
ce qui garantit que le jeu a partout le même style.

### Couleurs

Les couleurs sont des **codes ANSI** : des chaînes que le terminal interprète
au lieu de les afficher.

```go
const Red = "\033[38;5;196m"   // à partir d'ici, écris en rouge
const Reset = "\033[0m"        // reviens à la couleur normale
```

Règle d'or : **toute couleur ouverte doit être refermée par `Reset`** sur la
même ligne. Sinon la couleur « bave » sur le texte suivant.

Les **dégradés** sont de simples listes de couleurs (`fireColors`,
`iceColors`, `goldColors`, `steelColors`…) appliquées du haut vers le bas d'un
dessin.

### Dessins

```go
printArt(artMerchant, goldColors...)   // dégradé doré
printArt(artBag, Brown)                // une seule couleur = uni
printArtSlow(artLogo, fireColors...)   // apparition ligne par ligne
```

Les deux passent par `drawArt`, qui découpe le dessin en lignes et attribue à
chacune une couleur du dégradé :

```go
colors[i*len(colors)/len(lines)]
```

Avec 12 lignes et 6 couleurs, chaque couleur couvre 2 lignes. Avec une seule
couleur, l'indice vaut toujours 0 : le dessin est uni. **Il n'existe qu'un seul
chemin d'affichage** pour les dessins, donc ils ont tous exactement le même
rendu (même graisse, même façon de refermer les couleurs).

### Textes

Le terminal fait environ 80 colonnes. Les textes longs sont repliés à
`textWidth` (76) par `wrap`, utilisé par :

- `typewrite(couleur, texte)` — apparition lettre par lettre ;
- `paragraph(couleur, texte)` — affichage immédiat ;
- `say(nom, couleur, réplique)` — dialogue, aligné sous le nom du personnage.

### Cadres, listes et jauges

| Fonction | Ce qu'elle dessine |
|---|---|
| `banner(texte, couleur)` | Le cadre `╔═══╗` des titres d'écran |
| `section(texte, couleur)` | Le séparateur `── Titre ────────` |
| `option(n, libellé, couleur)` | Une entrée de menu `[1] Attaquer` |
| `itemLine(n, nom, couleur, détail)` | Une ligne de liste d'objets |
| `bar(actuel, max, largeur, couleur)` | Une jauge `████░░░░` |
| `hpBar(...)` | Une jauge de PV qui passe du vert au rouge |
| `showHP` / `showMana` | Une jauge nommée, prête à afficher |
| `success` / `fail` / `info` / `warn` | Les messages `√ × • !` |

### Saisie

`readChoice(low, high)` est le seul point d'entrée du clavier pour les menus :
il redemande tant que la réponse n'est pas un nombre dans l'intervalle. On ne
peut donc pas faire planter le jeu en tapant n'importe quoi.
`ask("question")` est le raccourci oui/non.

---

## 6. Les dessins (`ascii.go`)

Tous les dessins sont en **braille Unicode** (bloc `U+2800`). Un caractère
braille contient une grille de 2 × 4 points, ce qui donne quatre fois plus de
précision qu'un caractère ordinaire. Les « vides » ne sont pas des espaces mais
`⠀` (U+2800, le braille vide) : l'alignement est ainsi garanti.

```go
const artPotion = `
⠀⠀⠀⠀⢸⣿⣿⡇
⠀⠀⠀⠀⠀⣿⣿⠁
...
`
```

Deux contraintes à respecter en ajoutant un dessin :

1. **la largeur** : 33 caractères au maximum (≈ 66 points), pour tenir dans un
   terminal de 80 colonnes ;
2. **la hauteur** : environ 12 lignes pour un monstre, sinon l'écran de combat
   (dessin + jauges + menu ≈ 15 lignes d'interface) ne tient plus dans une
   fenêtre courte.

Les portraits des lignées sont un cas à part : `sideBySide` les affiche en
trois colonnes de 24 caractères, ils doivent donc rester sous cette largeur.

---

## 7. Le déroulement d'un combat (`combat.go`)

```
fight(héros, monstre, entraînement)
 │
 ├─ clearEffects        on nettoie les effets du combat précédent
 ├─ initiative          le plus rapide commencera
 │
 └─ pour chaque tour :
     ├─ showFight       dessin du monstre, jauges, effets
     ├─ tour du plus rapide, puis tour de l'autre
     │   ├─ characterTurn : Attaquer / Sorts / Inventaire / Fuir
     │   └─ monsterTurn   : m.Pattern, ou basicPattern par défaut
     ├─ après chaque demi-tour : fightOver ? (victoire, défaite)
     ├─ updateEffects   saignement, brûlure : dégâts de fin de tour
     └─ fightOver ?
```

`fight` renvoie `Victory`, `Defeat` ou `Fled`. L'appelant (`exploreFloor`) s'en
sert pour décider s'il continue l'étage.

Quelques règles utiles à connaître :

- **coup critique** : 10 % de chances, dégâts doublés (`damageMonster`) ;
- **attaque puissante** : tous les 3 tours, le monstre frappe deux fois plus
  fort et inflige son `Effect` ;
- **l'arène ne rapporte rien** : ni Y-Coins, ni XP, ni victoire comptabilisée —
  sinon on y gagnerait des niveaux sans aucun risque. On n'y meurt pas non plus ;
- **les boss ne peuvent pas être fuis** (`flee` refuse si `m.IsBoss`).

---

## 8. La sauvegarde (`save.go`)

Une partie est un fichier JSON à côté de l'exécutable :
`dungo_save_1.json`, `dungo_save_2.json`, `dungo_save_3.json`.

Le contenu est **exactement** la structure `Character`, écrite par
`encoding/json`. Conséquences pratiques :

- ajouter un champ à `Character` suffit à le sauvegarder ;
- un champ supprimé du code est simplement ignoré au chargement ;
- une ancienne sauvegarde reste lisible, les nouveaux champs valant zéro.

---

## 9. Ajouter du contenu

### Un monstre

1. Dessiner l'art dans `ascii.go` (≤ 33 colonnes, ≈ 12 lignes).
2. Ajouter une ligne à `bestiary` dans `monster.go`.
3. Citer sa clé dans les `Monsters` d'un étage, dans `floors` (`dungeon.go`).

Pour lui donner un comportement particulier, écrire une fonction
`func monNomPattern(m *Monster, c *Character, turn int)` dans `combat.go` et la
renseigner dans le champ `Pattern`.

### Un objet

1. Une constante `ItemQuelqueChose = "Nom affiché"` dans `inventory.go`.
2. Une ligne dans la table `items` (catégorie, couleur, prix, `InFight`).
3. Pour qu'il soit vendu au camp : l'ajouter à `shopItems` (`merchant.go`).
4. Pour qu'il ait un effet : un `case` dans `useItem`.

### Une arme ou une armure

1. Une constante et une ligne dans `gear` (`equipment.go`).
2. Une ligne dans `recipes` (`forge.go`) pour pouvoir la fabriquer.

### Une mission

Une ligne dans `quests` (`missions.go`). Le champ `Target` doit contenir le
**nom exact** du monstre, celui du bestiaire : c'est sur cette chaîne que
`updateQuest` compte les victoires.

### Un étage

Une ligne dans `floors` (`dungeon.go`) : nom, texte d'ambiance, dégradé, nombre
de salles, liste des monstres et clé du boss.

### Un sort

1. Une constante `SpellX` dans `character.go`.
2. Une ligne dans `spells` (`combat.go`). `Damage` suffit pour un sort
   offensif ; sinon écrire une fonction et la mettre dans `Extra`.
3. Le donner au héros : via une lignée (`classes`) ou un objet.

---

## 10. Conventions du code

- **Français partout** : noms d'affichage, commentaires, messages.
  Les identifiants Go restent en anglais (`Character`, `fight`, `damage`).
- **Les commentaires expliquent le pourquoi**, pas le quoi. Un commentaire qui
  paraphrase la ligne suivante est inutile.
- **Pas d'affichage brut** : on passe par `utils.go`. Un `fmt.Println` avec des
  couleurs écrites à la main finit toujours par oublier un `Reset`.
- **Pas de nombre magique dupliqué** : si une valeur apparaît dans un texte
  *et* dans un calcul, elle doit venir d'une table (voir `Difficulty`).
- **Une fonction = un écran, ou un calcul.** Si une fonction fait les deux,
  elle est trop grosse.
- `go vet ./src` doit rester silencieux.
