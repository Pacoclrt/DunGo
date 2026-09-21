package main

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

// Floor décrit un étage : son ambiance, ses monstres et son boss. Chaque étage
// a sa couleur, que prennent son titre et tous les monstres qu'on y croise.
type Floor struct {
	Name     string
	Intro    string
	Color    string
	Colors   []string
	Rooms    int // la dernière salle contient toujours le boss
	Monsters []string
	Boss     string
}

var floors = []Floor{
	{
		Name:     "Étage 1 · Les Galeries Gobelines",
		Intro:    "Des tunnels qui sentent la chaussette. Un panneau gobelin : « HUMAINS DEHORS. SAUF SI HUMAINS DONNER Y-COINS. » Tout au fond : Grukk, le Roi Gobelin.",
		Color:    Green,
		Colors:   mossColors,
		Rooms:    5,
		Monsters: []string{"goblin", "wolf", "raven", "boar"},
		Boss:     "goblin_king",
	},
	{
		Name:     "Étage 2 · Les Cryptes Englouties",
		Intro:    "Des cryptes inondées et des squelettes qui claquent des dents (de froid ?). Tout au fond : Mor'Vath, une liche qui adore voler le mana des autres.",
		Color:    Sky,
		Colors:   iceColors,
		Rooms:    6,
		Monsters: []string{"skeleton", "troll", "boar", "raven"},
		Boss:     "lich",
	},
	{
		Name:     "Étage 3 · L'Antre d'Ignarok",
		Intro:    "Il fait chaud. Très chaud. Au fond, quelque chose ronfle : c'est Ignarok. Réveillez-le poliment… à coups d'épée.",
		Color:    Orange,
		Colors:   fireColors,
		Rooms:    5,
		Monsters: []string{"krokmou", "troll", "skeleton"},
		Boss:     "dragon",
	},
}

func exploreDungeon(c *Character) {
	clearScreen()
	printArt(artDungeonGate, stoneColors...)
	banner("LE DONJON", Gold)
	showStatus(c)
	section("Où aller ?")

	for i, floor := range floors {
		status, color := DarkGray+"■ verrouillé", DarkGray
		switch {
		case i < c.FloorsCleared:
			status, color = Green+"√ terminé", floor.Color
		case i == c.FloorsCleared:
			status, color = Gold+"► boss : "+bestiary[floor.Boss].Name, floor.Color
		}
		fmt.Printf("   %s[%d]%s %s%-34s%s %s%s\n", Gold+Bold, i+1, Reset, color+Bold, floor.Name, Reset, status, Reset)
	}
	back("Retour au camp")

	choice := readChoice(0, len(floors))
	if choice == 0 {
		return
	}
	if choice-1 > c.FloorsCleared {
		fail("Terminez d'abord l'étage %d ! Le dragon n'est pas si pressé.", c.FloorsCleared+1)
		pause()
		return
	}
	exploreFloor(c, choice-1)
}

// roomTrail dessine la progression dans l'étage : [√]──[►]──[ ]──[☠]
func roomTrail(room, total int) string {
	trail := "  "
	for i := 1; i <= total; i++ {
		if i > 1 {
			trail += DarkGray + "──" + Reset
		}
		switch {
		case i < room:
			trail += Green + "[√]" + Reset
		case i == room:
			trail += Gold + Bold + "[►]" + Reset
		case i == total:
			trail += Red + Bold + "[☠]" + Reset
		default:
			trail += DarkGray + "[ ]" + Reset
		}
	}
	return trail
}

// exploreFloor enchaîne les salles d'un étage. Le joueur peut remonter au
// camp entre deux salles, mais il perdra sa progression dans l'étage.
func exploreFloor(c *Character, index int) {
	floor := floors[index]

	clearScreen()
	dots("Vous descendez les marches, en évitant la troisième qui grince", Gray)
	fmt.Println()
	banner(floor.Name, floor.Color)
	fmt.Println()
	typewrite(floor.Color+Italic, floor.Intro)
	pause()

	for room := 1; room <= floor.Rooms; room++ {
		isBossRoom := room == floor.Rooms

		// Une salle ordinaire sur trois cache un événement plutôt qu'un monstre.
		needFight := true
		if !isBossRoom && rand.IntN(100) < 30 {
			needFight = roomEvent(c, index)
		}

		if needFight {
			kind := floor.Monsters[rand.IntN(len(floor.Monsters))]
			if isBossRoom {
				kind = floor.Boss
			}
			m := newMonster(kind, c.Difficulty)
			m.Color, m.Colors = floor.Color, floor.Colors
			if isBossRoom {
				bossIntro(m)
			}

			clearScreen()
			banner(fmt.Sprintf("%s · Salle %d / %d", floor.Name, room, floor.Rooms), floor.Color)
			fmt.Println(roomTrail(room, floor.Rooms))
			fmt.Println()
			printArt(m.Art, m.Color)
			fmt.Printf("\n  %s%s surgit !%s  %s♥ %d PV · » %d attaque%s\n", Bold+m.Color, m.Name, Reset, Gray, m.MaxHP, m.Attack, Reset)
			if !isBossRoom {
				fmt.Println("  " + Gray + Italic + "« " + m.Cry + " »" + Reset)
			}
			fmt.Println()
			showStatus(c)
			option(1, "Combattre")
			option(2, "Fuir vers le camp")
			if readChoice(1, 2) == 2 || fight(c, m, false) != Victory {
				return
			}
		}

		if !isBossRoom {
			clearScreen()
			banner(floor.Name, floor.Color)
			fmt.Println(roomTrail(room+1, floor.Rooms))
			fmt.Println()
			success("Salle %d / %d terminée !", room, floor.Rooms)
			showStatus(c)
			option(1, "Salle suivante")
			optionHint(2, "Retour au camp", "l'étage recommencera du début")
			if readChoice(1, 2) == 2 {
				return
			}
		}
	}

	c.FloorsCleared = max(c.FloorsCleared, index+1)
	if floor.Boss == "dragon" {
		victoryScreen(c)
		return
	}
	clearScreen()
	printArt(artTrophy, campColors...)
	banner(fmt.Sprintf("ÉTAGE %d TERMINÉ !", index+1), Gold)
	success("Boss vaincu ! L'étage %d est débloqué.", index+2)
	pause()
}

func bossIntro(m *Monster) {
	clearScreen()
	warn("Le sol tremble… Quelque chose de très gros approche.")
	wait(1500)
	clearScreen()
	banner("☠ BOSS : "+strings.ToUpper(m.Name)+" ☠", Red)
	fmt.Println()
	printArtSlow(m.Art, m.Colors...)
	fmt.Println()
	say(m.Name, m.Color, m.Cry)
	pause()
}

func victoryScreen(c *Character) {
	c.DragonSlain = true
	clearScreen()
	printArtSlow(artVictory, campColors...)
	fmt.Println()
	printArt(artTrophy, campColors...)
	typewrite(Gold, fmt.Sprintf("Ignarok s'effondre dans un dernier ronflement. Au village, tout le monde peut enfin dormir. On grave votre nom sur la grande place : %s, Terreur des Dragons. On vous offre aussi un mouton. Il n'a rien demandé.", c.Name))

	section("Vos exploits")
	for _, line := range [][2]string{
		{"Difficulté", difficultyOf(c.Difficulty).Name},
		{"Niveau atteint", fmt.Sprint(c.Level)},
		{"Monstres vaincus", fmt.Sprint(c.Victories)},
		{"Y-Coins amassés", fmt.Sprint(c.GoldEarned)},
		{"Morts", fmt.Sprint(c.Deaths)},
	} {
		fmt.Printf("   %-20s %s%s%s\n", line[0], Gold+Bold, line[1], Reset)
	}
	pause()
}
