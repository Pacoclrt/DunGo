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
	switch rand.IntN(3) {
	case 0:
		fountainEvent(c)
	case 1:
		peddlerEvent(c)
	default:
		return !riddleEvent(c, floorIndex)
	}
	return false
}

func fountainEvent(c *Character) {
	printArt(artFountain, iceColors...)
	banner("UNE FONTAINE MAGIQUE", Cyan)
	paragraph(Sky+Italic, "Une eau claire et lumineuse coule au milieu des ténèbres…")
	c.HP = min(c.HP+c.MaxHP/2, c.MaxHP)
	c.Mana = min(c.Mana+c.MaxMana/2, c.MaxMana)
	success("Vous buvez à la fontaine : PV et mana en partie restaurés !")
	showStatus(c)
	pause()
}

// peddlerEvent : les mêmes potions qu'au camp, mais au double du prix.
func peddlerEvent(c *Character) {
	offers := []string{ItemHealthPotion, ItemManaPotion, ItemPoisonPotion}
	greeting := "Pssst… Des potions, l'ami ? Un peu chères, mais ici, pas de concurrence !"

	for {
		clearScreen()
		printArt(artPeddler, goldColors...)
		banner("UN MARCHAND AMBULANT", Gold)
		say("Filou le colporteur", Gold, greeting)
		greeting = "Autre chose avant que je file ?"
		showStatus(c)

		for i, item := range offers {
			itemLine(i+1, item, itemColor(item), fmt.Sprintf("%s%d Y-Coins%s", Gold, peddlerPrice(item), Reset))
		}
		fmt.Println()
		option(0, "Continuer l'exploration", Gray)

		choice := readChoice(0, len(offers))
		if choice == 0 {
			return
		}
		item, price := offers[choice-1], peddlerPrice(offers[choice-1])
		fmt.Println()
		if c.Gold < price {
			fail("Il vous manque %d Y-Coins.", price-c.Gold)
		} else if addInventory(c, item, 1) {
			c.Gold -= price
			success("Vous achetez : %s (-%d Y-Coins)", item, price)
		}
		pause()
	}
}

func peddlerPrice(item string) int {
	return items[item].Price * 2
}

// riddleEvent renvoie true si le joueur a trouvé la bonne réponse.
func riddleEvent(c *Character, floorIndex int) bool {
	riddle := riddles[rand.IntN(len(riddles))]
	printArt(artSphinx, Yellow, Gold, Brown, Brown)
	banner("LE SPHINX DE PIERRE", Yellow)
	say("Le Sphinx", Yellow, "Réponds à mon énigme, voyageur. Une erreur… et tu le regretteras.")
	fmt.Println()
	paragraph(Bold, riddle.Question)
	fmt.Println()
	for i, answer := range riddle.Answers {
		option(i+1, answer, White)
	}

	if readChoice(1, len(riddle.Answers)) != riddle.Correct {
		fmt.Println()
		fail("Mauvaise réponse ! La bonne réponse était : %s", riddle.Answers[riddle.Correct-1])
		say("Le Sphinx", Yellow, "Hé hé hé… Un monstre va s'occuper de toi.")
		pause()
		return false
	}

	reward := 15 * (floorIndex + 1)
	fmt.Println()
	printArt(artCoins, goldColors...)
	success("Bonne réponse ! Le sphinx vous laisse passer avec %d Y-Coins.", reward)
	c.Gold += reward
	c.GoldEarned += reward
	gainXP(c, 10*(floorIndex+1))
	pause()
	return true
}
