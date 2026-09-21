// ════════════════════════════════════════════════════════════════════════
//   sphinx.go · [SPHINX]
//   Le sphinx se cache dans 15 % des salles ordinaires (dungeon.go). Il
//   pose une énigme à 3 réponses :
//     • bonne réponse   → des Y-Coins, de l'XP, et pas de combat
//     • mauvaise réponse → le monstre de la salle attaque quand même
// ════════════════════════════════════════════════════════════════════════

package main

import (
	"fmt"
	"math/rand/v2"
)

type Riddle struct { // [SPHINX] sert à décrire une énigme : la question, les 3 réponses proposées et le numéro de la bonne
	Question      string
	Answers       []string
	CorrectAnswer int // le NUMÉRO de la bonne réponse : 1, 2 ou 3 (pas sa position dans la liste, qui commence à 0)
}

var riddles = []Riddle{ // [SPHINX] sert à lister les 6 énigmes ; le sphinx en tire une au hasard
	{"Combien de pattes a une araignée ?", []string{"Six", "Huit", "Dix"}, 2},
	{"Je tombe du ciel en hiver et je suis toute blanche. Qui suis-je ?", []string{"La neige", "La pluie", "Le sable"}, 1},
	{"Quel animal crache du feu ?", []string{"Le loup", "Le corbeau", "Le dragon"}, 3},
	{"Combien de jours y a-t-il dans une semaine ?", []string{"Cinq", "Sept", "Dix"}, 2},
	{"Qu'est-ce qui a des dents mais ne mord jamais ?", []string{"Un peigne", "Un loup", "Un gobelin"}, 1},
	{"Quelle potion faut-il boire pour regagner des PV ?", []string{"La Potion de poison", "La Potion de mana", "La Potion de vie"}, 3},
}

func sphinxRiddle(character *Character, floor Floor, floorIndex int) bool { // [SPHINX] sert à poser une énigme au hasard ; renvoie true si la réponse est bonne (récompense, pas de combat)
	riddle := riddles[rand.IntN(len(riddles))]
	clearScreen()
	printArt(artSphinx, floor.ArtColors...)
	banner("LE SPHINX", floor.Color)
	say("Le Sphinx", floor.Color, "Réponds à mon énigme, voyageur. Une erreur… et tu le regretteras.")
	fmt.Println()
	paragraph(Bold, riddle.Question)
	fmt.Println()
	for index, answer := range riddle.Answers {
		option(index+1, answer)
	}

	answer := readChoice(1, len(riddle.Answers))
	fmt.Println()

	// Mauvaise réponse : le monstre de la salle arrive (géré par exploreFloor)
	if answer != riddle.CorrectAnswer {
		fail("Raté ! La bonne réponse était : %s", riddle.Answers[riddle.CorrectAnswer-1])
		say("Le Sphinx", floor.Color, "Hé hé hé… Un ami monstre veut te dire bonjour.")
		pause()
		return false
	}

	// Bonne réponse : la récompense grandit avec l'étage (15, 30, 45 Y-Coins et 10, 20, 30 XP)
	floorNumber := floorIndex + 1
	reward := 15 * floorNumber
	character.YCoins += reward
	printArt(artCoins, campColors...)
	success("Bonne réponse ! Le sphinx vous laisse passer avec %d Y-Coins.", reward)
	gainXP(character, 10*floorNumber)
	pause()
	return true
}
