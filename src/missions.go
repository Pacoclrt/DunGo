package main

import "fmt"

// Quest : un contrat de chasse, proposé une fois l'précédent terminé.
type Quest struct {
	Title       string
	Description string
	Target      string // nom exact du monstre à vaincre
	Goal        int
	Reward      int
	XP          int
}

var quests = []Quest{
	{"Nuisibles", "Les gobelins volent les réserves du camp. Chassez-en trois !", "Gobelin chapardeur", 3, 40, 30},
	{"La meute", "Des loups rôdent autour du camp la nuit. Éliminez-en trois.", "Loup des cavernes", 3, 50, 40},
	{"Os à ronger", "Les squelettes des cryptes effraient les marchands. Deux suffiront.", "Squelette maudit", 2, 70, 60},
	{"Gros bras", "Deux trolls bloquent la route des marchands. Débarrassez-nous-en.", "Troll des cavernes", 2, 100, 90},
	{"Le roi est mort", "Grukk, le Roi Gobelin, menace tout le camp. Abattez-le !", "Grukk, le Roi Gobelin", 1, 120, 100},
}

func missionBoard(c *Character) {
	clearScreen()
	printArt(artQuestBoard, Brown, Gold, Gold, Brown)
	banner("MISSIONS", Gold)

	if c.QuestIndex >= len(quests) {
		fmt.Println()
		success("Toutes les missions sont terminées. Le camp vous doit une fière chandelle !")
		pause()
		return
	}

	quest := quests[c.QuestIndex]
	section(fmt.Sprintf("Mission %d / %d : %s", c.QuestIndex+1, len(quests), quest.Title), Gold)
	paragraph(Italic, quest.Description)
	fmt.Printf("   Objectif : vaincre %s%d × %s%s\n", Bold, quest.Goal, quest.Target, Reset)
	fmt.Printf("   Récompense : %s%d Y-Coins%s et %s%d XP%s\n", Gold+Bold, quest.Reward, Reset, Purple+Bold, quest.XP, Reset)

	switch {
	case !c.QuestActive:
		fmt.Println()
		option(1, "Accepter la mission", Green)
		option(0, "Retour au camp", Gray)
		if readChoice(0, 1) == 1 {
			c.QuestActive, c.QuestProgress = true, 0
			success("Mission acceptée : %s", quest.Title)
			pause()
		}

	case c.QuestProgress >= quest.Goal:
		showProgress(c, quest, Green)
		fmt.Println()
		option(1, "Réclamer la récompense", Gold)
		option(0, "Retour au camp", Gray)
		if readChoice(0, 1) == 1 {
			completeQuest(c, quest)
			pause()
		}

	default:
		showProgress(c, quest, Yellow)
		info("Mission en cours. Revenez quand l'objectif sera atteint !")
		pause()
	}
}

func showProgress(c *Character, quest Quest, color string) {
	fmt.Printf("   Progression : %s %d / %d\n", bar(c.QuestProgress, quest.Goal, 20, color), c.QuestProgress, quest.Goal)
}

func completeQuest(c *Character, quest Quest) {
	printArt(artCoins, goldColors...)
	c.Gold += quest.Reward
	c.GoldEarned += quest.Reward
	success("Mission « %s » terminée ! +%d Y-Coins", quest.Title, quest.Reward)
	gainXP(c, quest.XP)

	c.QuestIndex++
	c.QuestActive, c.QuestProgress = false, 0
	c.QuestsCompleted++
}

// updateQuest est appelée après chaque victoire pour faire avancer le contrat.
func updateQuest(c *Character, monsterName string) {
	if !c.QuestActive || c.QuestIndex >= len(quests) {
		return
	}
	quest := quests[c.QuestIndex]
	if monsterName != quest.Target || c.QuestProgress >= quest.Goal {
		return
	}
	c.QuestProgress++
	info("Mission « %s » : %d / %d", quest.Title, c.QuestProgress, quest.Goal)
	if c.QuestProgress >= quest.Goal {
		success("Objectif atteint ! Réclamez votre récompense au tableau des missions.")
	}
}
