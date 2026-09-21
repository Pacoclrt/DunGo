package main

import (
	"fmt"
	"strings"
)

// Material : une ressource et la quantité qu'exige une recette.
type Material struct {
	Item     string
	Quantity int
}

// Recipe décrit un objet que Borin sait forger.
type Recipe struct {
	Item      string
	Cost      int
	Materials []Material
}

// Les recettes de Borin, rangées par étal. Chaque étal est un petit menu à
// part : l'écran de la forge n'affiche jamais plus de trois recettes à la fois.
var forgeShelves = []struct {
	Name    string
	Recipes []Recipe
}{
	{"Armures de l'aventurier", []Recipe{
		{ItemAdventurerHat, 12, []Material{{ItemRavenFeather, 1}, {ItemBoarLeather, 1}}},
		{ItemAdventurerTunic, 12, []Material{{ItemWolfFur, 2}, {ItemTrollSkin, 1}}},
		{ItemAdventurerBoots, 12, []Material{{ItemWolfFur, 1}, {ItemBoarLeather, 1}}},
	}},
	{"Armures de maître", []Recipe{
		{ItemTrollHelm, 30, []Material{{ItemTrollSkin, 2}, {ItemRavenFeather, 2}}},
		{ItemBoarCuirass, 30, []Material{{ItemBoarLeather, 3}, {ItemTrollSkin, 2}}},
		{ItemWerewolfBoots, 30, []Material{{ItemWolfFur, 3}, {ItemBoarLeather, 1}}},
	}},
	{"Armes de l'aventurier", []Recipe{
		{ItemAdventurerSword, 15, []Material{{ItemBoarLeather, 2}, {ItemTrollSkin, 1}}},
		{ItemAdventurerBow, 15, []Material{{ItemRavenFeather, 2}, {ItemWolfFur, 1}}},
		{ItemAdventurerHammer, 15, []Material{{ItemTrollSkin, 1}, {ItemWolfFur, 1}}},
	}},
	{"Armes de maître", []Recipe{
		{ItemKnightBlade, 35, []Material{{ItemTrollSkin, 3}, {ItemBoarLeather, 2}}},
		{ItemSylvanBow, 35, []Material{{ItemRavenFeather, 4}, {ItemWolfFur, 2}}},
		{ItemRunicHammer, 35, []Material{{ItemTrollSkin, 3}, {ItemWolfFur, 2}}},
	}},
}

func blacksmith(c *Character) {
	greeting := "Apportez-moi des peaux, je vous rends une armure. Ou une arme. Je ne suis pas compliqué."
	for {
		clearScreen()
		printArt(artBlacksmith, campColors...)
		banner("LE FORGERON", Gold)
		say("Borin", Gold, greeting)
		greeting = "Le fer est encore chaud. Autre chose ?"
		showStatus(c)

		section("Que voulez-vous forger ?")
		for i, shelf := range forgeShelves {
			optionHint(i+1, shelf.Name, readyText(c, shelf.Recipes))
		}
		back("Retour au camp")

		choice := readChoice(0, len(forgeShelves))
		if choice == 0 {
			return
		}
		forgeShelf(c, forgeShelves[choice-1].Name, forgeShelves[choice-1].Recipes)
	}
}

// readyText annonce combien de recettes de l'étal sont faisables tout de suite.
func readyText(c *Character, recipes []Recipe) string {
	ready := 0
	for _, recipe := range recipes {
		if canCraft(c, recipe) {
			ready++
		}
	}
	switch ready {
	case 0:
		return ""
	case 1:
		return Green + "√ 1 recette prête" + Reset
	}
	return Green + fmt.Sprintf("√ %d recettes prêtes", ready) + Reset
}

// forgeShelf affiche les recettes d'un seul étal et fabrique celle choisie.
func forgeShelf(c *Character, name string, recipes []Recipe) {
	for {
		clearScreen()
		banner(strings.ToUpper(name), Gold)
		showStatus(c)
		for i, recipe := range recipes {
			ready := Green + "√ " + Reset
			if !canCraft(c, recipe) {
				ready = Red + "× " + Reset
			}
			fmt.Println()
			itemLine(i+1, recipe.Item, itemColor(recipe.Item)+Bold, fmt.Sprintf("%s%s%s · %s · %s%d Y-Coins%s",
				Orange, itemBonusText(recipe.Item), Reset, itemSlotText(recipe.Item), Gold, recipe.Cost, Reset))
			fmt.Println("         " + ready + materialsText(c, recipe))
		}
		fmt.Println()
		back("Retour")

		choice := readChoice(0, len(recipes))
		if choice == 0 {
			return
		}
		fmt.Println()
		craft(c, recipes[choice-1])
		pause()
	}
}

// materialsText liste les ressources d'une recette, en vert si on les a.
func materialsText(c *Character, recipe Recipe) string {
	text := ""
	for i, material := range recipe.Materials {
		if i > 0 {
			text += Gray + " + " + Reset
		}
		color := Green
		if c.Inventory[material.Item] < material.Quantity {
			color = Red
		}
		text += fmt.Sprintf("%s%d %s (%d/%d)%s", color, material.Quantity, material.Item,
			c.Inventory[material.Item], material.Quantity, Reset)
	}
	return text
}

func canCraft(c *Character, recipe Recipe) bool {
	if c.Gold < recipe.Cost {
		return false
	}
	for _, material := range recipe.Materials {
		if c.Inventory[material.Item] < material.Quantity {
			return false
		}
	}
	return true
}

func craft(c *Character, recipe Recipe) {
	if c.Gold < recipe.Cost {
		fail("Il vous manque %d Y-Coins. La forge ne chauffe pas gratis !", recipe.Cost-c.Gold)
		return
	}
	used := 0
	for _, material := range recipe.Materials {
		if c.Inventory[material.Item] < material.Quantity {
			fail("Il vous manque : %d %s (vous en avez %d).", material.Quantity, material.Item, c.Inventory[material.Item])
			say("Borin", Gold, "Je forge du métal, pas des promesses.")
			return
		}
		used += material.Quantity
	}
	// Les ressources consommées libèrent de la place pour l'équipement forgé.
	if inventoryCount(c)-used+1 > c.InventoryMax {
		fail("Votre sac est plein, et Borin refuse de porter vos affaires.")
		return
	}

	for _, material := range recipe.Materials {
		removeInventory(c, material.Item, material.Quantity)
	}
	c.Gold -= recipe.Cost
	addInventory(c, recipe.Item, 1)

	dots("Borin frappe le métal", Orange)
	if isWeapon(recipe.Item) {
		printArt(artSword, stoneColors...)
	} else {
		printArt(artAnvil, fireColors...)
	}
	success("Borin vous tend : %s (%s). Encore tiède !", recipe.Item, itemBonusText(recipe.Item))
	info("Équipez-le depuis l'inventaire.")
}
