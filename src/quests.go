// ════════════════════════════════════════════════════════════════════════
//   quests.go · [MISSIONS]
//   Le tableau des missions du camp : des contrats de chasse (« vaincre
//   3 gobelins ») à faire un par un, dans l'ordre, contre des Y-Coins et
//   de l'XP.
//
//   Une mission passe par 3 états, rangés dans le héros :
//     à accepter   QuestActive = false
//     en cours     QuestActive = true, QuestProgress < Goal
//     réussie      QuestActive = true, QuestProgress = Goal → récompense,
//                  puis QuestIndex + 1 (on passe à la mission suivante)
// ════════════════════════════════════════════════════════════════════════

package main

import "fmt"

type Quest struct { // [MISSIONS] sert à décrire une mission : quel monstre vaincre, combien de fois, et la récompense
	Title        string
	Description  string
	Target       string // le NOM EXACT du monstre à vaincre (comparé à Monster.Name)
	Goal         int    // combien il faut en vaincre
	RewardYCoins int
	RewardXP     int
}

var quests = []Quest{ // [MISSIONS] sert à lister les 5 missions, dans l'ordre où elles sont proposées
	{Title: "Nuisibles", Description: "Les gobelins volent les réserves du camp. Chassez-en trois !",
		Target: "Gobelin chapardeur", Goal: 3, RewardYCoins: 40, RewardXP: 30},
	{Title: "La meute", Description: "Des loups rôdent autour du camp la nuit. Éliminez-en trois.",
		Target: "Loup des cavernes", Goal: 3, RewardYCoins: 50, RewardXP: 40},
	{Title: "Os à ronger", Description: "Les squelettes des cryptes effraient les marchands. Deux suffiront.",
		Target: "Squelette maudit", Goal: 2, RewardYCoins: 70, RewardXP: 60},
	{Title: "Gros bras", Description: "Deux trolls bloquent la route des marchands. Débarrassez-nous-en.",
		Target: "Troll des cavernes", Goal: 2, RewardYCoins: 100, RewardXP: 90},
	{Title: "Le roi est mort", Description: "Grukk, le Roi Gobelin, menace tout le camp. Abattez-le !",
		Target: "Grukk, le Roi Gobelin", Goal: 1, RewardYCoins: 120, RewardXP: 100},
}

func questMenu(character *Character) { // [MISSIONS] sert à afficher la mission en cours et, selon son état, à l'accepter ou à réclamer la récompense
	clearScreen()
	printArt(artQuestBoard, campColors...)
	banner("MISSIONS", Gold)

	// Toutes les missions sont faites : il n'y a plus rien à proposer.
	if character.QuestIndex >= len(quests) {
		fmt.Println()
		success("Toutes les missions sont terminées. Le camp vous doit une fière chandelle !")
		pause()
		return
	}

	// La fiche de la mission en cours
	quest := quests[character.QuestIndex]
	section(fmt.Sprintf("Mission %d / %d : %s", character.QuestIndex+1, len(quests), quest.Title))
	paragraph(Italic, quest.Description)
	fmt.Printf("   Objectif : vaincre %d × %s\n", quest.Goal, quest.Target)
	fmt.Printf("   Récompense : %s%d Y-Coins%s et %s%d XP%s\n", Gold, quest.RewardYCoins, Reset, Purple, quest.RewardXP, Reset)
	if character.QuestActive {
		fmt.Printf("   Progression : %d / %d\n", character.QuestProgress, quest.Goal)
	}
	fmt.Println()

	// État 1 : pas encore acceptée → on propose de l'accepter
	if !character.QuestActive {
		option(1, "Accepter la mission")
		backOption("Retour au camp")
		if readChoice(0, 1) == 1 {
			character.QuestActive = true
			character.QuestProgress = 0
			success("Mission acceptée !")
			pause()
		}
		return
	}

	// État 2 : en cours → rien à faire ici
	if character.QuestProgress < quest.Goal {
		info("Mission en cours. Revenez quand l'objectif sera atteint.")
		pause()
		return
	}

	// État 3 : objectif atteint → on réclame la récompense
	option(1, "Réclamer la récompense")
	backOption("Retour au camp")
	if readChoice(0, 1) == 1 {
		printArt(artCoins, campColors...)
		character.YCoins += quest.RewardYCoins
		success("Mission « %s » terminée ! +%d Y-Coins", quest.Title, quest.RewardYCoins)
		gainXP(character, quest.RewardXP)
		character.QuestIndex++ // on passe à la mission suivante…
		character.QuestActive = false
		character.QuestProgress = 0 // …qui n'est pas encore acceptée
		pause()
	}
}

func updateQuest(character *Character, monsterName string) { // [MISSIONS] sert à faire avancer la mission en cours si le monstre vaincu est la cible (appelée après chaque victoire)
	if !character.QuestActive || character.QuestIndex >= len(quests) {
		return
	}
	quest := quests[character.QuestIndex]
	if monsterName != quest.Target || character.QuestProgress >= quest.Goal {
		return
	}
	character.QuestProgress++
	info("Mission « %s » : %d / %d", quest.Title, character.QuestProgress, quest.Goal)
}
