// ════════════════════════════════════════════════════════════════════════
//   merchant.go · [MARCHAND]
//   Mordecai, le marchand du camp : acheter et vendre des objets.
//   Les prix sont rangés dans la table items (items.go).
// ════════════════════════════════════════════════════════════════════════

package main

import "fmt"

var shopItems = []string{ // [MARCHAND] sert à lister les objets vendus par le marchand, dans l'ordre d'affichage
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

func merchantMenu(character *Character) { // [MARCHAND] sert à afficher la boutique (Acheter / Vendre) jusqu'au retour au camp
	for {
		clearScreen()
		printArt(artMerchant, campColors...)
		banner("LE MARCHAND", Gold)
		say("Mordecai", Gold, "Bienvenue ! Tout est à vendre, sauf mon œil. Un troll l'a déjà pris.")
		fmt.Printf("\n   %s¤ Bourse : %d Y-Coins%s\n", Gold+Bold, character.YCoins, Reset)
		section("Que faire ?")
		option(1, "Acheter")
		option(2, "Vendre")
		backOption("Retour au camp")

		switch readChoice(0, 2) {
		case 0:
			return
		case 1:
			buyMenu(character)
		case 2:
			sellMenu(character)
		}
	}
}

// ════════════════════════════════════════════════════════════════════════
//   ACHETER
// ════════════════════════════════════════════════════════════════════════

func buyPrice(character *Character, item string) int { // [MARCHAND] sert à donner le prix d'achat d'un objet (0 pour la première Potion de vie, offerte) et le renvoie
	if item == ItemHealthPotion && !character.FreePotionTaken {
		return 0
	}
	return items[item].Price
}

func buyMenu(character *Character) { // [MARCHAND] sert à afficher les objets à vendre avec leur prix et à acheter celui choisi, jusqu'au Retour
	for {
		clearScreen()
		banner("ACHETER", Gold)
		showCharacterBar(character)
		fmt.Println()
		for index, item := range shopItems {
			price := buyPrice(character, item)
			if price == 0 {
				itemLine(index+1, item, Green+"GRATUIT !"+Reset)
			} else {
				itemLine(index+1, item, fmt.Sprintf("%s%d Y-Coins%s", Gold, price, Reset))
			}
		}
		fmt.Println()
		backOption("Retour")

		choice := readChoice(0, len(shopItems))
		if choice == 0 {
			return
		}
		fmt.Println()
		buyItem(character, shopItems[choice-1])
		pause()
	}
}

func buyItem(character *Character, item string) { // [MARCHAND] sert à acheter un objet : vérifie les Y-Coins et la place dans le sac, puis paie
	price := buyPrice(character, item)
	if character.YCoins < price {
		fail("Il vous manque %d Y-Coins. Mordecai ne fait pas crédit.", price-character.YCoins)
		return
	}

	// L'augmentation d'inventaire n'est pas un objet : elle agrandit le sac
	// tout de suite (+10 places, 3 fois au maximum).
	if item == ItemInventoryUpgrade {
		if character.InventoryUpgrades >= 3 {
			fail("Votre sac a déjà été agrandi 3 fois.")
			return
		}
		character.InventoryUpgrades++
		character.InventoryMax += 10
		character.YCoins -= price
		success("Votre sac peut maintenant contenir %d objets.", character.InventoryMax)
		return
	}

	// Les autres objets vont dans le sac : on ne paie que s'il y a la place.
	if !addToInventory(character, item, 1) {
		return
	}
	character.YCoins -= price
	if item == ItemHealthPotion {
		character.FreePotionTaken = true // la prochaine Potion de vie sera payante
	}
	printArt(artCoins, campColors...)
	success("Vous achetez : %s (-%d Y-Coins).", item, price)
}

// ════════════════════════════════════════════════════════════════════════
//   VENDRE
// ════════════════════════════════════════════════════════════════════════

func sellPrice(item string) int { // [MARCHAND] sert à donner le prix de revente d'un objet (équipement : PV + 4 × attaque ; le reste : moitié du prix) et le renvoie
	if equipment, isGear := gear[item]; isGear {
		return equipment.BonusHP + equipment.BonusAttack*4
	}
	return items[item].Price / 2
}

func sellMenu(character *Character) { // [MARCHAND] sert à afficher le sac avec le prix de revente et à vendre l'objet choisi, jusqu'au Retour
	for {
		clearScreen()
		banner("VENDRE", Gold)
		showCharacterBar(character)
		fmt.Println()

		list := sortedInventory(character)
		if len(list) == 0 {
			info("Votre sac est vide.")
		}
		for index, item := range list {
			itemLine(index+1, item, fmt.Sprintf("x%-3d %s+%d Y-Coins%s", character.Inventory[item], Gold, sellPrice(item), Reset))
		}
		fmt.Println()
		backOption("Retour")

		choice := readChoice(0, len(list))
		if choice == 0 {
			return
		}
		fmt.Println()
		sellItem(character, list[choice-1])
		pause()
	}
}

func sellItem(character *Character, item string) { // [MARCHAND] sert à vendre un objet du sac : il est retiré et le héros reçoit son prix de revente
	removeFromInventory(character, item, 1)
	character.YCoins += sellPrice(item)
	success("Vendu : %s (+%d Y-Coins).", item, sellPrice(item))
}
