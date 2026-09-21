// DunGo — un RPG en ligne de commande (Projet RED · Ymmersion).
package main

import "fmt"

func main() {
	for {
		clearScreen()
		fmt.Println()
		printArt(artLogo, fireColors...)
		fmt.Println()
		option(1, "Nouvelle partie")
		option(2, "Continuer")
		option(3, "Crédits")
		back("Quitter")

		switch readChoice(0, 3) {
		case 0:
			return
		case 1:
			newGame()
		case 2:
			continueGame()
		case 3:
			credits()
		}
	}
}

func newGame() {
	slot := chooseSlot("NOUVELLE PARTIE", false)
	if slot == 0 {
		return
	}
	c := characterCreation()
	c.SaveSlot = slot
	intro(c)
	campMenu(c)
}

func continueGame() {
	slot := chooseSlot("CONTINUER", true)
	if slot == 0 {
		return
	}
	c, err := loadGame(slot)
	if err != nil {
		fail("Cette sauvegarde est illisible.")
		pause()
		return
	}
	campMenu(c)
}

func intro(c *Character) {
	clearScreen()
	printArt(artDragon, fireColors...)
	banner("PROLOGUE", Gold)
	fmt.Println()
	paragraph(Silver, "Un dragon, Ignarok, s'est installé au fond du donjon. Il brûle les champs, mange les moutons et ronfle si fort que plus personne ne dort.")
	paragraph(Gold+Bold, "Le village cherche un héros. Un seul volontaire s'est présenté : vous, "+c.Name+".")
	fmt.Println()
	paragraph(Silver, "Traversez les 3 étages du donjon et battez leurs boss : Grukk, Mor'Vath, puis le dragon.")
	pause()
}

// campMenu est la boucle principale d'une partie : on y revient entre deux
// expéditions dans le donjon.
func campMenu(c *Character) {
	for {
		clearScreen()
		banner("LE CAMP", Gold)
		showStatus(c)
		section("Que faire ?")
		optionHint(1, "Explorer le donjon", "affronter les monstres et les boss")
		optionHint(2, "Marchand", "acheter et vendre")
		optionHint(3, "Forgeron", "fabriquer armes et armures")
		optionHint(4, "Missions", "contrats de chasse")
		optionHint(5, "Entraînement", "combat sans risque")
		optionHint(6, "Inventaire", "potions et équipement")
		optionHint(7, "Fiche du héros", "statistiques et sorts")
		option(8, "Sauvegarder")
		back("Quitter la partie")

		switch readChoice(0, 8) {
		case 0:
			if ask("Sauvegarder avant de partir ?") {
				saveGame(c)
			}
			return
		case 1:
			exploreDungeon(c)
		case 2:
			merchant(c)
		case 3:
			blacksmith(c)
		case 4:
			missionBoard(c)
		case 5:
			trainingFight(c)
		case 6:
			accessInventory(c)
		case 7:
			displayInfo(c)
		case 8:
			saveGame(c)
			pause()
		}
	}
}

func credits() {
	clearScreen()
	banner("CRÉDITS", Gold)
	fmt.Println()
	printArt(artPaco, fireColors...)
	fmt.Println()
	printArt(artSofiane, iceColors...)
	fmt.Println()
	printArt(artValentin, campColors...)
	fmt.Println()
	printArt(artAyman, mossColors...)
	pause()
}
