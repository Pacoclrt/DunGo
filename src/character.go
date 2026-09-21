package main

import (
	"fmt"
	"strings"
	"unicode"
)

// Character est le héros. C'est la seule structure sauvegardée sur disque :
// tout ce qui doit survivre à un « Continuer » vit ici.
type Character struct {
	Name         string
	Class        string
	Level        int
	MaxHP        int
	HP           int
	MaxMana      int
	Mana         int
	Attack       int
	XP           int
	XPMax        int
	Gold         int
	Skills       []string
	Equipment    Equipment
	Inventory    map[string]int
	InventoryMax int

	InventoryUpgrades int
	FreePotionTaken   bool
	FloorsCleared     int
	DragonSlain       bool
	Victories         int
	Deaths            int
	GoldEarned        int
	TestMode          bool
	Difficulty        string // « Facile », « Normal » ou « Difficile »
	SaveSlot          int

	QuestIndex      int
	QuestActive     bool
	QuestProgress   int
	QuestsCompleted int

	// Effets temporaires, remis à zéro à chaque combat.
	Bleeding  int // tours restants
	Burning   int
	Stunned   int
	Weakened  int // attaques affaiblies restantes
	StoneSkin int // attaques encore absorbées
}

const (
	SpellPunch       = "Coup de poing"
	SpellFireball    = "Boule de Feu"
	SpellSecondWind  = "Second Souffle"
	SpellSilverArrow = "Flèche d'Argent"
	SpellStoneSkin   = "Peau de Pierre"
)

// knows dit si le héros connaît déjà un sort.
func (c *Character) knows(spell string) bool {
	for _, known := range c.Skills {
		if known == spell {
			return true
		}
	}
	return false
}

// canCast dit si le héros a assez de mana (toujours vrai en mode test).
func (c *Character) canCast(spell string) bool {
	return c.TestMode || c.Mana >= spells[spell].Cost
}

// ----------------------------------------------------------------- classes

// Class décrit une classe jouable et sa progression par niveau.
type Class struct {
	Name     string
	MaxHP    int
	Mana     int
	Spell    string
	Summary  string
	Color    string
	Portrait string
	HPGain   int
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

// --------------------------------------------------------------- création

func newCharacter(name string, class Class) *Character {
	return &Character{
		Name:         name,
		Class:        class.Name,
		Level:        1,
		MaxHP:        class.MaxHP,
		HP:           class.MaxHP / 2,
		MaxMana:      class.Mana,
		Mana:         class.Mana,
		Attack:       5,
		XPMax:        50,
		Gold:         50,
		Skills:       []string{SpellPunch, class.Spell},
		Inventory:    map[string]int{ItemHealthPotion: 1},
		InventoryMax: 10,
	}
}

func characterCreation() *Character {
	clearScreen()
	banner("NOUVEAU HÉROS", Gold)
	fmt.Println()
	say("Le maire", Gold, "Un volontaire contre le dragon ? Merveilleux ! Votre nom, s'il vous plaît. C'est pour… euh… la plaque commémorative.")
	fmt.Println(DarkGray + "  (lettres uniquement, 12 maximum)" + Reset)

	name := readLine()
	for !isValidName(name) {
		fail("Des lettres uniquement, 12 maximum. Même les gobelins y arrivent.")
		name = readLine()
	}
	name = formatName(name)
	success("Bienvenue, %s !", name)
	wait(700)

	for {
		class := chooseClass(name)
		clearScreen()
		printArtSlow(class.Portrait, class.Color)
		fmt.Println()
		fmt.Printf("  %s%s%s, %s%s%s. Est-ce bien vous ?\n", Bold+Gold, name, Reset, class.Color, class.Name, Reset)
		if !ask("Confirmer ?") {
			continue
		}
		c := newCharacter(name, class)
		c.Difficulty = chooseDifficulty()
		if strings.EqualFold(name, "test") {
			activateTestMode(c)
		}
		return c
	}
}

// chooseClass affiche les trois portraits côte à côte avec leurs statistiques.
func chooseClass(name string) Class {
	clearScreen()
	banner("CHOISISSEZ VOTRE CLASSE, "+strings.ToUpper(name), Gold)
	fmt.Println()

	const column = 24
	rowColors := []string{Bold, Red, Blue, Purple, Gray}
	portraits, colors := []string{}, []string{}
	rows := make([]string, len(rowColors))
	for i, class := range classes {
		portraits = append(portraits, class.Portrait)
		colors = append(colors, class.Color)
		texts := []string{
			fmt.Sprintf("[%d] %s", i+1, strings.ToUpper(class.Name)),
			fmt.Sprintf("♥ %d PV", class.MaxHP),
			fmt.Sprintf("♦ %d mana", class.Mana),
			"✦ " + class.Spell,
			"  " + class.Summary,
		}
		for row, text := range texts {
			color := rowColors[row]
			if row == 0 {
				color = class.Color + Bold
			}
			// On complète le texte nu à la bonne largeur AVANT de le colorer :
			// les codes ANSI ne prennent aucune place à l'écran.
			rows[row] += color + padRight(text, column) + Reset
		}
	}
	sideBySide(portraits, colors, column)
	fmt.Println()
	for _, row := range rows {
		fmt.Println("   " + row)
	}
	return classes[readChoice(1, len(classes))-1]
}

// Difficulty règle la puissance des monstres et les gains du joueur.
// Les pourcentages servent à la fois au texte du menu et au calcul réel :
// impossible que les deux se contredisent.
type Difficulty struct {
	Name   string
	Detail string
	Color  string
	Power  int // % appliqué aux PV et à l'attaque du monstre
	Reward int // % appliqué à l'XP et aux Y-Coins
}

var difficulties = []Difficulty{
	{"Facile", "Monstres affaiblis (-25 % de PV et d'attaque). Pas de honte.", Green, 75, 100},
	{"Normal", "L'aventure telle qu'on l'a imaginée", Gold, 100, 100},
	{"Difficile", "Monstres +30 %, mais +50 % d'XP et de Y-Coins. Courage.", Red, 130, 150},
}

// difficultyOf retrouve une difficulté par son nom ; Normal par défaut.
func difficultyOf(name string) Difficulty {
	for _, level := range difficulties {
		if level.Name == name {
			return level
		}
	}
	return difficulties[1]
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

func activateTestMode(c *Character) {
	c.TestMode = true
	c.HP = c.MaxHP
	c.Gold = 999999
	c.Attack = 999
	c.InventoryMax = 999
	c.Skills = []string{SpellPunch, SpellFireball, SpellSecondWind, SpellSilverArrow, SpellStoneSkin}
	c.FloorsCleared = 2

	clearScreen()
	banner("MODE TEST ACTIVÉ", Gold)
	fmt.Println()
	for _, effect := range []string{
		"PV infinis : vous ne subissez aucun dégât",
		"Mana infini : les sorts ne coûtent rien",
		"999 999 Y-Coins et un sac de 999 places",
		"Attaque à 999 : chaque coup est mortel",
		"Tous les sorts connus et les 3 étages débloqués",
	} {
		info("%s", effect)
	}
	fmt.Println(Gray + Italic + "\n  Le dragon trouve ça un peu injuste." + Reset)
	pause()
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

func formatName(name string) string {
	letters := []rune(strings.ToLower(name))
	letters[0] = unicode.ToUpper(letters[0])
	return string(letters)
}

// ------------------------------------------------------------------ fiche

// showStatus est le bandeau compact affiché en haut de la plupart des écrans.
// Le niveau est accolé à la barre d'XP : c'est elle qui le fait monter.
func showStatus(c *Character) {
	tags := ""
	if c.TestMode {
		tags = Gold + Bold + "   ∞ MODE TEST" + Reset
	}
	fmt.Println(DarkGray + "  ╭" + strings.Repeat("─", 72) + Reset)
	fmt.Printf("  %s│%s %s%s%s · %s%s%s   %s¤ %d Y-Coins%s   %s■ Sac %d/%d%s%s\n",
		DarkGray, Reset, Bold+Gold, c.Name, Reset, classOf(c.Class).Color, c.Class, Reset,
		Gold, c.Gold, Reset, Gray, inventoryCount(c), c.InventoryMax, Reset, tags)
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
	printArt(class.Portrait, class.Color, class.Color, class.Color, Silver)

	section(c.Name + " · " + c.Class)
	fmt.Printf("   %s♥ PV        %s %s %d / %d\n", Red, Reset, hpBar(c.HP, c.MaxHP, 25), c.HP, c.MaxHP)
	fmt.Printf("   %s♦ Mana      %s %s %d / %d\n", Blue, Reset, bar(c.Mana, c.MaxMana, 25, Blue), c.Mana, c.MaxMana)
	fmt.Printf("   %s★ Niveau %-3d%s %s %d / %d XP\n", Purple+Bold, c.Level, Reset, bar(c.XP, c.XPMax, 25, Purple), c.XP, c.XPMax)
	fmt.Printf("   %s» Attaque   %s %d dégâts\n", Orange, Reset, c.Attack)
	fmt.Printf("   %s¤ Bourse    %s %d Y-Coins    %s■ Sac%s %d / %d\n", Gold, Reset, c.Gold, Gray, Reset, inventoryCount(c), c.InventoryMax)

	section("Équipement")
	for _, slot := range equipmentSlots {
		fmt.Printf("   %s%-7s%s %s\n", Silver, slot, Reset, equipmentLabel(c, *c.slot(slot)))
	}

	section("Sorts")
	for _, name := range c.Skills {
		spell := spells[name]
		fmt.Printf("   %s✦ %-16s%s %s(%d mana)%s %s\n", Purple, name, Reset, Blue, spell.Cost, Reset, Gray+spell.Description+Reset)
	}

	section("Progression")
	fmt.Printf("   Difficulté : %s   ·   Donjon : %d / 3 étages\n", difficultyOf(c.Difficulty).Name, c.FloorsCleared)
	fmt.Printf("   Monstres vaincus : %d   ·   Morts : %d   ·   Missions : %d / %d\n", c.Victories, c.Deaths, c.QuestsCompleted, len(quests))
	pause()
}

func equipmentLabel(c *Character, item string) string {
	switch {
	case item == "":
		return DarkGray + "— rien, à part du courage —" + Reset
	case isWeapon(item):
		return fmt.Sprintf("%s%s%s%s (+%d attaque)%s", White, Bold, item, Orange, weaponBonus(c, item), Reset)
	}
	return fmt.Sprintf("%s%s%s%s (+%d PV)%s", White, Bold, item, Red, gear[item].HP, Reset)
}

// ------------------------------------------------------------ mort et niveau

// isDead ressuscite le héros au camp s'il est tombé. Renvoie true s'il est mort.
func isDead(c *Character) bool {
	if c.HP > 0 {
		return false
	}
	c.Deaths++
	clearScreen()
	printArtSlow(artWasted, bloodColors...)
	fmt.Println()
	printArt(artTomb, stoneColors...)
	fmt.Println(Bold + White + centerText("~ "+strings.ToUpper(c.Name)+" ~", 46) + Reset)
	fmt.Println(Gray + centerText("ci-gît un héros pressé", 46) + Reset)
	fmt.Println()
	typewrite(Gray, "Les dieux de DunGo vous renvoient au camp. Ils trouvent l'histoire trop drôle pour qu'elle s'arrête là.")
	c.HP = c.MaxHP / 2
	success("Vous ressuscitez avec %d / %d PV.", c.HP, c.MaxHP)
	return true
}

// gainXP ajoute de l'expérience et fait monter de niveau autant que nécessaire.
func gainXP(c *Character, amount int) {
	c.XP += amount
	class, gained := classOf(c.Class), 0
	for c.XP >= c.XPMax {
		c.XP -= c.XPMax // l'excédent est conservé pour le niveau suivant
		c.XPMax = c.XPMax * 3 / 2
		c.Level++
		c.MaxHP += class.HPGain
		c.MaxMana += class.ManaGain
		c.Attack++
		gained++
	}
	fmt.Printf("  %s★ +%d XP%s   %sNiv.%d%s %s %d / %d\n", Purple+Bold, amount, Reset,
		Purple+Bold, c.Level, Reset, bar(c.XP, c.XPMax, 20, Purple), c.XP, c.XPMax)
	if gained == 0 {
		return
	}

	c.HP, c.Mana = c.MaxHP, c.MaxMana
	fmt.Println()
	printArtSlow(artLevelUp, campColors...)
	section(fmt.Sprintf("Niveau %d !", c.Level))
	fmt.Printf("   %s♥ +%d PV max   %s♦ +%d mana max   %s» +%d attaque%s\n",
		Red, class.HPGain*gained, Blue, class.ManaGain*gained, Orange, gained, Reset)
	success("PV et mana à fond. Ça fait du bien.")
}
