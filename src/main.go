// DunGo — L'Antre d'Ignarok, un RPG en ligne de commande (Projet RED · Ymmersion).
//
// Le jeu se lit comme une pile d'écrans : titre → camp → donjon → combat.
// Chaque écran est une boucle « efface, affiche, lis un choix ».
package main

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

var tips = []string{
	"Les loups et les corbeaux laissent des ressources pour la forge. Ils n'en ont plus besoin.",
	"Les fontaines du donjon rendent la moitié de vos PV et de votre mana. Gratuitement, en plus.",
	"La Potion de poison se jette sur l'ennemi. Pas dans votre bouche.",
	"Quand Ignarok inspire profondément, ce n'est pas pour chanter.",
	"Tous les 3 tours, les monstres frappent deux fois plus fort. Eux aussi savent compter.",
	"Une arme faite pour votre classe donne +2 attaque en bonus.",
	"Les missions paient bien. Mieux que les gobelins, en tout cas.",
	"Second Souffle soigne aussi les saignements et les brûlures.",
	"La Boule de Feu fait brûler l'ennemi pendant 3 tours. Ça sent le grillé.",
	"On ne fuit pas un boss. D'autres ont essayé.",
	"Mourir coûte 20 % de vos Y-Coins. Les dieux de DunGo ne font pas crédit.",
}

// L'écran titre : chaque entrée est un simple appel de fonction.
var titleActions = []struct {
	Label  string
	Action func()
}{
	{"Nouvelle partie", newGame},
	{"Continuer", continueGame},
	{"Crédits", credits},
}

// Le menu du camp. L'ordre de la table donne les numéros affichés : ce qu'on
// fait le plus souvent vient en premier.
var campActions = []struct {
	Label  string
	Hint   string
	Action func(*Character)
}{
	{"Explorer le donjon", "affronter les monstres et les boss", exploreDungeon},
	{"Marchand", "acheter et vendre", merchant},
	{"Forgeron", "fabriquer armes et armures", blacksmith},
	{"Missions", "contrats de chasse bien payés", missionBoard},
	{"Entraînement", "combat pour de faux, sans risque", trainingFight},
	{"Inventaire", "potions et équipement", accessInventory},
	{"Fiche du héros", "statistiques, sorts, équipement", displayInfo},
	{"Sauvegarder", "", saveAndPause},
}

func main() {
	for {
		titleScreen()
		choice := readChoice(0, len(titleActions))
		if choice == 0 {
			goodbye()
			return
		}
		titleActions[choice-1].Action()
	}
}

func titleScreen() {
	clearScreen()
	fmt.Println()
	printArt(artLogo, fireColors...)
	fmt.Println(Gold + Bold + "              ~ L'ANTRE D'IGNAROK ~" + Reset)
	fmt.Println(Gray + Italic + "      Un dragon. Trois étages. Un seul héros : vous." + Reset)
	fmt.Println()
	for i, entry := range titleActions {
		option(i+1, entry.Label)
	}
	back("Quitter")
	fmt.Println("\n" + DarkGray + "   Projet RED · Ymmersion · Paco · Sofiane · Valentin · Ayman" + Reset)
}

func newGame() {
	slot := chooseSlot("NOUVELLE PARTIE : CHOISISSEZ UN EMPLACEMENT", false)
	if slot == 0 {
		return
	}
	c := characterCreation()
	c.SaveSlot = slot
	intro(c)
	campMenu(c)
}

func continueGame() {
	slot := chooseSlot("CONTINUER UNE PARTIE", true)
	if slot == 0 {
		return
	}
	c, err := loadGame(slot)
	if err != nil {
		fail("Cette sauvegarde est illisible. Un gobelin a dû la mâchouiller.")
		pause()
		return
	}
	success("Bon retour parmi nous, %s !", c.Name)
	wait(1000)
	campMenu(c)
}

// intro raconte l'histoire en trois phrases, puis montre le chemin à suivre.
func intro(c *Character) {
	clearScreen()
	printArt(artDragon, fireColors...)
	banner("PROLOGUE", Gold)
	fmt.Println()
	typewrite(Silver, "Un dragon, Ignarok, s'est installé au fond du donjon, sous la montagne.")
	typewrite(Silver, "Depuis, il brûle les champs, mange les moutons et ronfle si fort que plus personne ne dort.")
	typewrite(Gold+Bold, fmt.Sprintf("Le village cherche un héros. Un seul volontaire s'est présenté : vous, %s.", c.Name))

	section("Votre mission")
	showQuestMap()
	fmt.Println()
	paragraph(Silver, "Préparez-vous au camp, puis descendez battre le boss de chaque étage. Le dernier, c'est le dragon.")
	pause()
}

// showQuestMap dessine le chemin du camp jusqu'au dragon, boss par boss.
func showQuestMap() {
	path, bosses := Gold+Bold+"   Camp"+Reset, "       "
	for i, floor := range floors {
		boss, _, _ := strings.Cut(bestiary[floor.Boss].Name, ",")
		path += DarkGray + "  ──►  " + Reset + floor.Color + Bold + padRight(fmt.Sprintf("Étage %d", i+1), 8) + Reset
		bosses += "       " + Gray + padRight(boss, 8) + Reset
	}
	fmt.Println(path)
	fmt.Println(bosses)
}

// campMenu est la boucle principale d'une partie : on y revient entre
// chaque expédition dans le donjon.
func campMenu(c *Character) {
	for {
		clearScreen()
		banner("LE CAMP", Gold)
		showStatus(c)
		paragraph(Gray+Italic, "Astuce : "+tips[rand.IntN(len(tips))])
		section("Que faire ?")
		for i, entry := range campActions {
			optionHint(i+1, entry.Label, entry.Hint)
		}
		back("Quitter la partie")

		choice := readChoice(0, len(campActions))
		if choice == 0 {
			if ask("Sauvegarder avant de partir ?") {
				saveGame(c)
				wait(800)
			}
			return
		}
		campActions[choice-1].Action(c)
	}
}

func saveAndPause(c *Character) {
	fmt.Println()
	saveGame(c)
	pause()
}

func credits() {
	authors := []struct {
		Art    string
		Colors []string
	}{
		{artPaco, fireColors},
		{artSofiane, iceColors},
		{artValentin, campColors},
		{artAyman, mossColors},
	}
	clearScreen()
	banner("CRÉDITS", Gold)
	fmt.Println(Gray + Italic + "  Les quatre aventuriers qui ont forgé DunGo (aucun dragon n'a été blessé)" + Reset)
	fmt.Println()
	for _, author := range authors {
		printArtSlow(author.Art, author.Colors...)
		fmt.Println()
	}
	pause()
}

func goodbye() {
	clearScreen()
	printArt(artCampfire, fireColors...)
	fmt.Println()
	typewrite(Gold+Bold, "Merci d'avoir joué à DunGo ! Ignarok en profite pour faire une sieste.")
	fmt.Println()
}
