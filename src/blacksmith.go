// ════════════════════════════════════════════════════════════════════════
//   blacksmith.go · [FORGERON]
//   Borin, le forgeron du camp : il fabrique des armures et des armes à
//   partir des ressources laissées par les monstres, contre des Y-Coins.
//   Une recette = un objet, un prix, et une liste de matériaux.
// ════════════════════════════════════════════════════════════════════════

package main

import (
	"fmt"
	"strings"
)

// ════════════════════════════════════════════════════════════════════════
//   LES RECETTES
// ════════════════════════════════════════════════════════════════════════

type Material struct { // [FORGERON] sert à décrire un matériau d'une recette : quelle ressource et combien
	Item     string
	Quantity int
}

type Recipe struct { // [FORGERON] sert à décrire une recette : l'objet fabriqué, son prix en Y-Coins et ses matériaux
	Item      string
	Cost      int // en Y-Coins
	Materials []Material
}

type Shelf struct { // [FORGERON] sert à décrire un étal de la forge : un nom et un petit groupe de recettes
	Name    string
	Recipes []Recipe
}

var shelves = []Shelf{ // [FORGERON] sert à lister les 4 étals de la forge et leurs recettes (3 chacun)
	{Name: "Armures de l'aventurier", Recipes: []Recipe{
		{Item: ItemAdventurerHat, Cost: 12, Materials: []Material{{ItemRavenFeather, 1}, {ItemBoarLeather, 1}}},
		{Item: ItemAdventurerTunic, Cost: 12, Materials: []Material{{ItemWolfFur, 2}, {ItemTrollSkin, 1}}},
		{Item: ItemAdventurerBoots, Cost: 12, Materials: []Material{{ItemWolfFur, 1}, {ItemBoarLeather, 1}}},
	}},
	{Name: "Armures de maître", Recipes: []Recipe{
		{Item: ItemTrollHelm, Cost: 30, Materials: []Material{{ItemTrollSkin, 2}, {ItemRavenFeather, 2}}},
		{Item: ItemBoarCuirass, Cost: 30, Materials: []Material{{ItemBoarLeather, 3}, {ItemTrollSkin, 2}}},
		{Item: ItemWerewolfBoots, Cost: 30, Materials: []Material{{ItemWolfFur, 3}, {ItemBoarLeather, 1}}},
	}},
	{Name: "Armes de l'aventurier", Recipes: []Recipe{
		{Item: ItemAdventurerSword, Cost: 15, Materials: []Material{{ItemBoarLeather, 2}, {ItemTrollSkin, 1}}},
		{Item: ItemAdventurerBow, Cost: 15, Materials: []Material{{ItemRavenFeather, 2}, {ItemWolfFur, 1}}},
		{Item: ItemAdventurerHammer, Cost: 15, Materials: []Material{{ItemTrollSkin, 1}, {ItemWolfFur, 1}}},
	}},
	{Name: "Armes de maître", Recipes: []Recipe{
		{Item: ItemKnightBlade, Cost: 35, Materials: []Material{{ItemTrollSkin, 3}, {ItemBoarLeather, 2}}},
		{Item: ItemSylvanBow, Cost: 35, Materials: []Material{{ItemRavenFeather, 4}, {ItemWolfFur, 2}}},
		{Item: ItemRunicHammer, Cost: 35, Materials: []Material{{ItemTrollSkin, 3}, {ItemWolfFur, 2}}},
	}},
}

// ════════════════════════════════════════════════════════════════════════
//   LES ÉCRANS DE LA FORGE
// ════════════════════════════════════════════════════════════════════════

func blacksmithMenu(character *Character) { // [FORGERON] sert à afficher les 4 étals de la forge et à ouvrir celui choisi, jusqu'au retour au camp
	for {
		clearScreen()
		printArt(artBlacksmith, campColors...)
		banner("LE FORGERON", Gold)
		say("Borin", Gold, "Apportez-moi des peaux, je vous rends une armure. Ou une arme. Je ne suis pas compliqué.")
		showCharacterBar(character)
		section("Que voulez-vous forger ?")
		for index, shelf := range shelves {
			option(index+1, shelf.Name)
		}
		backOption("Retour au camp")

		choice := readChoice(0, len(shelves))
		if choice == 0 {
			return
		}
		shelfMenu(character, shelves[choice-1])
	}
}

func shelfMenu(character *Character, shelf Shelf) { // [FORGERON] sert à afficher les recettes d'un étal (bonus, prix, matériaux) et à fabriquer celle choisie, jusqu'au Retour
	for {
		clearScreen()
		banner(strings.ToUpper(shelf.Name), Gold)
		showCharacterBar(character)
		for index, recipe := range shelf.Recipes {
			fmt.Println()
			itemLine(index+1, recipe.Item, fmt.Sprintf("%s · %s%d Y-Coins%s", gearBonusText(recipe.Item), Gold, recipe.Cost, Reset))
			fmt.Println("         " + materialsText(character, recipe))
		}
		fmt.Println()
		backOption("Retour")

		choice := readChoice(0, len(shelf.Recipes))
		if choice == 0 {
			return
		}
		fmt.Println()
		craft(character, shelf.Recipes[choice-1])
		pause()
	}
}

func materialsText(character *Character, recipe Recipe) string { // [FORGERON] sert à écrire les matériaux d'une recette : en vert ceux qu'on a, en rouge ceux qui manquent (ex : 2 Fourrure de Loup (1/2)), et renvoie le texte
	parts := []string{}
	for _, material := range recipe.Materials {
		owned := character.Inventory[material.Item] // 0 si l'objet n'est pas dans le sac
		color := Green
		if owned < material.Quantity {
			color = Red
		}
		parts = append(parts, fmt.Sprintf("%s%d %s (%d/%d)%s", color, material.Quantity, material.Item, owned, material.Quantity, Reset))
	}
	return strings.Join(parts, " + ")
}

// ════════════════════════════════════════════════════════════════════════
//   LA FABRICATION
// ════════════════════════════════════════════════════════════════════════

func canCraft(character *Character, recipe Recipe) bool { // [FORGERON] sert à vérifier que le héros a assez de Y-Coins et tous les matériaux ; renvoie true s'il peut fabriquer
	if character.YCoins < recipe.Cost {
		return false
	}
	for _, material := range recipe.Materials {
		if character.Inventory[material.Item] < material.Quantity {
			return false
		}
	}
	return true
}

func craft(character *Character, recipe Recipe) { // [FORGERON] sert à fabriquer un objet : retire les matériaux et les Y-Coins, puis range l'objet dans le sac
	if !canCraft(character, recipe) {
		fail("Il manque des Y-Coins ou des ressources (en rouge).")
		return
	}
	for _, material := range recipe.Materials {
		removeFromInventory(character, material.Item, material.Quantity)
	}
	character.YCoins -= recipe.Cost
	// Il y a forcément la place : on vient de retirer au moins 2 matériaux du sac.
	addToInventory(character, recipe.Item, 1)

	printArt(artAnvil, fireColors...)
	success("Borin vous tend : %s. Équipez-le depuis l'inventaire.", recipe.Item)
}
