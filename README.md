<div align="center">

```
██████╗ ██╗   ██╗███╗   ██╗ ██████╗  ██████╗
██╔══██╗██║   ██║████╗  ██║██╔════╝ ██╔═══██╗
██║  ██║██║   ██║██╔██╗ ██║██║  ███╗██║   ██║
██║  ██║██║   ██║██║╚██╗██║██║   ██║██║   ██║
██████╔╝╚██████╔╝██║ ╚████║╚██████╔╝╚██████╔╝
╚═════╝  ╚═════╝ ╚═╝  ╚═══╝ ╚═════╝  ╚═════╝
         ~ L'ANTRE D'IGNAROK ~
```

**Un jeu de rôle au tour par tour, dans le terminal, écrit en Go.**

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Projet RED](https://img.shields.io/badge/Projet-RED-C0392B?style=for-the-badge)

*Projet RED · Ymmersion — Paco, Sofiane, Valentin et Ayman*

</div>

---

## 🐉 Le but du jeu

Un dragon, **Ignarok**, s'est installé au fond du donjon. Il brûle les champs, mange les moutons et ronfle si fort que plus personne ne dort. Le village cherche un héros… et vous êtes le seul volontaire.

Votre mission : traverser les **3 étages du donjon** et battre leurs boss. Le dernier, c'est le dragon.

```
Le camp  ──►  Étage 1  ──►  Étage 2  ──►  Étage 3  ──►  Victoire !
              Grukk         Mor'Vath      Ignarok
```

---

## 🚀 Lancer le jeu

1. Installez **Go** (version 1.22 ou plus récente) : https://go.dev/dl/
2. Récupérez le projet, puis lancez le jeu **depuis le dossier du projet** :

```bash
git clone https://github.com/Pacoclrt/DunGo.git
cd DunGo
go run ./src
```

> 💡 Agrandissez la fenêtre du terminal : les dessins sont larges. Le jeu a besoin d'un terminal récent, qui affiche les couleurs et les caractères braille (Terminal de macOS, Windows Terminal, terminal de VS Code…).

Pour fabriquer un exécutable au lieu de lancer le jeu directement :

```bash
go build -o dungo ./src
./dungo
```

---

## 🎮 Comment jouer

Tout se joue avec des **numéros** :

- tapez le **numéro** d'une option, puis **Entrée** ;
- **0** permet toujours de revenir en arrière ;
- quand le jeu affiche `[ Entrée pour continuer ]`, appuyez sur **Entrée**.

---

## ⚙️ Comment fonctionne le jeu

### 1. Créer son héros

Choisissez un nom, une classe et une difficulté.

| Classe | PV | Mana | Sort de classe |
|---|:---:|:---:|---|
| **Humain** | 100 | 40 | *Second Souffle* : se soigne |
| **Elfe** | 80 | 60 | *Flèche d'Argent* : gros dégâts |
| **Nain** | 120 | 30 | *Peau de Pierre* : encaisse moitié moins |

Trois difficultés : **Facile**, **Normal** ou **Difficile** (monstres plus forts, mais plus d'XP et de Y-Coins).

### 2. Se préparer au camp

Le camp est le point de départ de chaque expédition : on y revient toujours.

| Menu | À quoi ça sert |
|---|---|
| 🏰 **Explorer le donjon** | partir à l'aventure |
| 🪙 **Marchand** | acheter des potions et des ressources, revendre ses objets (la 1re Potion de vie est offerte) |
| 🔨 **Forgeron** | fabriquer des armures (+PV) et des armes (+attaque) avec les ressources laissées par les monstres |
| 📜 **Missions** | chasser des monstres contre des Y-Coins et de l'XP |
| ⚔️ **Entraînement** | un combat pour de faux : aucun risque, aucun butin |
| 🎒 **Inventaire** | boire une potion, équiper une arme ou une armure |
| 📋 **Fiche du héros** | statistiques, équipement et sorts |
| 💾 **Sauvegarder** | 3 emplacements de sauvegarde |

### 3. Explorer le donjon

Chaque étage est une suite de salles, avec un monstre dans chacune. La **dernière salle cache un boss** : le battre débloque l'étage suivant.

| Étage | Boss | Son pouvoir |
|---|---|---|
| 1 · Les Galeries Gobelines | **Grukk**, le Roi Gobelin | appelle un garde qui le soigne |
| 2 · Les Cryptes Englouties | **Mor'Vath**, la Liche | vole votre mana |
| 3 · L'Antre d'Ignarok | **Ignarok**, le dragon | crache du feu, et enrage quand il est blessé |

- Certaines salles cachent un **sphinx** : bonne réponse à son énigme = une récompense, mauvaise = un monstre.
- Entre deux salles, vous pouvez **remonter au camp**… mais l'étage recommencera du début.

### 4. Combattre

C'est du tour par tour, et vous jouez toujours en premier : **Attaquer**, lancer un **Sort** (coûte du mana), utiliser un **Objet** ou **Fuir** (1 chance sur 2, impossible contre un boss).

- Tous les **3 tours**, les monstres lancent une attaque puissante.
- **10 %** de chance de faire un coup critique (dégâts doublés).
- Certains monstres font **saigner**, **brûler**, **étourdir** ou **affaiblir**.
- Chaque victoire rapporte de l'**XP** (pour monter de niveau), des **Y-Coins** et parfois un objet.
- Si vos PV tombent à 0, vous revenez au camp… en perdant **20 %** de vos Y-Coins.

---

## 🧪 Tester le jeu sans farmer : la sauvegarde « Test Dev »

Le projet contient une sauvegarde de test, `dungo_save_1.json`, avec un héros surpuissant. Elle permet de visiter tout le jeu sans passer des heures à gagner de l'XP et des Y-Coins.

**Pour la lancer :** `go run ./src`, puis **Continuer**, puis l'**emplacement 1**.

| « Test Dev » (Elfe, Facile) | |
|---|---|
| PV et mana | 999 |
| Attaque | 999 (tous les monstres tombent en un coup) |
| Y-Coins | 9 999 (tout acheter chez le marchand) |
| Sac | 100 places, avec 10 Potions de vie |

> ⚠️ Lancez le jeu depuis le dossier du projet : c'est là que se trouve le fichier de sauvegarde. Et évitez de choisir **Nouvelle partie** sur l'emplacement 1, sinon la sauvegarde de test sera écrasée à la première sauvegarde.

Le fichier est du texte (JSON) : on peut l'ouvrir et changer les valeurs à la main, par exemple `"FloorsCleared": 2` pour aller directement à l'étage du dragon.

---

## 💡 Conseils pour débuter

1. Récupérez votre **potion gratuite** chez le marchand : vous commencez avec seulement la moitié de vos PV.
2. Acceptez une **mission** avant de partir.
3. Remontez au camp si vos PV sont bas.
4. Gardez des **potions** pour les boss !

---

## 🗂️ Pour les développeurs

```
src/     le code du jeu : 17 fichiers, un par thème (Go, sans aucune dépendance)
docs/    DunGo-guide-complet.md (+ PDF) : comment marche le code, pour débutants en Go
         comment-ca-marche.md : l'organisation du code, en bref
         DunGo-antiseche-oral.md : l'antisèche pour l'oral
```

Chaque fonction porte un commentaire sur sa ligne, qui commence par l'étiquette de son thème : `// [COMBAT] sert à …`. Une recherche sur `[COMBAT]`, `[MARCHAND]` ou `[SAUVEGARDE]` dans tout le projet liste toutes les fonctions de ce thème. Tout est expliqué dans le [guide du code](docs/DunGo-guide-complet.md).

---

<div align="center">

**Paco · Sofiane · Valentin · Ayman**

*Projet RED — Ymmersion*

</div>
