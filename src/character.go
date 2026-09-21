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
	Initiative   int
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

// ----------------------------------------------------------------- lignées

// Class décrit une lignée jouable et sa progression par niveau.
type Class struct {
	Name       string
	MaxHP      int
	Mana       int
	Initiative int
	Spell      string
	Summary    string
	Color      string
	Portrait   string
	HPGain     int
	ManaGain   int
}

var classes = []Class{
	{"Humain", 100, 40, 10, SpellSecondWind, "Polyvalent, se soigne", Sky, artHuman, 10, 5},
	{"Elfe", 80, 60, 14, SpellSilverArrow, "Rapide, magicien", Green, artElf, 6, 8},
	{"Nain", 120, 30, 7, SpellStoneSkin, "Un roc, encaisse tout", Orange, artDwarf, 12, 3},
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
		Initiative:   class.Initiative,
		XPMax:        50,
		Gold:         50,
		Skills:       []string{SpellPunch, class.Spell},
		Inventory:    map[string]int{ItemHealthPotion: 1},
		InventoryMax: 10,
	}
}

func characterCreation() *Character {
	clearScreen()
	banner("REGISTRE DES AVENTURIERS", Gold)
	fmt.Println()
	say("L'intendant", Gold, "Encore un fou qui veut affronter Ignarok ? Soit. Quel est votre nom ?")
	fmt.Println(DarkGray + "  (lettres uniquement, 12 maximum)" + Reset)

	name := readLine()
	for !isValidName(name) {
		fail("Nom invalide : utilisez uniquement des lettres (12 maximum).")
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
		fmt.Printf("  %s%s%s, %s%s%s de niveau 1. Est-ce bien vous ?\n", Bold+Gold, name, Reset, class.Color, class.Name, Reset)
		if !ask("Confirmer la lignée") {
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
	banner("CHOISISSEZ VOTRE LIGNÉE, "+strings.ToUpper(name), Gold)
	fmt.Println()

	const column = 24
	rowColors := []string{Bold, Red, Blue, Yellow, Purple, Gray}
	portraits, colors := []string{}, []string{}
	rows := make([]string, len(rowColors))
	for i, class := range classes {
		portraits = append(portraits, class.Portrait)
		colors = append(colors, class.Color)
		texts := []string{
			fmt.Sprintf("[%d] %s", i+1, strings.ToUpper(class.Name)),
			fmt.Sprintf("♥ %d PV", class.MaxHP),
			fmt.Sprintf("♦ %d mana", class.Mana),
			fmt.Sprintf("» Initiative %d", class.Initiative),
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
	{"Facile", "Monstres affaiblis (-25 % de PV et d'attaque)", Green, 75, 100},
	{"Normal", "L'aventure telle qu'elle a été pensée", Yellow, 100, 100},
	{"Difficile", "Monstres renforcés (+30 %), mais +50 % d'XP et de Y-Coins", Red, 130, 150},
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
		option(i+1, padRight(level.Name, 12)+Gray+level.Detail+Reset, level.Color)
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
	banner("MODE TEST ACTIVÉ", Cyan)
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
func showStatus(c *Character) {
	tags := ""
	if c.TestMode {
		tags = Cyan + Bold + "   ∞ MODE TEST" + Reset
	}
	fmt.Println(DarkGray + "  ╭" + strings.Repeat("─", 70) + Reset)
	fmt.Printf("  %s│%s %s%s%s %s%s niv.%d%s   %s¤ %d Y-Coins%s   %sSac %d/%d%s%s\n",
		DarkGray, Reset, Bold+Gold, c.Name, Reset, classOf(c.Class).Color, c.Class, c.Level, Reset,
		Gold, c.Gold, Reset, Gray, inventoryCount(c), c.InventoryMax, Reset, tags)
	fmt.Printf("  %s│%s %s♥%s %s %3d/%-3d  %s♦%s %s %3d/%-3d  %s★%s %s %d/%d\n",
		DarkGray, Reset,
		Red, Reset, hpBar(c.HP, c.MaxHP, 14), c.HP, c.MaxHP,
		Blue, Reset, bar(c.Mana, c.MaxMana, 10, Blue), c.Mana, c.MaxMana,
		Purple, Reset, bar(c.XP, c.XPMax, 10, Purple), c.XP, c.XPMax)
	fmt.Println(DarkGray + "  ╰" + strings.Repeat("─", 70) + Reset)
}

func displayInfo(c *Character) {
	clearScreen()
	class := classOf(c.Class)
	banner("FICHE DU HÉROS", Gold)
	printArt(class.Portrait, class.Color, class.Color, class.Color, Silver)

	section(fmt.Sprintf("%s · %s de niveau %d", c.Name, c.Class, c.Level), class.Color)
	fmt.Printf("   %s♥ PV        %s %s %d / %d\n", Red, Reset, hpBar(c.HP, c.MaxHP, 25), c.HP, c.MaxHP)
	fmt.Printf("   %s♦ Mana      %s %s %d / %d\n", Blue, Reset, bar(c.Mana, c.MaxMana, 25, Blue), c.Mana, c.MaxMana)
	fmt.Printf("   %s★ Expérience%s %s %d / %d\n", Purple, Reset, bar(c.XP, c.XPMax, 25, Purple), c.XP, c.XPMax)
	fmt.Printf("   %s» Attaque   %s %d dégâts     %s» Initiative%s %d\n", Orange, Reset, c.Attack, Yellow, Reset, c.Initiative)
	fmt.Printf("   %s¤ Bourse    %s %d Y-Coins    %s■ Sac%s %d / %d\n", Gold, Reset, c.Gold, Brown, Reset, inventoryCount(c), c.InventoryMax)

	section("Équipement", Silver)
	for _, slot := range equipmentSlots {
		fmt.Printf("   %s%-7s%s %s\n", Silver, slot, Reset, equipmentLabel(c, *c.slot(slot)))
	}

	section("Sorts", Purple)
	for _, name := range c.Skills {
		spell := spells[name]
		fmt.Printf("   %s✦ %-16s%s %s(%d mana)%s %s\n", Purple, name, Reset, Blue, spell.Cost, Reset, Gray+spell.Description+Reset)
	}

	section("Progression", Gold)
	fmt.Printf("   Difficulté : %s   ·   Donjon : %d / 3 étages\n", difficultyOf(c.Difficulty).Name, c.FloorsCleared)
	fmt.Printf("   Monstres vaincus : %d   ·   Morts : %d   ·   Missions : %d / %d\n", c.Victories, c.Deaths, c.QuestsCompleted, len(quests))
	pause()
}

func equipmentLabel(c *Character, item string) string {
	switch {
	case item == "":
		return DarkGray + "— vide —" + Reset
	case isWeapon(item):
		return fmt.Sprintf("%s%s%s%s (+%d attaque)%s", Silver, Bold, item, Orange, weaponBonus(c, item), Reset)
	}
	return fmt.Sprintf("%s%s%s%s (+%d PV)%s", Silver, Bold, item, Green, gear[item].HP, Reset)
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
	typewrite(Gray, c.Name+" est tombé… mais les dieux de DunGo n'en ont pas fini avec lui.")
	c.HP = c.MaxHP / 2
	success("Vous ressuscitez avec %d / %d PV.", c.HP, c.MaxHP)
	return true
}

// gainXP ajoute de l'expérience et fait monter de niveau autant que nécessaire.
func gainXP(c *Character, amount int) {
	c.XP += amount
	fmt.Printf("  %s★ +%d XP%s  %s %d / %d\n", Purple+Bold, amount, Reset, bar(c.XP, c.XPMax, 20, Purple), c.XP, c.XPMax)

	class := classOf(c.Class)
	for c.XP >= c.XPMax {
		c.XP -= c.XPMax // l'excédent est conservé pour le niveau suivant
		c.XPMax = c.XPMax * 3 / 2
		c.Level++
		c.MaxHP += class.HPGain
		c.MaxMana += class.ManaGain
		c.Attack++
		c.Initiative++
		c.HP, c.Mana = c.MaxHP, c.MaxMana

		fmt.Println()
		printArtSlow(artLevelUp, goldColors...)
		section(fmt.Sprintf("Vous passez niveau %d !", c.Level), Gold)
		fmt.Printf("   %s♥ +%d PV max   %s♦ +%d mana max   %s» +1 attaque   %s» +1 initiative%s\n",
			Red, class.HPGain, Blue, class.ManaGain, Orange, Yellow, Reset)
		success("PV et mana entièrement restaurés !")
	}
}
