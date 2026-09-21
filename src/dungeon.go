package main

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

// Floor décrit un étage : son ambiance, ses monstres et son boss.
type Floor struct {
	Name     string
	Intro    string
	Colors   []string
	Rooms    int // la dernière salle contient toujours le boss
	Monsters []string
	Boss     string
}

var floors = []Floor{
	{
		Name:     "Étage 1 · Les Galeries Gobelines",
		Intro:    "Des tunnels creusés à coups de griffes. Sur un mur, en lettres maladroites : « GOBELINS ICI CHEZ NOUS. HUMAINS DEHORS. SAUF SI HUMAINS DONNER Y-COINS. »",
		Colors:   poisonColors,
		Rooms:    5,
		Monsters: []string{"goblin", "wolf", "raven", "boar"},
		Boss:     "goblin_king",
	},
	{
		Name:     "Étage 2 · Les Cryptes Englouties",
		Intro:    "La chaleur du dragon a fait fondre les glaciers. Dans l'eau noire, des os s'entrechoquent… et une voix glaciale murmure votre nom.",
		Colors:   iceColors,
		Rooms:    6,
		Monsters: []string{"skeleton", "troll", "boar", "raven"},
		Boss:     "lich",
	},
	{
		Name:     "Étage 3 · L'Antre d'Ignarok",
		Intro:    "L'air brûle. Le sol tremble. Une inscription naine à moitié fondue : « Quand le Fléau inspire, que le brave se protège… ou boive. »",
		Colors:   fireColors,
		Rooms:    5,
		Monsters: []string{"krokmou", "troll", "skeleton"},
		Boss:     "dragon",
	},
}

func exploreDungeon(c *Character) {
	clearScreen()
	printArt(artDungeonGate, stoneColors...)
	banner("L'ENTRÉE DU DONJON DE KARAK-DÛM", Silver)
	showStatus(c)
	fmt.Println()

	for i, floor := range floors {
		status, color := DarkGray+"■ verrouillé", DarkGray
		switch {
		case i < c.FloorsCleared:
			status, color = Green+"√ terminé", Green
		case i == c.FloorsCleared:
			status, color = Yellow+"► à explorer", floor.Colors[1]
		}
		fmt.Printf("   %s[%d]%s %s%-36s%s %s%s\n", Gold+Bold, i+1, Reset, color+Bold, floor.Name, Reset, status, Reset)
	}
	fmt.Println()
	option(0, "Retour au camp", Gray)

	choice := readChoice(0, len(floors))
	if choice == 0 {
		return
	}
	if choice-1 > c.FloorsCleared {
		fail("Terminez d'abord l'étage %d !", c.FloorsCleared+1)
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
			trail += Yellow + Bold + "[►]" + Reset
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
	dots("Vous descendez les marches glissantes", Gray)
	printArtSlow(artStairs, floor.Colors...)
	banner(floor.Name, Gold)
	typewrite(floor.Colors[1]+Italic, floor.Intro)
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
			if isBossRoom {
				bossIntro(m)
			}

			clearScreen()
			banner(fmt.Sprintf("%s · Salle %d / %d", floor.Name, room, floor.Rooms), Gold)
			fmt.Println(roomTrail(room, floor.Rooms))
			fmt.Println()
			printArt(m.Art, m.Color)
			fmt.Printf("\n  %s%s surgit !%s\n", Bold+m.Color, m.Name, Reset)
			fmt.Printf("  %s♥ %d PV   » %d attaque   » %d initiative%s\n\n", Gray, m.MaxHP, m.Attack, m.Initiative, Reset)
			showStatus(c)
			option(1, "Combattre", Red)
			option(2, "Fuir vers le camp", Gray)
			if readChoice(1, 2) == 2 || fight(c, m, false) != Victory {
				return
			}
		}

		if !isBossRoom {
			clearScreen()
			banner(floor.Name, Gold)
			fmt.Println(roomTrail(room+1, floor.Rooms))
			fmt.Println()
			success("Salle %d / %d terminée !", room, floor.Rooms)
			showStatus(c)
			option(1, "Continuer vers la salle suivante", Yellow)
			option(2, "Remonter au camp", Gray)
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
	printArt(artTrophy, goldColors...)
	banner(floor.Name+" : TERMINÉ !", Green)
	success("Le boss est vaincu. L'étage suivant est débloqué !")
	pause()
}

func bossIntro(m *Monster) {
	clearScreen()
	warn("Le sol tremble…")
	wait(800)
	warn("Une présence terrifiante approche…")
	wait(800)
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
	printArtSlow(artVictory, goldColors...)
	fmt.Println()
	printArt(artTrophy, goldColors...)
	typewrite(Gold, fmt.Sprintf("Ignarok s'effondre dans un fracas qui fait trembler la montagne. Au camp, on grave enfin un nom sur la treizième stèle : %s, Tueur de Dragon.", c.Name))

	section("Vos exploits", Gold)
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
