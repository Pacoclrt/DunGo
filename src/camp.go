package main

import (
	"fmt"
	"strings"
)

// --- Marchand ---

var shopItems = []string{
	ItemHealthPotion,
	ItemManaPotion,
	ItemPoisonPotion,
	ItemFireballBook,
	ItemWolfFur,
	ItemTrollSkin,
	ItemBoarLeather,
	ItemRavenFeather,
	ItemInventoryUpgrade,
}

func merchant(c *Character) {
	for {
		clearScreen()
		printArt(artMerchant, campColors...)
		banner("LE MARCHAND", Gold)
		say("Mordecai", Gold, "Bienvenue ! Tout est à vendre, sauf mon œil. Un troll l'a déjà pris.")
		fmt.Printf("\n   %s¤ Bourse : %d Y-Coins%s\n", Gold+Bold, c.Gold, Reset)
		section("Que faire ?")
		option(1, "Acheter")
		option(2, "Vendre")
		back("Retour au camp")

		switch readChoice(0, 2) {
		case 0:
			return
		case 1:
			buyMenu(c)
		case 2:
			sellMenu(c)
		}
	}
}

// La première Potion de vie est offerte.
func priceOf(c *Character, item string) int {
	if item == ItemHealthPotion && !c.FreePotionTaken {
		return 0
	}
	return items[item].Price
}

func buyMenu(c *Character) {
	for {
		clearScreen()
		banner("ACHETER", Gold)
		showStatus(c)
		fmt.Println()
		for i, item := range shopItems {
			price := priceOf(c, item)
			if price == 0 {
				itemLine(i+1, item, Green+"GRATUIT !"+Reset)
			} else {
				itemLine(i+1, item, fmt.Sprintf("%s%d Y-Coins%s", Gold, price, Reset))
			}
		}
		fmt.Println()
		back("Retour")

		choice := readChoice(0, len(shopItems))
		if choice == 0 {
			return
		}
		fmt.Println()
		buyItem(c, shopItems[choice-1])
		pause()
	}
}

func buyItem(c *Character, item string) {
	price := priceOf(c, item)
	if c.Gold < price {
		fail("Il vous manque %d Y-Coins. Mordecai ne fait pas crédit.", price-c.Gold)
		return
	}

	// L'augmentation d'inventaire n'est pas un objet : elle agrandit le sac.
	if item == ItemInventoryUpgrade {
		if c.InventoryUpgrades >= 3 {
			fail("Votre sac a déjà été agrandi 3 fois.")
			return
		}
		c.InventoryUpgrades++
		c.InventoryMax += 10
		c.Gold -= price
		success("Votre sac peut maintenant contenir %d objets.", c.InventoryMax)
		return
	}

	if !addInventory(c, item, 1) {
		return
	}
	c.Gold -= price
	if item == ItemHealthPotion {
		c.FreePotionTaken = true
	}
	printArt(artCoins, campColors...)
	success("Vous achetez : %s (-%d Y-Coins).", item, price)
}

// Un équipement se revend selon ses bonus, le reste à la moitié de son prix.
func sellPrice(item string) int {
	if g, isGear := gear[item]; isGear {
		return g.HP + g.Attack*4
	}
	return items[item].Price / 2
}

func sellMenu(c *Character) {
	for {
		clearScreen()
		banner("VENDRE", Gold)
		showStatus(c)
		fmt.Println()

		list := sortedItems(c)
		if len(list) == 0 {
			info("Votre sac est vide.")
		}
		for i, item := range list {
			itemLine(i+1, item, fmt.Sprintf("x%-3d %s+%d Y-Coins%s", c.Inventory[item], Gold, sellPrice(item), Reset))
		}
		fmt.Println()
		back("Retour")

		choice := readChoice(0, len(list))
		if choice == 0 {
			return
		}
		item := list[choice-1]
		removeInventory(c, item, 1)
		c.Gold += sellPrice(item)
		fmt.Println()
		success("Vendu : %s (+%d Y-Coins).", item, sellPrice(item))
		pause()
	}
}

// --- Forgeron ---

type Material struct {
	Item     string
	Quantity int
}

type Recipe struct {
	Item      string
	Cost      int
	Materials []Material
}

// Shelf est un étal de la forge : un petit groupe de recettes.
type Shelf struct {
	Name    string
	Recipes []Recipe
}

var shelves = []Shelf{
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
	for {
		clearScreen()
		printArt(artBlacksmith, campColors...)
		banner("LE FORGERON", Gold)
		say("Borin", Gold, "Apportez-moi des peaux, je vous rends une armure. Ou une arme. Je ne suis pas compliqué.")
		showStatus(c)
		section("Que voulez-vous forger ?")
		for i, shelf := range shelves {
			option(i+1, shelf.Name)
		}
		back("Retour au camp")

		choice := readChoice(0, len(shelves))
		if choice == 0 {
			return
		}
		forgeShelf(c, shelves[choice-1])
	}
}

func forgeShelf(c *Character, shelf Shelf) {
	for {
		clearScreen()
		banner(strings.ToUpper(shelf.Name), Gold)
		showStatus(c)
		for i, recipe := range shelf.Recipes {
			fmt.Println()
			itemLine(i+1, recipe.Item, fmt.Sprintf("%s · %s%d Y-Coins%s", gearText(recipe.Item), Gold, recipe.Cost, Reset))
			fmt.Println("         " + materialsText(c, recipe))
		}
		fmt.Println()
		back("Retour")

		choice := readChoice(0, len(shelf.Recipes))
		if choice == 0 {
			return
		}
		fmt.Println()
		craft(c, shelf.Recipes[choice-1])
		pause()
	}
}

func materialsText(c *Character, recipe Recipe) string {
	parts := []string{}
	for _, material := range recipe.Materials {
		owned := c.Inventory[material.Item]
		color := Green
		if owned < material.Quantity {
			color = Red
		}
		parts = append(parts, fmt.Sprintf("%s%d %s (%d/%d)%s", color, material.Quantity, material.Item, owned, material.Quantity, Reset))
	}
	return strings.Join(parts, " + ")
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
	if !canCraft(c, recipe) {
		fail("Il manque des Y-Coins ou des ressources (en rouge).")
		return
	}
	for _, material := range recipe.Materials {
		removeInventory(c, material.Item, material.Quantity)
	}
	c.Gold -= recipe.Cost
	addInventory(c, recipe.Item, 1)

	printArt(artAnvil, fireColors...)
	success("Borin vous tend : %s. Équipez-le depuis l'inventaire.", recipe.Item)
}

// --- Missions ---

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

// missionBoard affiche la mission en cours. Les missions se font une à une,
// dans l'ordre : QuestIndex est le numéro de la mission en cours.
func missionBoard(c *Character) {
	clearScreen()
	printArt(artQuestBoard, campColors...)
	banner("MISSIONS", Gold)

	if c.QuestIndex >= len(quests) {
		fmt.Println()
		success("Toutes les missions sont terminées. Le camp vous doit une fière chandelle !")
		pause()
		return
	}

	quest := quests[c.QuestIndex]
	section(fmt.Sprintf("Mission %d / %d : %s", c.QuestIndex+1, len(quests), quest.Title))
	paragraph(Italic, quest.Description)
	fmt.Printf("   Objectif : vaincre %d × %s\n", quest.Goal, quest.Target)
	fmt.Printf("   Récompense : %s%d Y-Coins%s et %s%d XP%s\n", Gold, quest.Reward, Reset, Purple, quest.XP, Reset)
	if c.QuestActive {
		fmt.Printf("   Progression : %d / %d\n", c.QuestProgress, quest.Goal)
	}
	fmt.Println()

	if !c.QuestActive {
		option(1, "Accepter la mission")
		back("Retour au camp")
		if readChoice(0, 1) == 1 {
			c.QuestActive = true
			c.QuestProgress = 0
			success("Mission acceptée !")
			pause()
		}
		return
	}

	if c.QuestProgress < quest.Goal {
		info("Mission en cours. Revenez quand l'objectif sera atteint.")
		pause()
		return
	}

	option(1, "Réclamer la récompense")
	back("Retour au camp")
	if readChoice(0, 1) == 1 {
		printArt(artCoins, campColors...)
		c.Gold += quest.Reward
		success("Mission « %s » terminée ! +%d Y-Coins", quest.Title, quest.Reward)
		gainXP(c, quest.XP)
		c.QuestIndex++
		c.QuestActive = false
		c.QuestProgress = 0
		pause()
	}
}

// updateQuest est appelée après chaque victoire pour faire avancer la mission.
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
}
