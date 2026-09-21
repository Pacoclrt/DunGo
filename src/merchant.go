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
}

func merchant(c *Character) {
	greeting := merchantGreetings[rand.IntN(len(merchantGreetings))]
	for {
		clearScreen()
		printArt(artMerchant, goldColors...)
		banner("L'ÉCHOPPE DE MORDECAI LE BORGNE", Gold)
		say("Mordecai", Gold, greeting)
		greeting = "Autre chose pour votre service ?"
		fmt.Printf("\n   %s¤ Bourse : %d Y-Coins%s\n\n", Gold+Bold, c.Gold, Reset)
		option(1, "Acheter", Green)
		option(2, "Vendre", Orange)
		option(0, "Retour au camp", Gray)

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
		banner("ACHETER", Green)
		showStatus(c)
		shown := ""
		for i, item := range shopItems {
			if category := itemCategory(item); category != shown {
				section(category, itemColor(item))
				shown = category
			}
			itemLine(i+1, item, itemColor(item), priceTag(c, item))
		}
		fmt.Println()
		option(0, "Retour", Gray)

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
		return Lime + Bold + "GRATUIT !" + Reset
	}
	color := Gold
	if c.Gold < items[item].Price {
		color = DarkRed
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
		fail("Il vous manque %d Y-Coins.", price-c.Gold)
		say("Mordecai", Gold, "Pas de Y-Coins, pas d'objet, l'ami !")
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
		success("Cadeau de bienvenue ! Vous recevez : %s", item)
		say("Mordecai", Gold, "La première est offerte. Les suivantes, beaucoup moins.")
		return
	}
	printArt(artCoins, goldColors...)
	success("Vous avez acheté : %s (-%d Y-Coins). Il vous reste %d Y-Coins.", item, price, c.Gold)
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
		banner("VENDRE", Orange)
		showStatus(c)
		list := showItems(c, func(item string) string {
			return fmt.Sprintf("x%-3d %s+%d Y-Coins%s", c.Inventory[item], Gold, sellPrice(item), Reset)
		})
		fmt.Println()
		option(0, "Retour", Gray)

		choice := readChoice(0, len(list))
		if choice == 0 {
			return
		}
		item := list[choice-1]
		removeInventory(c, item, 1)
		c.Gold += sellPrice(item)
		fmt.Println()
		success("Vous vendez %s (+%d Y-Coins).", item, sellPrice(item))
		say("Mordecai", Gold, "Un prix honnête. Pour moi, surtout.")
		pause()
	}
}
