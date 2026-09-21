package main

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

type Floor struct {
	Name     string
	Intro    string
	Color    string
	Colors   []string
	Rooms    int      // la dernière salle est celle du boss
	Monsters []string // clés du bestiaire
	Boss     string
}

var floors = []Floor{
	{
		Name:     "Étage 1 · Les Galeries Gobelines",
		Intro:    "Des tunnels qui sentent la chaussette. Un panneau gobelin : « HUMAINS DEHORS. SAUF SI HUMAINS DONNER Y-COINS. »",
		Color:    Green,
		Colors:   mossColors,
		Rooms:    5,
		Monsters: []string{"goblin", "wolf", "raven", "boar"},
		Boss:     "goblin_king",
	},
	{
		Name:     "Étage 2 · Les Cryptes Englouties",
		Intro:    "Des cryptes inondées et des squelettes qui claquent des dents. Quelque part, une liche compte son mana.",
		Color:    Sky,
		Colors:   iceColors,
		Rooms:    6,
		Monsters: []string{"skeleton", "troll", "boar", "raven"},
		Boss:     "lich",
	},
	{
		Name:     "Étage 3 · L'Antre d'Ignarok",
		Intro:    "Il fait chaud. Très chaud. Au fond, quelque chose ronfle : c'est Ignarok.",
		Color:    Orange,
		Colors:   fireColors,
		Rooms:    5,
		Monsters: []string{"krokmou", "troll", "skeleton"},
		Boss:     "dragon",
	},
}

// exploreDungeon affiche les étages. FloorsCleared dit combien d'étages
// le héros a déjà terminés : les suivants sont verrouillés.
func exploreDungeon(c *Character) {
	clearScreen()
	printArt(artDungeonGate, stoneColors...)
	banner("LE DONJON", Gold)
	showStatus(c)
	section("Où aller ?")

	for i, floor := range floors {
		status := DarkGray + "verrouillé" + Reset
		if i < c.FloorsCleared {
			status = Green + "√ terminé" + Reset
		}
		if i == c.FloorsCleared {
			status = Gold + "► boss : " + bestiary[floor.Boss].Name + Reset
		}
		option(i+1, padRight(floor.Name, 34)+" "+status)
	}
	back("Retour au camp")

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

// exploreFloor enchaîne les salles d'un étage. Si le héros fuit, meurt ou
// remonte au camp, il devra recommencer l'étage depuis la première salle.
func exploreFloor(c *Character, index int) {
	floor := floors[index]
	clearScreen()
	banner(floor.Name, floor.Color)
	fmt.Println()
	paragraph(floor.Color+Italic, floor.Intro)
	pause()

	for room := 1; room <= floor.Rooms; room++ {
		isBossRoom := room == floor.Rooms

		// 15 % des salles ordinaires cachent le sphinx à la place d'un monstre.
		sphinxSolved := false
		if !isBossRoom && rand.IntN(100) < 15 {
			sphinxSolved = riddleEvent(c, floor, index)
		}

		if !sphinxSolved {
			won := roomFight(c, floor, room)
			if !won {
				return
			}
		}

		if !isBossRoom {
			goOn := nextRoom(c, floor, room)
			if !goOn {
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

// roomFight fait apparaître le monstre de la salle. Renvoie true si le héros gagne.
func roomFight(c *Character, floor Floor, room int) bool {
	isBossRoom := room == floor.Rooms
	kind := floor.Monsters[rand.IntN(len(floor.Monsters))]
	if isBossRoom {
		kind = floor.Boss
	}
	m := newMonster(kind, c.Difficulty)
	m.Color = floor.Color
	m.Colors = floor.Colors

	clearScreen()
	if isBossRoom {
		banner("☠ BOSS : "+strings.ToUpper(m.Name)+" ☠", Red)
	} else {
		banner(fmt.Sprintf("%s · Salle %d / %d", floor.Name, room, floor.Rooms), floor.Color)
	}
	fmt.Println()
	printArt(m.Art, m.Colors...)
	fmt.Println()
	fmt.Printf("  %s%s surgit !%s  %s♥ %d PV · » %d attaque%s\n", Bold+m.Color, m.Name, Reset, Gray, m.MaxHP, m.Attack, Reset)
	say(m.Name, m.Color, m.Cry)
	fmt.Println()
	showStatus(c)
	option(1, "Combattre")
	option(2, "Fuir vers le camp")

	if readChoice(1, 2) == 2 {
		return false
	}
	result := fight(c, m, false)
	return result == Victory
}

// nextRoom demande au joueur s'il continue. Renvoie true pour continuer.
func nextRoom(c *Character, floor Floor, room int) bool {
	clearScreen()
	banner(floor.Name, floor.Color)
	fmt.Println()
	success("Salle %d / %d terminée !", room, floor.Rooms)
	showStatus(c)
	option(1, "Salle suivante")
	optionHint(2, "Retour au camp", "l'étage recommencera du début")
	return readChoice(1, 2) == 1
}

func victoryScreen(c *Character) {
	clearScreen()
	printArt(artVictory, campColors...)
	fmt.Println()
	printArt(artTrophy, campColors...)
	fmt.Println()
	paragraph(Gold, "Ignarok s'effondre dans un dernier ronflement. Le village peut enfin dormir, et votre nom est gravé sur la grande place : "+c.Name+", Terreur des Dragons.")

	section("Vos exploits")
	fmt.Printf("   Niveau atteint     %d\n", c.Level)
	fmt.Printf("   Monstres vaincus   %d\n", c.Victories)
	fmt.Printf("   Morts              %d\n", c.Deaths)
	pause()
}

type Riddle struct {
	Question string
	Answers  []string
	Correct  int // numéro de la bonne réponse : 1, 2 ou 3
}

var riddles = []Riddle{
	{"Combien de pattes a une araignée ?", []string{"Six", "Huit", "Dix"}, 2},
	{"Je tombe du ciel en hiver et je suis toute blanche. Qui suis-je ?", []string{"La neige", "La pluie", "Le sable"}, 1},
	{"Quel animal crache du feu ?", []string{"Le loup", "Le corbeau", "Le dragon"}, 3},
	{"Combien de jours y a-t-il dans une semaine ?", []string{"Cinq", "Sept", "Dix"}, 2},
	{"Qu'est-ce qui a des dents mais ne mord jamais ?", []string{"Un peigne", "Un loup", "Un gobelin"}, 1},
	{"Quelle potion faut-il boire pour regagner des PV ?", []string{"La Potion de poison", "La Potion de mana", "La Potion de vie"}, 3},
}

// riddleEvent pose une énigme. Renvoie true si la réponse est bonne.
func riddleEvent(c *Character, floor Floor, index int) bool {
	riddle := riddles[rand.IntN(len(riddles))]
	clearScreen()
	printArt(artSphinx, floor.Colors...)
	banner("LE SPHINX", floor.Color)
	say("Le Sphinx", floor.Color, "Réponds à mon énigme, voyageur. Une erreur… et tu le regretteras.")
	fmt.Println()
	paragraph(Bold, riddle.Question)
	fmt.Println()
	for i, answer := range riddle.Answers {
		option(i+1, answer)
	}

	answer := readChoice(1, len(riddle.Answers))
	fmt.Println()
	if answer != riddle.Correct {
		fail("Raté ! La bonne réponse était : %s", riddle.Answers[riddle.Correct-1])
		say("Le Sphinx", floor.Color, "Hé hé hé… Un ami monstre veut te dire bonjour.")
		pause()
		return false
	}

	reward := 15 * (index + 1)
	c.Gold += reward
	printArt(artCoins, campColors...)
	success("Bonne réponse ! Le sphinx vous laisse passer avec %d Y-Coins.", reward)
	gainXP(c, 10*(index+1))
	pause()
	return true
}
