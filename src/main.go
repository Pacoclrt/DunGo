// ════════════════════════════════════════════════════════════════════════
//   main.go · [ACCUEIL]
//   Le point de départ du programme : l'écran titre, la nouvelle partie,
//   « Continuer », le prologue et les crédits.
//
//   DunGo — un RPG en ligne de commande (Projet RED · Ymmersion).
//   Pour lancer le jeu depuis le dossier du projet : go run ./src
// ════════════════════════════════════════════════════════════════════════

package main

import "fmt"

func main() { // [ACCUEIL] sert à afficher l'écran titre en boucle (Nouvelle partie, Continuer, Crédits) jusqu'à ce que le joueur quitte : c'est ici que le programme commence
	for {
		clearScreen()
		fmt.Println()
		printArt(artLogo, fireColors...)
		fmt.Println()
		option(1, "Nouvelle partie")
		option(2, "Continuer")
		option(3, "Crédits")
		backOption("Quitter")

		switch readChoice(0, 3) {
		case 0:
			return // sortir de main = fermer le programme
		case 1:
			newGame()
		case 2:
			continueGame()
		case 3:
			credits()
		}
	}
}

func newGame() { // [ACCUEIL] sert à lancer une nouvelle partie : choix de l'emplacement, création du héros, prologue, puis le camp
	slot := chooseSlot("NOUVELLE PARTIE", false)
	if slot == 0 {
		return
	}
	character := createCharacter()
	character.SaveSlot = slot
	prologue(character)
	campMenu(character)
}

func continueGame() { // [ACCUEIL] sert à reprendre une partie sauvegardée : choix de l'emplacement, chargement, puis le camp
	slot := chooseSlot("CONTINUER", true)
	if slot == 0 {
		return
	}
	character, err := loadGame(slot)
	if err != nil {
		fail("Cette sauvegarde est illisible.")
		pause()
		return
	}
	campMenu(character)
}

func prologue(character *Character) { // [ACCUEIL] sert à raconter le début de l'histoire (le dragon Ignarok) à une nouvelle partie
	clearScreen()
	printArt(artDragon, fireColors...)
	banner("PROLOGUE", Gold)
	fmt.Println()
	paragraph(Silver, "Un dragon, Ignarok, s'est installé au fond du donjon. Il brûle les champs, mange les moutons et ronfle si fort que plus personne ne dort.")
	paragraph(Gold+Bold, "Le village cherche un héros. Un seul volontaire s'est présenté : vous, "+character.Name+".")
	fmt.Println()
	paragraph(Silver, "Traversez les 3 étages du donjon et battez leurs boss : Grukk, Mor'Vath, puis le dragon.")
	pause()
}

func credits() { // [ACCUEIL] sert à afficher les prénoms de l'équipe, chacun avec son dégradé
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
