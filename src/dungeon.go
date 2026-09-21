// ════════════════════════════════════════════════════════════════════════
//   dungeon.go · [DONJON]
//   Les 3 étages du donjon. Un étage = une suite de salles ; la dernière
//   salle est toujours celle du boss. Battre le boss débloque l'étage
//   suivant, et battre le dragon (étage 3) termine l'aventure.
// ════════════════════════════════════════════════════════════════════════

package main

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

// ════════════════════════════════════════════════════════════════════════
//   LES ÉTAGES
// ════════════════════════════════════════════════════════════════════════

type Floor struct { // [DONJON] sert à décrire un étage : son nom, ses couleurs, son nombre de salles, ses monstres et son boss
	Name      string
	Intro     string   // le texte d'ambiance affiché en entrant
	Color     string   // la couleur des titres et des noms de monstres de l'étage
	ArtColors []string // le dégradé des dessins de l'étage
	Rooms     int      // le nombre de salles ; la dernière est celle du boss
	Monsters  []string // les clés des monstres qui peuvent apparaître (bestiary)
	Boss      string   // la clé du boss (bestiary)
}

var floors = []Floor{ // [DONJON] sert à lister les 3 étages, dans l'ordre
	{
		Name:      "Étage 1 · Les Galeries Gobelines",
		Intro:     "Des tunnels qui sentent la chaussette. Un panneau gobelin : « HUMAINS DEHORS. SAUF SI HUMAINS DONNER Y-COINS. »",
		Color:     Green,
		ArtColors: mossColors,
		Rooms:     5,
		Monsters:  []string{"goblin", "wolf", "raven", "boar"},
		Boss:      "goblin_king",
	},
	{
		Name:      "Étage 2 · Les Cryptes Englouties",
		Intro:     "Des cryptes inondées et des squelettes qui claquent des dents. Quelque part, une liche compte son mana.",
		Color:     Sky,
		ArtColors: iceColors,
		Rooms:     6,
		Monsters:  []string{"skeleton", "troll", "boar", "raven"},
		Boss:      "lich",
	},
	{
		Name:      "Étage 3 · L'Antre d'Ignarok",
		Intro:     "Il fait chaud. Très chaud. Au fond, quelque chose ronfle : c'est Ignarok.",
		Color:     Orange,
		ArtColors: fireColors,
		Rooms:     5,
		Monsters:  []string{"krokmou", "troll", "skeleton"},
		Boss:      "dragon",
	},
}

// ════════════════════════════════════════════════════════════════════════
//   LE CHOIX DE L'ÉTAGE
// ════════════════════════════════════════════════════════════════════════

func dungeonMenu(character *Character) { // [DONJON] sert à afficher les 3 étages (terminé, en cours ou verrouillé) et à entrer dans celui choisi
	clearScreen()
	printArt(artDungeonGate, stoneColors...)
	banner("LE DONJON", Gold)
	showCharacterBar(character)
	section("Où aller ?")

	// FloorsCleared suffit à tout savoir : les étages avant lui sont
	// terminés, celui à sa place est le prochain, ceux d'après sont fermés.
	// Ex : FloorsCleared = 1 → étage 1 terminé, étage 2 ouvert, étage 3 verrouillé.
	for index, floor := range floors {
		status := DarkGray + "verrouillé" + Reset
		if index < character.FloorsCleared {
			status = Green + "√ terminé" + Reset
		}
		if index == character.FloorsCleared {
			status = Gold + "► boss : " + bestiary[floor.Boss].Name + Reset
		}
		option(index+1, padRight(floor.Name, 34)+" "+status)
	}
	backOption("Retour au camp")

	choice := readChoice(0, len(floors))
	if choice == 0 {
		return
	}
	floorIndex := choice - 1 // le joueur tape 1 à 3, la liste commence à 0
	if floorIndex > character.FloorsCleared {
		fail("Terminez d'abord l'étage %d !", character.FloorsCleared+1)
		pause()
		return
	}
	exploreFloor(character, floorIndex)
}

// ════════════════════════════════════════════════════════════════════════
//   L'EXPLORATION D'UN ÉTAGE
// ════════════════════════════════════════════════════════════════════════

func exploreFloor(character *Character, floorIndex int) { // [DONJON] sert à enchaîner les salles d'un étage jusqu'au boss ; fuir, mourir ou remonter au camp oblige à recommencer l'étage
	floor := floors[floorIndex]
	clearScreen()
	banner(floor.Name, floor.Color)
	fmt.Println()
	paragraph(floor.Color+Italic, floor.Intro)
	pause()

	for room := 1; room <= floor.Rooms; room++ {
		isBossRoom := room == floor.Rooms

		// 1. 15 % des salles ordinaires cachent le sphinx à la place du
		//    monstre. Bonne réponse = pas de combat dans cette salle.
		sphinxSolved := false
		if !isBossRoom && rand.IntN(100) < 15 {
			sphinxSolved = sphinxRiddle(character, floor, floorIndex)
		}

		// 2. Le combat de la salle (sauf si le sphinx a été vaincu)
		if !sphinxSolved {
			won := roomFight(character, floor, room)
			if !won {
				return // fuite ou mort : retour au camp, l'étage est à refaire
			}
		}

		// 3. Entre deux salles, le joueur peut remonter au camp
		if !isBossRoom {
			goOn := askNextRoom(character, floor, room)
			if !goOn {
				return
			}
		}
	}

	// 4. Le boss est vaincu. L'étage suivant se débloque, mais seulement la
	//    première fois (refaire l'étage 1 plus tard ne débloque rien de plus).
	firstTime := floorIndex == character.FloorsCleared
	if firstTime {
		character.FloorsCleared++
	}
	if floor.Boss == "dragon" {
		victoryScreen(character)
		return
	}
	clearScreen()
	printArt(artTrophy, campColors...)
	banner(fmt.Sprintf("ÉTAGE %d TERMINÉ !", floorIndex+1), Gold)
	if firstTime {
		success("Boss vaincu ! L'étage %d est débloqué.", floorIndex+2)
	} else {
		success("Boss vaincu, encore une fois !")
	}
	pause()
}

func roomFight(character *Character, floor Floor, room int) bool { // [DONJON] sert à faire apparaître le monstre de la salle (le boss dans la dernière) et à lancer le combat ; renvoie true si le héros gagne
	// 1. Quel monstre ? Un au hasard parmi ceux de l'étage, ou le boss.
	isBossRoom := room == floor.Rooms
	key := floor.Monsters[rand.IntN(len(floor.Monsters))]
	if isBossRoom {
		key = floor.Boss
	}
	monster := newMonster(key, character.Difficulty)
	monster.Color = floor.Color // le monstre prend les couleurs de l'étage
	monster.ArtColors = floor.ArtColors

	// 2. L'apparition
	clearScreen()
	if isBossRoom {
		banner("☠ BOSS : "+strings.ToUpper(monster.Name)+" ☠", Red)
	} else {
		banner(fmt.Sprintf("%s · Salle %d / %d", floor.Name, room, floor.Rooms), floor.Color)
	}
	fmt.Println()
	printArt(monster.Art, monster.ArtColors...)
	fmt.Println()
	fmt.Printf("  %s%s surgit !%s  %s♥ %d PV · » %d attaque%s\n", Bold+monster.Color, monster.Name, Reset, Gray, monster.MaxHP, monster.Attack, Reset)
	say(monster.Name, monster.Color, monster.Cry)
	fmt.Println()
	showCharacterBar(character)

	// 3. Combattre ou repartir (possible même devant un boss, avant le combat)
	option(1, "Combattre")
	option(2, "Fuir vers le camp")
	if readChoice(1, 2) == 2 {
		return false
	}
	return fight(character, monster, false) == Victory
}

func askNextRoom(character *Character, floor Floor, room int) bool { // [DONJON] sert à demander au joueur s'il continue vers la salle suivante ; renvoie true pour continuer, false pour remonter au camp
	clearScreen()
	banner(floor.Name, floor.Color)
	fmt.Println()
	success("Salle %d / %d terminée !", room, floor.Rooms)
	showCharacterBar(character)
	option(1, "Salle suivante")
	optionHint(2, "Retour au camp", "l'étage recommencera du début")
	return readChoice(1, 2) == 1
}

func victoryScreen(character *Character) { // [DONJON] sert à afficher l'écran de fin quand Ignarok est vaincu, avec les exploits du héros
	clearScreen()
	printArt(artVictory, campColors...)
	fmt.Println()
	printArt(artTrophy, campColors...)
	fmt.Println()
	paragraph(Gold, "Ignarok s'effondre dans un dernier ronflement. Le village peut enfin dormir, et votre nom est gravé sur la grande place : "+character.Name+", Terreur des Dragons.")

	section("Vos exploits")
	fmt.Printf("   Niveau atteint     %d\n", character.Level)
	fmt.Printf("   Monstres vaincus   %d\n", character.Victories)
	fmt.Printf("   Morts              %d\n", character.Deaths)
	pause()
}
