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

**Un jeu de rôle dans le terminal, écrit en Go.**

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Projet RED](https://img.shields.io/badge/Projet-RED-C0392B?style=for-the-badge)

*Projet RED · Ymmersion — Paco, Sofiane, Valentin et Ayman*

</div>

---

## 🐉 L'histoire

Un dragon, **Ignarok**, s'est installé au fond du donjon. Il brûle les champs, mange les moutons et ronfle si fort que plus personne ne dort. Le village cherche un héros… et vous êtes le seul volontaire.

```
Camp  ──►  Étage 1  ──►  Étage 2  ──►  Étage 3
           Grukk        Mor'Vath      Ignarok
```

Préparez-vous au camp, puis descendez battre le boss de chaque étage. Le dernier, c'est le dragon.

---

## 🚀 Lancer le jeu

1. Installez **Go** : https://go.dev/dl/
2. Récupérez le projet et lancez-le :

```bash
git clone https://github.com/Pacoclrt/DunGo.git
cd DunGo
go run ./src
```

> 💡 Agrandissez la fenêtre du terminal pour voir les dessins en entier.

---

## 🎮 Jouer

- Tapez le **numéro** d'une option, puis **Entrée**.
- **0** permet toujours de revenir en arrière.
- Quand le jeu affiche `[ Entrée pour continuer ]`, appuyez sur **Entrée**.

### Votre héros

| Classe | PV | Mana | Sort spécial |
|---|:---:|:---:|---|
| **Humain** | 100 | 40 | *Second Souffle* : se soigne |
| **Elfe** | 80 | 60 | *Flèche d'Argent* : gros dégâts |
| **Nain** | 120 | 30 | *Peau de Pierre* : encaisse moitié moins |

Trois difficultés : **Facile**, **Normal** ou **Difficile** (monstres plus forts, mais plus d'XP et de Y-Coins).

### Le camp

| Menu | À quoi ça sert |
|---|---|
| 🏰 **Explorer le donjon** | Partir à l'aventure |
| 🪙 **Marchand** | Acheter des potions et des ressources, revendre des objets. La première Potion de vie est offerte ! |
| 🔨 **Forgeron** | Fabriquer des armures (+PV) et des armes (+attaque) avec les ressources laissées par les monstres |
| 📜 **Missions** | Chasser des monstres contre des Y-Coins et de l'XP |
| ⚔️ **Entraînement** | Un combat pour de faux : aucun risque, aucun butin |
| 🎒 **Inventaire** | Boire une potion, équiper une arme ou une armure |
| 📋 **Fiche du héros** | Statistiques, équipement et sorts |
| 💾 **Sauvegarder** | 3 emplacements de sauvegarde |

### Le donjon

Chaque étage est une suite de salles. La dernière salle cache un **boss**, et le battre débloque l'étage suivant.

| Étage | Boss | Son pouvoir |
|---|---|---|
| 1 · Les Galeries Gobelines | **Grukk**, le Roi Gobelin | Appelle un garde qui le soigne |
| 2 · Les Cryptes Englouties | **Mor'Vath**, la Liche | Vole votre mana |
| 3 · L'Antre d'Ignarok | **Ignarok**, le dragon | Crache du feu |

Certaines salles cachent un **sphinx** : une énigme facile, bonne réponse = récompense, mauvaise = un monstre. Entre deux salles, vous pouvez remonter au camp, mais l'étage recommencera du début.

### Les combats

Vous jouez toujours en premier : **Attaquer**, lancer un **Sort** (coûte du mana), utiliser un **Objet** ou **Fuir** (1 chance sur 2, impossible contre un boss).

- Tous les **3 tours**, les monstres frappent deux fois plus fort.
- **10 %** de chance de faire un coup critique (dégâts doublés).
- Certains monstres font **saigner**, **brûler**, **étourdir** ou **affaiblir**.
- Si vos PV tombent à 0, vous revenez au camp… en perdant **20 %** de vos Y-Coins.

### Conseils pour débuter

1. Récupérez votre **potion gratuite** chez le marchand.
2. Acceptez une **mission** avant de partir.
3. Remontez au camp si vos PV sont bas.
4. Gardez des **potions** pour les boss !

---

## 🗂️ Pour les développeurs

```
src/        le code du jeu (un seul paquet Go, sans dépendance)
docs/       comment-ca-marche.md : l'organisation du code
```

```bash
go run ./src     # lancer le jeu
go vet ./src     # vérifier le code
```

Les données du jeu (classes, monstres, objets, sorts, recettes, missions, étages) sont rangées dans des tableaux faciles à modifier. Tout est expliqué dans [docs/comment-ca-marche.md](docs/comment-ca-marche.md).

---

<div align="center">

**Paco · Sofiane · Valentin · Ayman**

*Projet RED — Ymmersion*

</div>
