# Gestion de projet — DunGo

> **Projet RED · Ymmersion**
> Équipe : **Paco**, **Sofiane**, **Valentin**, **Ayman**
> Dépôt : `projet-red_DunGo`
>
> *Modèle pré-rempli : complétez les zones marquées « À compléter ».*

---

## 1. Présentation

| | |
|---|---|
| **Nom** | DunGo — L'Antre d'Ignarok |
| **Type** | RPG en ligne de commande, joué dans le terminal de VS Code |
| **Langage** | Go (bibliothèque standard uniquement) |
| **Univers** | Dark fantasy avec touches d'humour : un camp, 3 étages de donjon et le dragon Ignarok |

**Objectif** : utiliser les notions de Go vues pendant l'Ymmersion (structures, pointeurs, maps, slices, boucles, `switch`, fonctions, gestion des erreurs, JSON) pour créer un jeu complet : personnage, inventaire, marchand, forge et combat au tour par tour.

---

## 2. Organisation de l'équipe

| Membre | Rôle principal | Fichiers |
|---|---|---|
| Paco | *À compléter* | *ex. `character.go`, `inventory.go`* |
| Sofiane | *À compléter* | *ex. `combat.go`, `monster.go`, `dungeon.go`* |
| Valentin | *À compléter* | *ex. `merchant.go`, `forge.go`, `equipment.go`* |
| Ayman | *À compléter* | *ex. `ascii.go`, `utils.go`, `main.go`* |

| Besoin | Outil |
|---|---|
| Code et versionnement | Git + GitHub |
| Éditeur | VS Code + extension Go |
| Suivi des tâches | *À compléter (GitHub Projects, Trello, Notion…)* |
| Communication | *À compléter (Discord, Teams…)* |

**Conventions**
- Tout le code est dans `src/` (`package main`), un fichier par thème.
- Noms Go en anglais, commentaires en français.
- Les noms de fonctions imposés par le sujet sont conservés (`initCharacter`, `displayInfo`, `accessInventory`, `takePot`, `isDead`, `poisonPot`, `spellBook`, `characterCreation`, `upgradeInventorySlot`, `initGoblin`, `goblinPattern`, `characterTurn`, `trainingFight`).
- Avant chaque commit : `gofmt -w src` et `go vet ./src`.

---

## 3. Planning

| Phase | Contenu | Période | Statut |
|---|---|---|---|
| 1. Cadrage | Lecture du sujet, univers, répartition | *À compléter* | ☐ |
| 2. Partie 1 | Personnage, inventaire, marchand, sorts (tâches 1 à 12) | *À compléter* | ☐ |
| 3. Partie 2 | Économie, forge, équipement (tâches 13 à 18) | *À compléter* | ☐ |
| 4. Partie 3 | Combat au tour par tour (tâches 19 à 22) | *À compléter* | ☐ |
| 5. Missions | Initiative, XP, sorts, mana, donjon, crédits | *À compléter* | ☐ |
| 6. Finitions | ASCII art, équilibrage, README | *À compléter* | ☐ |
| 7. Oral | Répétition de la démonstration | *À compléter* | ☐ |

---

## 4. Suivi des tâches

| # | Tâche | Fonction(s) | Responsable | Statut |
|---|---|---|---|---|
| 1 | Personnage | `Character` | *À compléter* | ☑ |
| 2 | Initialisation | `initCharacter` | *À compléter* | ☑ |
| 3 | Infos du personnage | `displayInfo` | *À compléter* | ☑ |
| 4 | Inventaire | `accessInventory` | *À compléter* | ☑ |
| 5 | Potion de vie | `takePot` | *À compléter* | ☑ |
| 6 | Menu | `campMenu` | *À compléter* | ☑ |
| 7 | Marchand | `merchant`, `addInventory`, `removeInventory` | *À compléter* | ☑ |
| 8 | Wasted | `isDead` | *À compléter* | ☑ |
| 9 | Potion de poison | `poisonPot` | *À compléter* | ☑ |
| 10 | Sorts | `Skills`, `spellBook` | *À compléter* | ☑ |
| 11 | Création du personnage | `characterCreation`, `formatName` | *À compléter* | ☑ |
| 12 | Limite d'inventaire | `canAddItem` | *À compléter* | ☑ |
| 13 | Argent | `Gold` | *À compléter* | ☑ |
| 14 | Prix | `prices`, `buyItem` | *À compléter* | ☑ |
| 15 | Forgeron | `blacksmith`, `craft` | *À compléter* | ☑ |
| 16 | Equipment | `Equipment` | *À compléter* | ☑ |
| 17 | Équiper | `equipItem` | *À compléter* | ☑ |
| 18 | Agrandir le sac | `upgradeInventorySlot` | *À compléter* | ☑ |
| 19 | Monster | `Monster`, `initGoblin` | *À compléter* | ☑ |
| 20 | IA du gobelin | `goblinPattern` | *À compléter* | ☑ |
| 21 | Tour du joueur | `characterTurn` | *À compléter* | ☑ |
| 22 | Entraînement | `trainingFight`, `fight` | *À compléter* | ☑ |
| M1 | Initiative | `fight` | *À compléter* | ☑ |
| M2 | Expérience | `gainXP` | *À compléter* | ☑ |
| M3 | Combat magique | `spellMenu` | *À compléter* | ☑ |
| M4 | Mana | `spells`, `takeManaPot` | *À compléter* | ☑ |
| M5 | Améliorations | Boss, effets de statut, événements, armes, missions, difficultés, 3 sauvegardes, critiques, ASCII art braille animé | *À compléter* | ☑ |
| M6 | Qui sont-ils ? | `whoAreThey` | *À compléter* | ☑ |

---

## 5. Choix techniques

| Choix | Pourquoi |
|---|---|
| **Terminal intégré de VS Code** | Tout le monde lance le jeu de la même façon (`go run ./src` ou F5), sans configuration de console. |
| **Bibliothèque standard uniquement** | Rien à installer à part Go, et tout le code est écrit par l'équipe. |
| **`bufio.Scanner` + `readChoice`** | Une seule fonction lit les choix et redemande tant que la saisie est invalide. |
| **Menus en boucle `for` + `switch`** | Même structure pour tous les écrans : facile à lire et à étendre. |
| **Inventaire en `map[string]int`** | Nom de l'objet → quantité : ajout, retrait et comptage simples. |
| **ASCII art en constantes** | Les dessins braille sont regroupés dans `ascii.go` et affichés avec `printArt` (uni ou dégradé) ou `printArtSlow` (apparition ligne par ligne). |
| **Couleurs 256 en constantes** | `"\033[38;5;208m"` = orange : une constante par couleur et quelques dégradés (`fireColors`, `goldColors`…). |
| **Fonctions d'affichage réutilisables** | `banner`, `section`, `say`, `hpBar`, `showStatus` : chaque écran reste court et homogène. |
| **Sauvegarde JSON** | `encoding/json` enregistre directement la structure `Character`. |
| **Donjon en menus** | Une salle = un combat ou un événement (fontaine, marchand, sphinx), et un boss à la fin de chaque étage. |
| **Un fichier par fonctionnalité** | `missions.go`, `events.go`, `status.go` : chaque nouveauté reste dans son fichier, facile à présenter. |
| **Effets de statut en compteurs** | Un simple entier par effet (`Bleeding`, `Burning`…) qui diminue à chaque tour. |
| **3 emplacements de sauvegarde** | Un fichier JSON par emplacement (`dungo_save_1.json`…). |

**Équilibrage** : la monnaie du jeu s'appelle les **Y-Coins**. Départ à 50 Y-Coins au lieu de 100 et prix plus élevés. Les monstres laissent des Y-Coins et des ressources, et le camp soigne gratuitement.

**Mode test** : un héros nommé « test » obtient PV et mana infinis, 999 999 Y-Coins, une attaque à 999, tous les sorts et les 3 étages débloqués (`activateTestMode` dans `character.go`). Pratique pour tester et pour la démonstration.

---

## 6. Difficultés rencontrées

| Difficulté | Solution |
|---|---|
| *ex. Saisies invalides (lettres au lieu d'un nombre)* | *`readChoice` redemande tant que le nombre n'est pas valide* |
| *ex. Remplacer un équipement déjà porté* | *Retirer le bonus de l'ancien objet et le remettre dans le sac* |
| *ex. Équilibrer le combat contre le dragon* | *Attaque annoncée un tour à l'avance, parties de test* |
| *À compléter* | *À compléter* |

---

## 7. Préparation de l'oral

**Déroulé de la démonstration (≈ 10 min)**
1. Présentation de l'univers et de l'équipe. — *À compléter : qui parle*
2. `go run ./src` dans VS Code : écran titre, création (nom mal saisi pour montrer la mise en forme), choix de classe. — *À compléter*
3. Camp : fiche du héros, marchand (potion offerte, achat refusé faute d'or), inventaire. — *À compléter*
4. Entraînement : initiative, attaque ×2 au 3e tour, sorts, potions. — *À compléter*
5. Forge et équipement, puis donjon : combat, butin, montée de niveau. — *À compléter*
6. Code : organisation des fichiers, correspondance tâches → fonctions. — *À compléter*
7. « Qui sont-ils ? » et conclusion. — *À compléter*

**Questions probables du jury**
- Pourquoi utiliser des pointeurs `*Character` ?
- Comment vérifiez-vous la limite d'inventaire à chaque ajout ?
- Comment fonctionne l'attaque ×2 tous les 3 tours (`turn % 3`) ?
- Comment la sauvegarde fonctionne-t-elle ?
- Qu'amélioreriez-vous avec plus de temps ?

---

## 8. Bilan

- **Ce qui a bien fonctionné** : *À compléter*
- **Ce que nous améliorerions** : *À compléter*
- **Ce que chacun a appris** : *À compléter*
