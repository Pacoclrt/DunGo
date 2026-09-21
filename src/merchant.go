package main

import (
	"fmt"
	"math/rand/v2"
)

// Ce que Mordecai propose, dans l'ordre d'affichage. Les prix vivent dans
// la table items, avec le reste de la description des objets.
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

var merchantGreetings = []string{
	"Bienvenue ! Tout est à vendre, sauf mon œil. Un troll l'a déjà pris.",
	"Des potions, des grimoires, des fourrures… et zéro remboursement.",
	"Vous allez voir le dragon ? Prenez des potions. Beaucoup de potions.",
	"Ah, un client vivant ! Ça change des derniers.",
	"Satisfait ou… non, en fait, pas remboursé du tout.",
}

func merchant(c *Character) {
	greeting := merchantGreetings[rand.IntN(len(merchantGreetings))]
	for {
		clearScreen()
		printArt(artMerchant, campColors...)
		banner("LE MARCHAND", Gold)
		say("Mordecai", Gold, greeting)
		greeting = "Autre chose ? J'ai tout mon temps. Et vous, toutes vos pièces."
		fmt.Printf("\n   %s¤ Bourse : %d Y-Coins%s\n", Gold+Bold, c.Gold, Reset)
		section("Que faire ?")
		optionHint(1, "Acheter", "potions, grimoire, ressources")
		optionHint(2, "Vendre", "vider un peu votre sac")
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

func buyMenu(c *Character) {
	for {
		clearScreen()
		banner("ACHETER", Gold)
		showStatus(c)
		shown := ""
		for i, item := range shopItems {
			if category := itemCategory(item); category != shown {
				section(category)
				shown = category
			}
			itemLine(i+1, item, itemColor(item), priceTag(c, item))
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

// priceTag affiche le prix, en vert s'il est offert, en rouge s'il est trop cher.
func priceTag(c *Character, item string) string {
	if item == ItemHealthPotion && !c.FreePotionTaken {
		return Green + Bold + "GRATUIT !" + Reset
	}
	color := Gold
	if c.Gold < items[item].Price {
		color = DarkGray
	}
	return fmt.Sprintf("%s%d Y-Coins%s", color, items[item].Price, Reset)
}

func buyItem(c *Character, item string) {
	price := items[item].Price
	free := item == ItemHealthPotion && !c.FreePotionTaken
	if free {
		price = 0 // cadeau de bienvenue
	}
	if c.Gold < price {
		fail("Il vous manque %d Y-Coins. Mordecai ne fait pas crédit.", price-c.Gold)
		return
	}

	if item == ItemInventoryUpgrade {
		if upgradeInventorySlot(c) {
			c.Gold -= price
		}
		return
	}
	if !addInventory(c, item, 1) {
		return
	}
	c.Gold -= price
	if free {
		c.FreePotionTaken = true
		success("Cadeau de bienvenue : %s !", item)
		say("Mordecai", Gold, "La première est offerte. Les suivantes, beaucoup moins.")
		return
	}
	printArt(artCoins, campColors...)
	success("Vous achetez : %s (-%d Y-Coins). Mordecai jubile.", item, price)
}

// sellPrice : un équipement vaut ses bonus, le reste la moitié de son prix.
func sellPrice(item string) int {
	if g, ok := gear[item]; ok {
		return g.HP + g.Attack*4
	}
	return items[item].Price / 2
}

func sellMenu(c *Character) {
	for {
		clearScreen()
		banner("VENDRE", Gold)
		showStatus(c)
		list := showItems(c, func(item string) string {
			return fmt.Sprintf("x%-3d %s+%d Y-Coins%s", c.Inventory[item], Gold, sellPrice(item), Reset)
		})
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
		say("Mordecai", Gold, "Un prix honnête. Pour moi, surtout.")
		pause()
	}
}
