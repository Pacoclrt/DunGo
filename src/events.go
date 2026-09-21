package main

import (
	"fmt"
	"math/rand/v2"
)

// Riddle est une énigme posée par le sphinx : trois réponses, une seule bonne.
type Riddle struct {
	Question string
	Answers  []string
	Correct  int // 1, 2 ou 3
}

var riddles = []Riddle{
	{"Combien de pattes a une araignée ?", []string{"Six", "Huit", "Dix"}, 2},
	{"Je tombe du ciel en hiver et je suis toute blanche. Qui suis-je ?", []string{"La neige", "La pluie", "Le sable"}, 1},
	{"Quel animal crache du feu ?", []string{"Le loup", "Le corbeau", "Le dragon"}, 3},
	{"Combien de jours y a-t-il dans une semaine ?", []string{"Cinq", "Sept", "Dix"}, 2},
	{"Qu'est-ce qui a des dents mais ne mord jamais ?", []string{"Un peigne", "Un loup", "Un gobelin"}, 1},
	{"Quelle potion faut-il boire pour regagner des PV ?", []string{"La Potion de poison", "La Potion de mana", "La Potion de vie"}, 3},
}

// roomEvent tire un événement au hasard. Il renvoie true si un combat doit
// quand même avoir lieu (énigme ratée).
func roomEvent(c *Character, floorIndex int) bool {
	clearScreen()
	if rand.IntN(2) == 0 {
		fountainEvent(c, floorIndex)
		return false
	}
	return !riddleEvent(c, floorIndex)
}

func fountainEvent(c *Character, floorIndex int) {
	floor := floors[floorIndex]
	printArt(artFountain, floor.Colors...)
	banner("UNE FONTAINE MAGIQUE", floor.Color)
	paragraph(Silver+Italic, "Une eau claire et lumineuse coule au milieu des ténèbres. Elle a un petit goût de menthe.")
	c.HP = min(c.HP+c.MaxHP/2, c.MaxHP)
	c.Mana = min(c.Mana+c.MaxMana/2, c.MaxMana)
	success("Vous buvez : PV et mana en partie restaurés !")
	showStatus(c)
	pause()
}

// riddleEvent renvoie true si le joueur a trouvé la bonne réponse.
func riddleEvent(c *Character, floorIndex int) bool {
	riddle, floor := riddles[rand.IntN(len(riddles))], floors[floorIndex]
	printArt(artSphinx, floor.Colors...)
	banner("LE SPHINX", floor.Color)
	say("Le Sphinx", floor.Color, "Réponds à mon énigme, voyageur. Une erreur… et tu le regretteras.")
	fmt.Println()
	paragraph(Bold, riddle.Question)
	fmt.Println()
	for i, answer := range riddle.Answers {
		option(i+1, answer)
	}

	if readChoice(1, len(riddle.Answers)) != riddle.Correct {
		fmt.Println()
		fail("Raté ! La bonne réponse était : %s", riddle.Answers[riddle.Correct-1])
		say("Le Sphinx", floor.Color, "Hé hé hé… Un ami monstre veut te dire bonjour.")
		pause()
		return false
	}

	reward := 15 * (floorIndex + 1)
	fmt.Println()
	printArt(artCoins, campColors...)
	c.Gold += reward
	c.GoldEarned += reward
	success("Bonne réponse ! Le sphinx boude, mais vous laisse passer avec %d Y-Coins.", reward)
	gainXP(c, 10*(floorIndex+1))
	pause()
	return true
}
