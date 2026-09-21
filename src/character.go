// ════════════════════════════════════════════════════════════════════════
//   character.go · [HEROS]
//   Le héros : sa structure (la seule chose sauvegardée), les classes, les
//   difficultés, la création du personnage, sa fiche, l'XP et les niveaux,
//   et la mort.
// ════════════════════════════════════════════════════════════════════════

package main

import (
	"fmt"
	"strings"
	"unicode"
)

// ════════════════════════════════════════════════════════════════════════
//   LE HÉROS
// ════════════════════════════════════════════════════════════════════════

// Tous les champs commencent par une majuscule : c'est obligatoire pour que
// le paquet encoding/json puisse les lire et les écrire (save.go).
type Character struct { // [HEROS] sert à décrire le héros : tout ce qui doit être sauvegardé est ici, et nulle part ailleurs
	// Qui il est
	Name       string // son nom, ex : « Paco »
	Class      string // "Humain", "Elfe" ou "Nain"
	Difficulty string // "Facile", "Normal" ou "Difficile"
	SaveSlot   int    // l'emplacement de sauvegarde : 1, 2 ou 3

	// Ses statistiques
	Level       int // son niveau (commence à 1)
	XP          int // l'XP gagnée depuis le dernier niveau
	XPToLevelUp int // l'XP à atteindre pour passer au niveau suivant
	HP          int // PV = points de vie
	MaxHP       int // PV maximum
	Mana        int // sert à lancer des sorts
	MaxMana     int // mana maximum
	Attack      int // les dégâts de l'attaque de base
	YCoins      int // la monnaie du jeu

	// Ce qu'il possède
	Spells            []string          // les noms des sorts qu'il connaît
	Inventory         map[string]int    // le sac : nom de l'objet → quantité
	InventoryMax      int               // le nombre d'objets que le sac peut contenir
	InventoryUpgrades int               // combien de fois le sac a été agrandi (3 au maximum)
	Equipment         map[string]string // ce qu'il porte : emplacement ("Tête"…) → nom de l'objet
	FreePotionTaken   bool              // true dès qu'il a pris la Potion de vie offerte par le marchand

	// Sa progression
	FloorsCleared int // le nombre d'étages terminés (0 à 3) : l'étage suivant est le seul débloqué
	Victories     int // le nombre de monstres vaincus
	Deaths        int // le nombre de morts

	// Sa mission (quests.go)
	QuestIndex    int  // le numéro de la mission en cours dans la liste quests (= missions déjà terminées)
	QuestActive   bool // true si la mission en cours a été acceptée
	QuestProgress int  // le nombre de monstres cibles déjà vaincus

	// Les effets de combat (effects.go) : le nombre de tours (ou de coups)
	// qu'il leur reste. 0 = pas d'effet. Tout est remis à 0 au DÉBUT de
	// chaque combat (un reste d'effet après un combat ne fait donc rien).
	Bleeding  int // saignement : -3 PV à la fin de chaque tour
	Burning   int // brûlure : -5 PV à la fin de chaque tour
	Stunned   int // étourdissement : le héros passe son prochain tour
	Weakened  int // affaiblissement : ses prochaines attaques font moitié moins mal
	StoneSkin int // Peau de Pierre : les prochains coups reçus font moitié moins mal
}

// ════════════════════════════════════════════════════════════════════════
//   LES CLASSES
// ════════════════════════════════════════════════════════════════════════

type Class struct { // [HEROS] sert à décrire une classe jouable : ses stats de départ, son sort et ce qu'elle gagne à chaque niveau
	Name         string
	MaxHP        int
	MaxMana      int
	Spell        string // le sort de départ propre à la classe
	Summary      string // la petite phrase affichée au choix de la classe
	Color        string
	Portrait     string // le dessin de la fiche du héros
	HPPerLevel   int    // PV max gagnés à chaque niveau
	ManaPerLevel int    // mana max gagné à chaque niveau
}

var classes = []Class{ // [HEROS] sert à lister les 3 classes jouables, dans l'ordre du menu
	{Name: "Humain", MaxHP: 100, MaxMana: 40, Spell: SpellSecondWind, Summary: "Polyvalent, se soigne",
		Color: Sky, Portrait: artHuman, HPPerLevel: 10, ManaPerLevel: 5},
	{Name: "Elfe", MaxHP: 80, MaxMana: 60, Spell: SpellSilverArrow, Summary: "Magie et flèches",
		Color: Green, Portrait: artElf, HPPerLevel: 6, ManaPerLevel: 8},
	{Name: "Nain", MaxHP: 120, MaxMana: 30, Spell: SpellStoneSkin, Summary: "Un mur, mais barbu",
		Color: Orange, Portrait: artDwarf, HPPerLevel: 12, ManaPerLevel: 3},
}

func findClass(name string) Class { // [HEROS] sert à retrouver une classe à partir de son nom ("Elfe" → la classe Elfe) et la renvoie
	for _, class := range classes {
		if class.Name == name {
			return class
		}
	}
	return classes[0] // nom inconnu (sauvegarde modifiée à la main…) : Humain par défaut
}

// ════════════════════════════════════════════════════════════════════════
//   LES DIFFICULTÉS
// ════════════════════════════════════════════════════════════════════════

type Difficulty struct { // [HEROS] sert à décrire une difficulté : deux pourcentages appliqués à chaque monstre créé (monster.go)
	Name           string
	Detail         string
	Color          string
	MonsterPercent int // % appliqué aux PV et à l'attaque des monstres
	RewardPercent  int // % appliqué à l'XP et aux Y-Coins qu'ils rapportent
}

var difficulties = []Difficulty{ // [HEROS] sert à lister les 3 difficultés, dans l'ordre du menu
	{Name: "Facile", Detail: "Monstres affaiblis (-25 % de PV et d'attaque)", Color: Green, MonsterPercent: 75, RewardPercent: 100},
	{Name: "Normal", Detail: "L'aventure telle qu'on l'a imaginée", Color: Gold, MonsterPercent: 100, RewardPercent: 100},
	{Name: "Difficile", Detail: "Monstres +30 %, mais +50 % d'XP et de Y-Coins", Color: Red, MonsterPercent: 130, RewardPercent: 150},
}

func findDifficulty(name string) Difficulty { // [HEROS] sert à retrouver une difficulté à partir de son nom et la renvoie (Normal si le nom est inconnu)
	for _, difficulty := range difficulties {
		if difficulty.Name == name {
			return difficulty
		}
	}
	return difficulties[1] // Normal par défaut
}

// ════════════════════════════════════════════════════════════════════════
//   LA CRÉATION DU HÉROS
// ════════════════════════════════════════════════════════════════════════

func createCharacter() *Character { // [HEROS] sert à créer le héros avec le joueur (nom → classe → difficulté) et le renvoie
	clearScreen()
	banner("NOUVEAU HÉROS", Gold)
	fmt.Println()
	say("Le maire", Gold, "Un volontaire contre le dragon ? Merveilleux ! Votre nom, s'il vous plaît. C'est pour… la plaque commémorative.")
	fmt.Println(Gray + "  (lettres uniquement, 12 maximum)" + Reset)

	// 1. Le nom : on redemande tant qu'il n'est pas valide
	name := readLine()
	for !isValidName(name) {
		fail("Des lettres uniquement, 12 maximum.")
		name = readLine()
	}
	name = formatName(name)

	// 2. La classe, puis 3. la difficulté
	character := startingCharacter(name, chooseClass(name))
	character.Difficulty = chooseDifficulty()
	return character
}

func startingCharacter(name string, class Class) *Character { // [HEROS] sert à fabriquer le héros de départ (niveau 1, moitié de ses PV, 50 Y-Coins, 1 potion) et le renvoie
	return &Character{
		Name:         name,
		Class:        class.Name,
		Level:        1,
		XPToLevelUp:  50,
		HP:           class.MaxHP / 2, // il commence blessé : la potion offerte par le marchand tombe bien !
		MaxHP:        class.MaxHP,
		Mana:         class.MaxMana,
		MaxMana:      class.MaxMana,
		Attack:       5,
		YCoins:       50,
		Spells:       []string{SpellPunch, class.Spell},
		Inventory:    map[string]int{ItemHealthPotion: 1},
		InventoryMax: 10,
		Equipment:    map[string]string{},
	}
}

func chooseClass(name string) Class { // [HEROS] sert à afficher les 3 classes et renvoie celle que le joueur choisit
	clearScreen()
	banner("CHOISISSEZ VOTRE CLASSE, "+strings.ToUpper(name), Gold)
	fmt.Println()
	for index, class := range classes {
		option(index+1, fmt.Sprintf("%s%-7s%s %s♥ %3d PV%s  %s♦ %2d mana%s  %s✦ %-15s%s %s%s%s",
			class.Color+Bold, class.Name, Reset,
			Red, class.MaxHP, Reset,
			Blue, class.MaxMana, Reset,
			Purple, class.Spell, Reset,
			Gray, class.Summary, Reset))
	}
	// Le joueur tape 1, 2 ou 3, mais une liste commence à 0 : d'où le -1.
	return classes[readChoice(1, len(classes))-1]
}

func chooseDifficulty() string { // [HEROS] sert à afficher les 3 difficultés et renvoie le nom de celle choisie
	clearScreen()
	banner("CHOISISSEZ LA DIFFICULTÉ", Gold)
	fmt.Println()
	for index, difficulty := range difficulties {
		option(index+1, difficulty.Color+padRight(difficulty.Name, 12)+Reset+Gray+difficulty.Detail+Reset)
	}
	return difficulties[readChoice(1, len(difficulties))-1].Name
}

func isValidName(name string) bool { // [HEROS] sert à vérifier un nom : entre 1 et 12 lettres, rien d'autre (ni chiffre, ni espace) ; renvoie true s'il est valide
	// []rune(name) découpe le texte en lettres : « Zoé » fait 3 lettres
	// (alors que len("Zoé") vaut 4 octets, car « é » prend 2 octets).
	if name == "" || len([]rune(name)) > 12 {
		return false
	}
	for _, letter := range name {
		if !unicode.IsLetter(letter) {
			return false
		}
	}
	return true
}

func formatName(name string) string { // [HEROS] sert à mettre une majuscule au début et des minuscules ensuite (« pACO » → « Paco ») et renvoie le nom
	letters := []rune(strings.ToLower(name))
	letters[0] = unicode.ToUpper(letters[0])
	return string(letters)
}

// ════════════════════════════════════════════════════════════════════════
//   L'AFFICHAGE DU HÉROS
// ════════════════════════════════════════════════════════════════════════

func showCharacterBar(character *Character) { // [HEROS] sert à afficher le bandeau du héros (nom, Y-Coins, sac, PV, mana, niveau) en haut des écrans du camp
	fmt.Println(DarkGray + "  ╭" + strings.Repeat("─", 72) + Reset)
	fmt.Printf("  %s│%s %s%s%s · %s%s%s   %s¤ %d Y-Coins%s   %s■ Sac %d/%d%s\n",
		DarkGray, Reset,
		Gold+Bold, character.Name, Reset,
		findClass(character.Class).Color, character.Class, Reset,
		Gold, character.YCoins, Reset,
		Gray, inventoryCount(character), character.InventoryMax, Reset)
	fmt.Printf("  %s│%s %s♥%s %s %3d/%-3d  %s♦%s %s %3d/%-3d  %s★ Niv.%d%s %s %d/%d\n",
		DarkGray, Reset,
		Red, Reset, hpBar(character.HP, character.MaxHP, 14), character.HP, character.MaxHP,
		Blue, Reset, bar(character.Mana, character.MaxMana, 10, Blue), character.Mana, character.MaxMana,
		Purple+Bold, character.Level, Reset, bar(character.XP, character.XPToLevelUp, 10, Purple), character.XP, character.XPToLevelUp)
	fmt.Println(DarkGray + "  ╰" + strings.Repeat("─", 72) + Reset)
}

func showCharacterSheet(character *Character) { // [HEROS] sert à afficher la fiche complète du héros : stats, équipement, sorts et progression
	clearScreen()
	class := findClass(character.Class)
	banner("FICHE DU HÉROS", Gold)
	printArt(class.Portrait, class.Color)

	section(character.Name + " · " + character.Class)
	fmt.Printf("   %s♥ PV        %s %s %d / %d\n", Red, Reset, hpBar(character.HP, character.MaxHP, 25), character.HP, character.MaxHP)
	fmt.Printf("   %s♦ Mana      %s %s %d / %d\n", Blue, Reset, bar(character.Mana, character.MaxMana, 25, Blue), character.Mana, character.MaxMana)
	fmt.Printf("   %s★ Niveau %-3d%s %s %d / %d XP\n", Purple+Bold, character.Level, Reset, bar(character.XP, character.XPToLevelUp, 25, Purple), character.XP, character.XPToLevelUp)
	fmt.Printf("   %s» Attaque   %s %d dégâts\n", Orange, Reset, character.Attack)
	fmt.Printf("   %s¤ Bourse    %s %d Y-Coins\n", Gold, Reset, character.YCoins)

	section("Équipement")
	for _, slot := range equipmentSlots {
		item := character.Equipment[slot]
		if item == "" {
			item = DarkGray + "— rien —" + Reset
		}
		fmt.Printf("   %s%-7s%s %s\n", Silver, slot, Reset, item)
	}

	section("Sorts")
	for _, name := range character.Spells {
		spell := spells[name]
		fmt.Printf("   %s✦ %-16s%s %s(%d mana)%s %s%s%s\n", Purple, name, Reset, Blue, spell.ManaCost, Reset, Gray, spell.Description, Reset)
	}

	section("Progression")
	fmt.Printf("   Difficulté : %s   ·   Étages terminés : %d / %d\n", character.Difficulty, character.FloorsCleared, len(floors))
	fmt.Printf("   Monstres vaincus : %d   ·   Morts : %d   ·   Missions : %d / %d\n", character.Victories, character.Deaths, character.QuestIndex, len(quests))
	pause()
}

// ════════════════════════════════════════════════════════════════════════
//   L'EXPÉRIENCE ET LES NIVEAUX
// ════════════════════════════════════════════════════════════════════════

func gainXP(character *Character, amount int) { // [HEROS] sert à donner de l'XP au héros et à le faire monter de niveau (une ou plusieurs fois) si besoin
	// 1. On ajoute l'XP. Tant qu'elle dépasse le palier, on monte d'un
	//    niveau : l'XP en trop est gardée et le palier suivant est 1,5 fois
	//    plus haut (50 → 75 → 112 → 168…).
	character.XP += amount
	levelsGained := 0
	for character.XP >= character.XPToLevelUp {
		character.XP -= character.XPToLevelUp
		character.XPToLevelUp = character.XPToLevelUp * 3 / 2
		character.Level++
		levelsGained++
	}
	fmt.Printf("  %s★ +%d XP   Niv.%d%s %s %d / %d\n", Purple+Bold, amount, character.Level, Reset, bar(character.XP, character.XPToLevelUp, 20, Purple), character.XP, character.XPToLevelUp)
	if levelsGained == 0 {
		return
	}

	// 2. Chaque niveau gagné donne les bonus de la classe + 1 d'attaque,
	//    puis les PV et le mana sont remis au maximum.
	class := findClass(character.Class)
	character.MaxHP += class.HPPerLevel * levelsGained
	character.MaxMana += class.ManaPerLevel * levelsGained
	character.Attack += levelsGained
	character.HP = character.MaxHP
	character.Mana = character.MaxMana

	fmt.Println()
	printArt(artLevelUp, campColors...)
	section(fmt.Sprintf("Niveau %d !", character.Level))
	fmt.Printf("   %s♥ +%d PV max   %s♦ +%d mana max   %s» +%d attaque%s\n",
		Red, class.HPPerLevel*levelsGained, Blue, class.ManaPerLevel*levelsGained, Orange, levelsGained, Reset)
	success("PV et mana à fond.")
}

// ════════════════════════════════════════════════════════════════════════
//   LA MORT
// ════════════════════════════════════════════════════════════════════════

func characterDies(character *Character) { // [HEROS] sert à gérer la mort du héros : il perd 20 % de ses Y-Coins et ressuscite au camp avec la moitié de ses PV
	character.Deaths++
	lost := character.YCoins / 5 // 1/5 = 20 %
	character.YCoins -= lost
	character.HP = character.MaxHP / 2

	clearScreen()
	printArt(artTomb, stoneColors...)
	fmt.Println(Gray + "        ci-gît " + character.Name + ", un héros pressé" + Reset)
	fmt.Println()
	paragraph(Gray, "Les dieux de DunGo vous renvoient au camp. Ils trouvent l'histoire trop drôle pour qu'elle s'arrête là.")
	success("Vous ressuscitez avec %d / %d PV.", character.HP, character.MaxHP)
	fail("Vous perdez %d Y-Coins en chemin.", lost)
}
