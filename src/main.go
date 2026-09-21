// DunGo — L'Antre d'Ignarok, un RPG en ligne de commande (Projet RED · Ymmersion).
//
// Le jeu se lit comme une pile d'écrans : titre → camp → donjon → combat.
// Chaque écran est une boucle « efface, affiche, lis un choix ».
package main

import (
	"fmt"
	"math/rand/v2"
)

var tips = []string{
	"Les loups et les corbeaux laissent souvent des ressources pour la forge.",
	"Les fontaines du donjon rendent la moitié de vos PV et de votre mana.",
	"Une Potion de poison lancée en combat inflige 30 dégâts au total.",
	"Quand Ignarok inspire, le Souffle Infernal arrive au tour suivant.",
	"Tous les 3 tours, les monstres frappent deux fois plus fort.",
	"Une arme faite pour votre lignée donne +2 attaque en bonus.",
	"Le tableau des missions propose des contrats bien payés. Jetez-y un œil !",
	"Second Souffle soigne aussi les saignements et les brûlures.",
	"La Boule de Feu fait brûler l'ennemi pendant 3 tours.",
	"Les boss ne vous laissent pas fuir. Préparez-vous avant la dernière salle !",
}

// L'écran titre : chaque entrée est un simple appel de fonction.
var titleActions = []struct {
	Label  string
	Color  string
	Action func()
}{
	{"Nouvelle partie", Green, newGame},
	{"Continuer", Sky, continueGame},
	{"Comment jouer", Yellow, howToPlay},
	{"Qui sont-ils ?", Purple, whoAreThey},
}

// Le menu du camp. L'ordre de la table donne les numéros affichés.
var campActions = []struct {
	Label  string
	Color  string
	Action func(*Character)
}{
	{"Informations du personnage", Sky, displayInfo},
	{"Inventaire", Brown, accessInventory},
	{"Marchand · Mordecai le Borgne", Gold, merchant},
	{"Forgeron · Borin Poing-de-Fer", Orange, blacksmith},
	{"Missions · Tableau du camp", Yellow, missionBoard},
	{"Entraînement · Arène du Sergent Grol", Red, trainingFight},
	{"Explorer le donjon", Purple, exploreDungeon},
	{"Sauvegarder", Silver, saveAndPause},
	{"Comment jouer", Gray, func(*Character) { howToPlay() }},
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
	fmt.Println(Orange + Bold + "              ~ L'ANTRE D'IGNAROK ~" + Reset)
	fmt.Println(Gray + Italic + "     « Sous les ruines de Karak-Dûm, le Fléau Écarlate s'éveille. »" + Reset)
	fmt.Println()
	for i, entry := range titleActions {
		option(i+1, entry.Label, entry.Color)
	}
	option(0, "Quitter", Gray)
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
		fail("Impossible de charger cette sauvegarde.")
		pause()
		return
	}
	success("Bon retour parmi nous, %s !", c.Name)
	wait(1000)
	campMenu(c)
}

func intro(c *Character) {
	clearScreen()
	printArt(artCampfire, fireColors...)
	banner("PROLOGUE", Gold)
	fmt.Println()
	typewrite(Silver, "Il y a mille ans, les nains de Karak-Dûm creusèrent trop profond et réveillèrent Ignarok, le Fléau Écarlate.")
	typewrite(Silver, "Aujourd'hui, la montagne tremble à nouveau. Douze héros sont partis tuer le dragon. Aucun n'est revenu.")
	typewrite(Gold+Bold, fmt.Sprintf("Vous êtes %s. Vous avez 50 Y-Coins, une potion et beaucoup trop de courage.", c.Name))
	pause()
}

// campMenu est la boucle principale d'une partie : on y revient entre
// chaque expédition dans le donjon.
func campMenu(c *Character) {
	for {
		clearScreen()
		banner("LE CAMP DE LA DERNIÈRE LUEUR", Gold)
		showStatus(c)
		fmt.Println(Gray + Italic + "  Astuce : " + tips[rand.IntN(len(tips))] + Reset)
		fmt.Println()
		for i, entry := range campActions {
			option(i+1, entry.Label, entry.Color)
		}
		option(0, "Quitter vers l'écran titre", Gray)

		choice := readChoice(0, len(campActions))
		if choice == 0 {
			if ask("Sauvegarder avant de quitter ?") {
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

// L'aide du jeu : un tableau de sections, faciles à compléter.
var helpSections = []struct {
	Title string
	Color string
	Lines []string
}{
	{"Commandes", Gold, []string{
		"Tapez le " + Gold + "numéro" + Reset + " d'une option puis " + Gold + "Entrée" + Reset + ".",
		Gold + "0" + Reset + " permet toujours de revenir en arrière.",
	}},
	{"But du jeu", Red, []string{
		"Traversez les 3 étages du donjon, battez leurs boss et terrassez " + Red + Bold + "Ignarok" + Reset + ".",
	}},
	{"Au camp", Orange, []string{
		Gold + "Marchand" + Reset + " : potions, grimoire, ressources (la 1re Potion de vie est offerte).",
		Orange + "Forgeron" + Reset + " : armures (+PV) et armes (+attaque), plus fortes pour votre lignée.",
		Yellow + "Missions" + Reset + " : des contrats de chasse récompensés en Y-Coins et en XP.",
		Red + "Arène" + Reset + "    : entraînement contre le gobelin, sans danger… et sans butin.",
	}},
	{"Dans le donjon", Purple, []string{
		"Chaque salle cache un monstre… ou une " + Cyan + "fontaine" + Reset + ", un " + Gold + "marchand ambulant" + Reset + " ou un " + Yellow + "sphinx" + Reset + ".",
		"La dernière salle de chaque étage contient un " + Red + Bold + "boss" + Reset + " : impossible de le fuir !",
	}},
	{"En combat", Red, []string{
		"Le plus rapide (" + Yellow + "initiative" + Reset + ") joue en premier. Les monstres frappent ×2 tous les 3 tours.",
		Red + "Saignement" + Reset + " et " + Orange + "brûlure" + Reset + " : dégâts à chaque tour.",
		Yellow + "Étourdi" + Reset + " : vous passez votre tour. " + Purple + "Affaibli" + Reset + " : vos dégâts sont divisés par 2.",
		Gold + "10 %" + Reset + " de chance de coup critique.",
		"Mort : résurrection au camp et perte de 20 % des Y-Coins.",
	}},
}

func howToPlay() {
	clearScreen()
	banner("COMMENT JOUER", Yellow)
	for _, help := range helpSections {
		section(help.Title, help.Color)
		for _, line := range help.Lines {
			fmt.Println("   " + line)
		}
	}
	pause()
}
func whoAreThey() {
	authors := []struct {
		Art    string
		Colors []string
	}{
		{artPaco, fireColors},
		{artSofiane, iceColors},
		{artValentin, goldColors},
		{artAyman, poisonColors},
	}
	clearScreen()
	banner("QUI SONT-ILS ?", Purple)
	fmt.Println(Gray + Italic + "  Les quatre aventuriers qui ont forgé DunGo" + Reset)
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
	typewrite(Gold+Bold, "Merci d'avoir joué à DunGo ! Le feu du camp vous attendra, aventurier.")
	fmt.Println()
}
