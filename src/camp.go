// ════════════════════════════════════════════════════════════════════════
//   camp.go · [CAMP]
//   Le camp est le cœur d'une partie : on y revient toujours, entre deux
//   expéditions dans le donjon. Son menu mène à tous les autres écrans.
//   Ici aussi : l'entraînement (un combat sans risque contre un gobelin).
// ════════════════════════════════════════════════════════════════════════

package main

func campMenu(character *Character) { // [CAMP] sert à afficher le menu du camp et à ouvrir l'écran choisi ; c'est la boucle principale d'une partie
	for {
		clearScreen()
		banner("LE CAMP", Gold)
		showCharacterBar(character)
		section("Que faire ?")
		optionHint(1, "Explorer le donjon", "affronter les monstres et les boss")
		optionHint(2, "Marchand", "acheter et vendre")
		optionHint(3, "Forgeron", "fabriquer armes et armures")
		optionHint(4, "Missions", "contrats de chasse")
		optionHint(5, "Entraînement", "combat sans risque")
		optionHint(6, "Inventaire", "potions et équipement")
		optionHint(7, "Fiche du héros", "statistiques et sorts")
		option(8, "Sauvegarder")
		backOption("Quitter la partie")

		switch readChoice(0, 8) {
		case 0:
			if askYesNo("Sauvegarder avant de partir ?") {
				saveGame(character)
				pause() // sinon le message « Partie sauvegardée » serait effacé aussitôt
			}
			return
		case 1:
			dungeonMenu(character)
		case 2:
			merchantMenu(character)
		case 3:
			blacksmithMenu(character)
		case 4:
			questMenu(character)
		case 5:
			trainingFight(character)
		case 6:
			inventoryMenu(character)
		case 7:
			showCharacterSheet(character)
		case 8:
			saveGame(character)
			pause()
		}
	}
}

func trainingFight(character *Character) { // [CAMP] sert à lancer un combat d'entraînement : pas de butin, pas d'objets, et le héros récupère ses PV et son mana à la fin
	clearScreen()
	printArt(artArena, campColors...)
	banner("L'ENTRAÎNEMENT", Gold)
	say("Sergent Grol", Gold, "Mon gobelin cogne mou, mais il cogne. Ici, pas de risque… et pas de butin !")
	pause()

	// Le gobelin d'entraînement est toujours en Normal ("" = difficulté inconnue → Normal).
	monster := newMonster("training_goblin", "")
	monster.Color = Gold
	monster.ArtColors = campColors

	// On retient les PV et le mana d'avant, pour les rendre après le combat.
	hpBefore := character.HP
	manaBefore := character.Mana
	fight(character, monster, true)
	character.HP = hpBefore
	character.Mana = manaBefore
}
