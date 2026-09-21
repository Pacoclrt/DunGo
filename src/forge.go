package main

import "fmt"

// Material : une ressource et la quantité qu'exige une recette.
type Material struct {
	Item     string
	Quantity int
}

// Recipe décrit ce que Borin sait forger. Section n'est rempli que sur la
// première recette d'un groupe : elle sert de titre dans la liste.
type Recipe struct {
	Section   string
	Item      string
	Cost      int
	Materials []Material
}

var recipes = []Recipe{
	{"Armures de l'aventurier", ItemAdventurerHat, 12, []Material{{ItemRavenFeather, 1}, {ItemBoarLeather, 1}}},
	{"", ItemAdventurerTunic, 12, []Material{{ItemWolfFur, 2}, {ItemTrollSkin, 1}}},
	{"", ItemAdventurerBoots, 12, []Material{{ItemWolfFur, 1}, {ItemBoarLeather, 1}}},
	{"Armures de maître", ItemTrollHelm, 30, []Material{{ItemTrollSkin, 2}, {ItemRavenFeather, 2}}},
	{"", ItemBoarCuirass, 30, []Material{{ItemBoarLeather, 3}, {ItemTrollSkin, 2}}},
	{"", ItemWerewolfBoots, 30, []Material{{ItemWolfFur, 3}, {ItemBoarLeather, 1}}},
	{"Armes de l'aventurier", ItemAdventurerSword, 15, []Material{{ItemBoarLeather, 2}, {ItemTrollSkin, 1}}},
	{"", ItemAdventurerBow, 15, []Material{{ItemRavenFeather, 2}, {ItemWolfFur, 1}}},
	{"", ItemAdventurerHammer, 15, []Material{{ItemTrollSkin, 1}, {ItemWolfFur, 1}}},
	{"Armes de maître", ItemKnightBlade, 35, []Material{{ItemTrollSkin, 3}, {ItemBoarLeather, 2}}},
	{"", ItemSylvanBow, 35, []Material{{ItemRavenFeather, 4}, {ItemWolfFur, 2}}},
	{"", ItemRunicHammer, 35, []Material{{ItemTrollSkin, 3}, {ItemWolfFur, 2}}},
}

func blacksmith(c *Character) {
	greeting := "Apportez-moi des peaux, je vous rends une armure. Ou une arme, si vous préférez."
	for {
		clearScreen()
		printArt(artBlacksmith, fireColors...)
		banner("LA FORGE DE BORIN POING-DE-FER", Orange)
		say("Borin", Orange, greeting)
		greeting = "Le fer est encore chaud. Autre chose ?"
		showStatus(c)

		for i, recipe := range recipes {
			if recipe.Section != "" {
				section(recipe.Section, Silver)
			}
			ready := Green + "√ prêt" + Reset
			if !canCraft(c, recipe) {
				ready = DarkRed + "× manquant" + Reset
			}
			itemLine(i+1, recipe.Item, Silver+Bold, fmt.Sprintf("%s%-11s%s %s(%s)%s  %s%d Y-Coins%s  %s",
				Green, itemBonusText(recipe.Item), Reset,
				Gray, itemSlotText(recipe.Item), Reset, Gold, recipe.Cost, Reset, ready))
			fmt.Println("         " + materialsText(c, recipe))
		}
		fmt.Println()
		option(0, "Retour au camp", Gray)

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
		fail("Il vous manque %d Y-Coins.", recipe.Cost-c.Gold)
		say("Borin", Orange, "La forge ne chauffe pas gratis !")
		return
	}
	used := 0
	for _, material := range recipe.Materials {
		if c.Inventory[material.Item] < material.Quantity {
			fail("Ressources manquantes : %d %s (vous en avez %d).", material.Quantity, material.Item, c.Inventory[material.Item])
			say("Borin", Orange, "Sans les bons matériaux, même moi je ne fais pas de miracles.")
			return
		}
		used += material.Quantity
	}
	// Les ressources consommées libèrent de la place pour l'équipement forgé.
	if inventoryCount(c)-used+1 > c.InventoryMax {
		fail("Pas de place dans votre sac pour l'équipement !")
		return
	}

	for _, material := range recipe.Materials {
		removeInventory(c, material.Item, material.Quantity)
	}
	c.Gold -= recipe.Cost
	addInventory(c, recipe.Item, 1)

	dots("Borin frappe le métal", Orange)
	if isWeapon(recipe.Item) {
		printArt(artSword, steelColors...)
	} else {
		printArt(artAnvil, fireColors...)
	}
	success("Borin vous tend : %s (%s) !", recipe.Item, itemBonusText(recipe.Item))
	info("Équipez-le depuis l'inventaire.")
}
