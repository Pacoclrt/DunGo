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

**Un jeu de rôle en ligne de commande, écrit en Go.**

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![VS Code](https://img.shields.io/badge/VS%20Code-terminal-007ACC?style=for-the-badge&logo=visualstudiocode&logoColor=white)
![Projet RED](https://img.shields.io/badge/Projet-RED-C0392B?style=for-the-badge)

*Projet RED · Ymmersion — Paco, Sofiane, Valentin et Ayman*

</div>

---

## 🐉 C'est quoi DunGo ?

Un dragon, **Ignarok**, s'est réveillé sous la montagne. Douze héros sont partis le combattre… aucun n'est revenu.

Vous êtes le treizième. Votre mission :

1. **Créer** votre héros et choisir la difficulté.
2. **Vous préparer** au camp : marchand, forgeron, missions, entraînement.
3. **Traverser** les 3 étages du donjon et vaincre leurs **boss**.
4. **Abattre** le dragon Ignarok.

Tout se passe dans le terminal, avec des dessins en **ASCII art** et des **couleurs**.

---

## 🚀 Lancer le jeu

### Étape 1 : installer les outils (une seule fois)

| Outil | Où le trouver |
|---|---|
| **Go** (le langage du jeu) | https://go.dev/dl/ |
| **VS Code** (l'éditeur) | https://code.visualstudio.com/ |
| **Extension Go** pour VS Code | Dans VS Code : onglet Extensions → chercher « Go » → Installer |

Pour vérifier que Go est bien installé, ouvrez un terminal et tapez :

```bash
go version
```

### Étape 2 : récupérer le projet

```bash
git clone https://github.com/Pacoclrt/DunGo.git
```

Puis dans VS Code : **Fichier → Ouvrir le dossier…** → choisir le dossier `DunGo`.

### Étape 3 : jouer

1. Ouvrez le terminal de VS Code : menu **Terminal → Nouveau terminal**.
2. Tapez :

```bash
go run ./src
```

3. C'est parti ! 🎮

> 💡 **Plus rapide** : appuyez sur **F5** dans VS Code, le jeu se lance tout seul.
>
> 💡 **Plus beau** : agrandissez le panneau du terminal pour voir les dessins en entier.

---

## 🎮 Comment jouer

### Les commandes

Tout se joue avec des **numéros** :

```
   [1] Marchand
   [2] Forgeron
   [0] Retour

  ► 1        ← tapez le numéro, puis appuyez sur Entrée
```

- **0** permet toujours de **revenir en arrière**.
- Quand le jeu affiche `[ Appuyez sur Entrée pour continuer ]`, appuyez simplement sur **Entrée**.

### Nouvelle partie

1. Choisissez un **emplacement de sauvegarde** (3 disponibles).
2. Donnez un **nom** à votre héros et choisissez sa **classe**.
3. Choisissez la **difficulté** :

| Difficulté | Effet |
|---|---|
| **Facile** | Monstres plus faibles |
| **Normal** | L'aventure classique |
| **Difficile** | Monstres plus forts, mais plus d'XP et de Y-Coins |

### Les classes

| Classe | Points forts | Sort spécial | Arme de prédilection |
|---|---|---|---|
| **Humain** | Équilibré (100 PV) | *Second Souffle* : se soigne | Épée |
| **Elfe** | Rapide et magicien (80 PV) | *Flèche d'Argent* : gros dégâts | Arc |
| **Nain** | Très résistant (120 PV) | *Peau de Pierre* : encaisse moins de dégâts | Marteau |

Tout le monde commence avec **50 Y-Coins** (la monnaie du jeu), une **Potion de vie** et le sort **Coup de poing**.

### Le camp

C'est votre base. Vous y revenez entre deux expéditions.

| Lieu | À quoi ça sert |
|---|---|
| 🪙 **Marchand** | Acheter des potions et des objets, revendre ce que vous ne voulez plus |
| 🔨 **Forgeron** | Fabriquer des **armures** (+PV) et des **armes** (+attaque) |
| 📜 **Missions** | Accepter des **missions** de chasse récompensées en Y-Coins et en XP |
| ⚔️ **Entraînement** | Combattre un gobelin pour s'exercer : aucun risque, mais aucun butin |
| 🏰 **Donjon** | Partir à l'aventure : 3 étages, chacun gardé par un boss |
| 💾 **Sauvegarder** | Enregistrer la partie dans son emplacement |

### Le donjon

Chaque étage est une suite de **salles**. Dans une salle, vous pouvez tomber sur :

- 👹 un **monstre** à combattre ;
- 💧 une **fontaine** qui soigne ;
- 🧳 un **marchand ambulant** qui vend des potions (un peu plus cher) ;
- 🗿 un **sphinx** qui pose une énigme toute simple : bonne réponse = récompense, mauvaise = un monstre surgit !

La **dernière salle** de chaque étage contient un **boss** (impossible de le fuir) :

| Étage | Boss | Son pouvoir |
|---|---|---|
| 1 · Les Galeries Gobelines | **Grukk, le Roi Gobelin** | Appelle un garde qui le soigne |
| 2 · Les Cryptes Englouties | **Mor'Vath, la Liche** | Aspire votre mana et vous affaiblit |
| 3 · L'Antre d'Ignarok | **Ignarok, le dragon** | Souffle de feu qui brûle |

Entre deux salles, vous pouvez **remonter au camp** pour vous soigner. Battre un boss débloque l'étage suivant.

### Le combat

Les combats se jouent **chacun son tour**. À votre tour, vous choisissez :

| Action | Effet |
|---|---|
| **Attaquer** | Frapper l'ennemi |
| **Sorts** | Lancer un sort (coûte du **mana**) |
| **Inventaire** | Boire une potion ou lancer du poison |
| **Fuir** | Tenter de s'échapper (1 chance sur 2, impossible contre un boss) |

Certains ennemis infligent des **effets** :

| Effet | Conséquence | Qui l'inflige |
|---|---|---|
| 🩸 **Saignement** | Perte de PV à chaque tour | Loup |
| 🔥 **Brûlure** | Perte de PV à chaque tour | Krokmou, Ignarok *(et votre Boule de Feu !)* |
| 💫 **Étourdissement** | Vous passez votre tour | Troll |
| 💀 **Affaiblissement** | Vos dégâts sont divisés par 2 | Squelette, Liche |

À savoir :
- ❤️ Si vos **PV** tombent à 0, vous ressuscitez au camp… mais vous perdez **20 % de vos Y-Coins**.
- ⚡ Les monstres frappent **deux fois plus fort tous les 3 tours**.
- ⭐ Vous avez **10 % de chance** de faire un **coup critique** (dégâts doublés).

Chaque victoire rapporte de l'**expérience** (pour monter de niveau), des **Y-Coins** et parfois des **ressources**.

### Le marchand

| Objet | Prix | Effet |
|---|:---:|---|
| Potion de vie | 8 Y-Coins | +50 PV *(la première est offerte !)* |
| Potion de mana | 10 Y-Coins | +30 mana |
| Potion de poison | 15 Y-Coins | 10 dégâts par seconde pendant 3 secondes |
| Livre de Sort : Boule de Feu | 60 Y-Coins | Apprend le sort Boule de Feu |
| Fourrure, Peau, Cuir, Plume | 3 à 18 Y-Coins | Ressources pour le forgeron |
| Augmentation d'inventaire | 75 Y-Coins | +10 places dans le sac |

### Le forgeron

Apportez-lui des **ressources** (laissées par les monstres ou achetées au marchand) :

| Équipement | Bonus | Ressources |
|---|:---:|---|
| Chapeau de l'aventurier | +10 PV | 1 Plume de Corbeau, 1 Cuir de Sanglier |
| Tunique de l'aventurier | +25 PV | 2 Fourrures de Loup, 1 Peau de Troll |
| Bottes de l'aventurier | +15 PV | 1 Fourrure de Loup, 1 Cuir de Sanglier |
| Épée / Arc / Marteau de l'aventurier | +3 attaque | Cuir, plumes, fourrure ou peau de troll |

Il connaît aussi des **recettes de maître**, plus puissantes. Une arme faite pour **votre classe** donne **+2 attaque** en bonus. L'équipement se porte depuis l'**inventaire**.

### Les missions

Le tableau des missions du camp propose **5 missions**, une à la fois (par exemple « Vaincre 3 loups »). Acceptez-la, remplissez l'objectif dans le donjon, puis revenez au tableau pour toucher votre **récompense**.

### 💡 Conseils pour débuter

1. Allez chez le **marchand** récupérer votre **potion gratuite**.
2. Faites un **entraînement** pour comprendre les combats.
3. Passez au **tableau des missions** : acceptez une mission pour gagner des Y-Coins en plus.
4. Au donjon, **remontez au camp** entre deux salles si vos PV sont bas.
5. Gardez des **potions** pour les boss !

### 🧪 Mode test

Pour découvrir tout le jeu rapidement, appelez votre héros **`test`** : PV et mana infinis, 999 999 Y-Coins, attaque à 999, tous les sorts et les 3 étages débloqués.

---

## 🗂️ Organisation du projet

```
DunGo/
├── src/                Le code du jeu
│   ├── main.go         Écran titre et menu du camp
│   ├── character.go    Le héros (création, difficulté, fiche, niveaux)
│   ├── inventory.go    L'inventaire et les potions
│   ├── equipment.go    Les armures et les armes
│   ├── merchant.go     Le marchand
│   ├── forge.go        Le forgeron
│   ├── missions.go     Les missions
│   ├── monster.go      Les monstres et les boss
│   ├── combat.go       Les combats
│   ├── status.go       Les effets (saignement, brûlure…)
│   ├── dungeon.go      Le donjon
│   ├── events.go       Fontaine, marchand ambulant et sphinx
│   ├── save.go         Les sauvegardes
│   ├── utils.go        Couleurs, menus et affichage
│   └── ascii.go        Tous les dessins
├── docs/               Documentation
│   ├── comment-ca-marche.md   Comment le code fonctionne
│   └── gestion-de-projet.md   Suivi du projet
└── .vscode/            Réglages VS Code (touche F5)
```

Les données du jeu sont rangées dans des **tableaux** faciles à modifier : `classes` (lignées), `gear` (armures et armes), `bestiary` (monstres), `spells` (sorts), `recipes` (forge), `quests` (missions) et `floors` (étages).

Le jeu utilise uniquement **Go**, sans rien d'autre à installer.

---
