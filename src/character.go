package main

import (
	"fmt"
	"strings"
	"unicode"
)

// Character est le héros. C'est la seule structure sauvegardée sur disque.
type Character struct {
	Name         string
	Class        string
	Difficulty   string
	Level        int
	XP           int
	XPMax        int
	HP           int
	MaxHP        int
	Mana         int
	MaxMana      int
	Attack       int
	Gold         int
	Skills       []string
	Inventory    map[string]int // objet → quantité
	InventoryMax int
	Equipment    map[string]string // emplacement → objet porté

	InventoryUpgrades int
	FreePotionTaken   bool
	FloorsCleared     int
	Victories         int
	Deaths            int
	SaveSlot          int

	QuestIndex    int
	QuestActive   bool
	QuestProgress int

	// Effets de combat : nombre de tours (ou de coups) restants.
	Bleeding  int
	Burning   int
	Stunned   int
	Weakened  int
	StoneSkin int
}

const (
	SpellPunch       = "Coup de poing"
	SpellFireball    = "Boule de Feu"
	SpellSecondWind  = "Second Souffle"
	SpellSilverArrow = "Flèche d'Argent"
	SpellStoneSkin   = "Peau de Pierre"
)

type Class struct {
	Name     string
	MaxHP    int
	Mana     int
	Spell    string
	Summary  string
	Color    string
	Portrait string
	HPGain   int // gagnés à chaque niveau
	ManaGain int
}

var classes = []Class{
	{"Humain", 100, 40, SpellSecondWind, "Polyvalent, se soigne", Sky, artHuman, 10, 5},
	{"Elfe", 80, 60, SpellSilverArrow, "Magie et flèches", Green, artElf, 6, 8},
	{"Nain", 120, 30, SpellStoneSkin, "Un mur, mais barbu", Orange, artDwarf, 12, 3},
}

func classOf(name string) Class {
	for _, class := range classes {
		if class.Name == name {
			return class
		}
	}
	return classes[0]
}

type Difficulty struct {
	Name   string
	Detail string
	Color  string
	Power  int // % appliqué aux PV et à l'attaque des monstres
	Reward int // % appliqué à l'XP et aux Y-Coins gagnés
}

var difficulties = []Difficulty{
	{"Facile", "Monstres affaiblis (-25 % de PV et d'attaque)", Green, 75, 100},
	{"Normal", "L'aventure telle qu'on l'a imaginée", Gold, 100, 100},
	{"Difficile", "Monstres +30 %, mais +50 % d'XP et de Y-Coins", Red, 130, 150},
}

func difficultyOf(name string) Difficulty {
	for _, level := range difficulties {
		if level.Name == name {
			return level
		}
	}
	return difficulties[1]
}

func newCharacter(name string, class Class) *Character {
	return &Character{
		Name:         name,
		Class:        class.Name,
		Level:        1,
		XPMax:        50,
		HP:           class.MaxHP / 2,
		MaxHP:        class.MaxHP,
		Mana:         class.Mana,
		MaxMana:      class.Mana,
		Attack:       5,
		Gold:         50,
		Skills:       []string{SpellPunch, class.Spell},
		Inventory:    map[string]int{ItemHealthPotion: 1},
		InventoryMax: 10,
		Equipment:    map[string]string{},
	}
}

func characterCreation() *Character {
	clearScreen()
	banner("NOUVEAU HÉROS", Gold)
	fmt.Println()
	say("Le maire", Gold, "Un volontaire contre le dragon ? Merveilleux ! Votre nom, s'il vous plaît. C'est pour… la plaque commémorative.")
	fmt.Println(Gray + "  (lettres uniquement, 12 maximum)" + Reset)

	name := readLine()
	for !isValidName(name) {
		fail("Des lettres uniquement, 12 maximum.")
		name = readLine()
	}
	name = formatName(name)

	c := newCharacter(name, chooseClass(name))
	c.Difficulty = chooseDifficulty()
	return c
}

func chooseClass(name string) Class {
	clearScreen()
	banner("CHOISISSEZ VOTRE CLASSE, "+strings.ToUpper(name), Gold)
	fmt.Println()
	for i, class := range classes {
		option(i+1, fmt.Sprintf("%s%-7s%s %s♥ %3d PV%s  %s♦ %2d mana%s  %s✦ %-15s%s %s%s%s",
			class.Color+Bold, class.Name, Reset,
			Red, class.MaxHP, Reset,
			Blue, class.Mana, Reset,
			Purple, class.Spell, Reset,
			Gray, class.Summary, Reset))
	}
	return classes[readChoice(1, len(classes))-1]
}

func chooseDifficulty() string {
	clearScreen()
	banner("CHOISISSEZ LA DIFFICULTÉ", Gold)
	fmt.Println()
	for i, level := range difficulties {
		option(i+1, level.Color+padRight(level.Name, 12)+Reset+Gray+level.Detail+Reset)
	}
	return difficulties[readChoice(1, len(difficulties))-1].Name
}

func isValidName(name string) bool {
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

// formatName met une majuscule au début et des minuscules ensuite : « pACO » → « Paco ».
func formatName(name string) string {
	letters := []rune(strings.ToLower(name))
	letters[0] = unicode.ToUpper(letters[0])
	return string(letters)
}

func showStatus(c *Character) {
	fmt.Println(DarkGray + "  ╭" + strings.Repeat("─", 72) + Reset)
	fmt.Printf("  %s│%s %s%s%s · %s%s%s   %s¤ %d Y-Coins%s   %s■ Sac %d/%d%s\n",
		DarkGray, Reset,
		Gold+Bold, c.Name, Reset,
		classOf(c.Class).Color, c.Class, Reset,
		Gold, c.Gold, Reset,
		Gray, inventoryCount(c), c.InventoryMax, Reset)
	fmt.Printf("  %s│%s %s♥%s %s %3d/%-3d  %s♦%s %s %3d/%-3d  %s★ Niv.%d%s %s %d/%d\n",
		DarkGray, Reset,
		Red, Reset, hpBar(c.HP, c.MaxHP, 14), c.HP, c.MaxHP,
		Blue, Reset, bar(c.Mana, c.MaxMana, 10, Blue), c.Mana, c.MaxMana,
		Purple+Bold, c.Level, Reset, bar(c.XP, c.XPMax, 10, Purple), c.XP, c.XPMax)
	fmt.Println(DarkGray + "  ╰" + strings.Repeat("─", 72) + Reset)
}

func displayInfo(c *Character) {
	clearScreen()
	class := classOf(c.Class)
	banner("FICHE DU HÉROS", Gold)
	printArt(class.Portrait, class.Color)

	section(c.Name + " · " + c.Class)
	fmt.Printf("   %s♥ PV        %s %s %d / %d\n", Red, Reset, hpBar(c.HP, c.MaxHP, 25), c.HP, c.MaxHP)
	fmt.Printf("   %s♦ Mana      %s %s %d / %d\n", Blue, Reset, bar(c.Mana, c.MaxMana, 25, Blue), c.Mana, c.MaxMana)
	fmt.Printf("   %s★ Niveau %-3d%s %s %d / %d XP\n", Purple+Bold, c.Level, Reset, bar(c.XP, c.XPMax, 25, Purple), c.XP, c.XPMax)
	fmt.Printf("   %s» Attaque   %s %d dégâts\n", Orange, Reset, c.Attack)
	fmt.Printf("   %s¤ Bourse    %s %d Y-Coins\n", Gold, Reset, c.Gold)

	section("Équipement")
	for _, slot := range equipmentSlots {
		item := c.Equipment[slot]
		if item == "" {
			item = DarkGray + "— rien —" + Reset
		}
		fmt.Printf("   %s%-7s%s %s\n", Silver, slot, Reset, item)
	}

	section("Sorts")
	for _, name := range c.Skills {
		spell := spells[name]
		fmt.Printf("   %s✦ %-16s%s %s(%d mana)%s %s%s%s\n", Purple, name, Reset, Blue, spell.Cost, Reset, Gray, spell.Description, Reset)
	}

	section("Progression")
	fmt.Printf("   Difficulté : %s   ·   Étages terminés : %d / 3\n", c.Difficulty, c.FloorsCleared)
	fmt.Printf("   Monstres vaincus : %d   ·   Morts : %d   ·   Missions : %d / %d\n", c.Victories, c.Deaths, c.QuestIndex, len(quests))
	pause()
}

// gainXP ajoute de l'expérience et fait monter de niveau si besoin.
// L'XP en trop est gardée pour le niveau suivant.
func gainXP(c *Character, amount int) {
	c.XP += amount
	levelsGained := 0
	for c.XP >= c.XPMax {
		c.XP -= c.XPMax
		c.XPMax = c.XPMax * 3 / 2
		c.Level++
		levelsGained++
	}
	fmt.Printf("  %s★ +%d XP   Niv.%d%s %s %d / %d\n", Purple+Bold, amount, c.Level, Reset, bar(c.XP, c.XPMax, 20, Purple), c.XP, c.XPMax)
	if levelsGained == 0 {
		return
	}

	class := classOf(c.Class)
	c.MaxHP += class.HPGain * levelsGained
	c.MaxMana += class.ManaGain * levelsGained
	c.Attack += levelsGained
	c.HP = c.MaxHP
	c.Mana = c.MaxMana

	fmt.Println()
	printArt(artLevelUp, campColors...)
	section(fmt.Sprintf("Niveau %d !", c.Level))
	fmt.Printf("   %s♥ +%d PV max   %s♦ +%d mana max   %s» +%d attaque%s\n",
		Red, class.HPGain*levelsGained, Blue, class.ManaGain*levelsGained, Orange, levelsGained, Reset)
	success("PV et mana à fond.")
}

// isDead ressuscite le héros avec la moitié de ses PV s'il est tombé à 0.
func isDead(c *Character) bool {
	if c.HP > 0 {
		return false
	}
	c.Deaths++
	c.HP = c.MaxHP / 2
	clearScreen()
	printArt(artTomb, stoneColors...)
	fmt.Println(Gray + "        ci-gît " + c.Name + ", un héros pressé" + Reset)
	fmt.Println()
	paragraph(Gray, "Les dieux de DunGo vous renvoient au camp. Ils trouvent l'histoire trop drôle pour qu'elle s'arrête là.")
	success("Vous ressuscitez avec %d / %d PV.", c.HP, c.MaxHP)
	return true
}
