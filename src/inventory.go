package main

import (
	"fmt"
	"sort"
)

const (
	ItemHealthPotion     = "Potion de vie"
	ItemManaPotion       = "Potion de mana"
	ItemPoisonPotion     = "Potion de poison"
	ItemFireballBook     = "Livre de Sort : Boule de Feu"
	ItemInventoryUpgrade = "Augmentation d'inventaire"
	ItemWolfFur          = "Fourrure de Loup"
	ItemTrollSkin        = "Peau de Troll"
	ItemBoarLeather      = "Cuir de Sanglier"
	ItemRavenFeather     = "Plume de Corbeau"
)

// Ordre d'affichage des catégories dans le sac et chez le marchand.
var categories = []string{"Potions", "Grimoires", "Armes", "Armures", "Ressources", "Services"}

// Item décrit tout ce qui n'est pas un équipement : sa catégorie, sa couleur,
// son prix chez Mordecai et s'il s'utilise en plein combat. Les armes et les
// armures sont décrites à part, dans gear.
type Item struct {
	Category string
	Color    string
	Price    int
	InFight  bool
}

var items = map[string]Item{
	ItemHealthPotion:     {"Potions", Red, 8, true},
	ItemManaPotion:       {"Potions", Blue, 10, true},
	ItemPoisonPotion:     {"Potions", Green, 15, true},
	ItemFireballBook:     {"Grimoires", Orange, 60, true},
	ItemWolfFur:          {"Ressources", Brown, 10, false},
	ItemTrollSkin:        {"Ressources", Brown, 18, false},
	ItemBoarLeather:      {"Ressources", Brown, 8, false},
	ItemRavenFeather:     {"Ressources", Brown, 3, false},
	ItemInventoryUpgrade: {"Services", Gold, 75, false},
}

func itemCategory(item string) string {
	if info, ok := items[item]; ok {
		return info.Category
	}
	if isWeapon(item) {
		return "Armes"
	}
	return "Armures"
}

func itemColor(item string) string {
	if info, ok := items[item]; ok {
		return info.Color
	}
	return White
}

// ------------------------------------------------------------------ le sac

func inventoryCount(c *Character) int {
	total := 0
	for _, quantity := range c.Inventory {
		total += quantity
	}
	return total
}

func addInventory(c *Character, item string, quantity int) bool {
	if inventoryCount(c)+quantity > c.InventoryMax {
		fail("Votre sac est plein (%d / %d) ! Il va falloir faire du tri.", inventoryCount(c), c.InventoryMax)
		return false
	}
	c.Inventory[item] += quantity
	return true
}

func removeInventory(c *Character, item string, quantity int) bool {
	if c.Inventory[item] < quantity {
		return false
	}
	c.Inventory[item] -= quantity
	if c.Inventory[item] == 0 {
		delete(c.Inventory, item)
	}
	return true
}

// sortedItems range le sac par catégorie puis par ordre alphabétique.
func sortedItems(c *Character) []string {
	names := []string{}
	for _, category := range categories {
		group := []string{}
		for name := range c.Inventory {
			if itemCategory(name) == category {
				group = append(group, name)
			}
		}
		sort.Strings(group)
		names = append(names, group...)
	}
	return names
}

// showItems affiche le sac trié, avec un titre par catégorie. Le détail de
// chaque ligne est fourni par l'appelant (quantité, prix de revente…).
func showItems(c *Character, detail func(item string) string) []string {
	list := sortedItems(c)
	if len(list) == 0 {
		fmt.Println("\n" + DarkGray + Italic + "   Votre sac est vide. Même les mites sont parties." + Reset)
	}
	shown := ""
	for i, item := range list {
		if category := itemCategory(item); category != shown {
			section(category)
			shown = category
		}
		itemLine(i+1, item, itemColor(item), detail(item))
	}
	return list
}

func accessInventory(c *Character) {
	for {
		clearScreen()
		printArt(artBag, campColors...)
		banner("INVENTAIRE", Gold)
		showStatus(c)
		list := showItems(c, func(item string) string {
			return fmt.Sprintf("%sx%d%s", White, c.Inventory[item], Reset)
		})
		if len(list) > 0 {
			fmt.Println("\n" + Gray + Italic + "  Choisissez un objet pour l'utiliser ou l'équiper." + Reset)
		}
		fmt.Println()
		back("Retour")

		choice := readChoice(0, len(list))
		if choice == 0 {
			return
		}
		fmt.Println()
		useItem(c, nil, list[choice-1])
		pause()
	}
}

// ------------------------------------------------------------- usage d'objet

// useItem utilise un objet. m vaut nil hors combat ; en combat la Potion de
// poison est jetée sur l'ennemi au lieu d'être bue.
func useItem(c *Character, m *Monster, item string) {
	switch item {
	case ItemHealthPotion, ItemManaPotion:
		drinkPotion(c, item)
	case ItemPoisonPotion:
		if m != nil {
			throwPoison(c, m)
			return
		}
		drinkPoison(c)
		isDead(c)
	case ItemFireballBook:
		learnFireball(c)
	default:
		if isEquipment(item) {
			equipItem(c, item)
			return
		}
		info("%s ne se mange pas. Apportez-la plutôt à Borin, le forgeron.", item)
	}
}

// drinkPotion rend des PV (potion de vie) ou du mana (potion de mana).
func drinkPotion(c *Character, item string) {
	if !removeInventory(c, item, 1) {
		fail("Vous n'avez pas de %s.", item)
		return
	}
	printArt(artPotion, itemColor(item))
	if item == ItemHealthPotion {
		c.HP = min(c.HP+50, c.MaxHP)
		success("Glou glou… +50 PV. Goût fraise, pas mal.")
		showHP("", c.HP, c.MaxHP)
		return
	}
	c.Mana = min(c.Mana+30, c.MaxMana)
	success("Ça pétille jusque dans les oreilles : +30 mana.")
	showMana(c.Mana, c.MaxMana)
}

// drinkPoison : boire son propre poison coûte 30 PV en trois secondes.
func drinkPoison(c *Character) {
	if !removeInventory(c, ItemPoisonPotion, 1) {
		fail("Vous n'avez pas de Potion de poison.")
		return
	}
	printArt(artPotion, Green)
	warn("Vous buvez la Potion de poison. Audacieux. Idiot, mais audacieux.")
	if c.TestMode {
		info("Mode test : le poison n'a aucun effet sur vous.")
		return
	}
	for second := 1; second <= 3; second++ {
		wait(1000)
		c.HP = max(c.HP-10, 0)
		fmt.Printf("  %sSeconde %d : ça brûle ! -10 PV%s\n", Green, second, Reset)
	}
	showHP("", c.HP, c.MaxHP)
}

// throwPoison : jetée sur l'ennemi, la potion ronge 30 PV en trois secondes.
func throwPoison(c *Character, m *Monster) {
	if !removeInventory(c, ItemPoisonPotion, 1) {
		fail("Vous n'avez pas de Potion de poison.")
		return
	}
	fmt.Printf("  %sVous lancez la Potion de poison. Splash !%s\n", Green, Reset)
	for second := 1; second <= 3; second++ {
		wait(600)
		hurtMonster(m, 10, Green+"  (poison)"+Reset)
	}
}

func learnFireball(c *Character) {
	if c.knows(SpellFireball) {
		fail("Vous connaissez déjà Boule de Feu.")
		return
	}
	removeInventory(c, ItemFireballBook, 1)
	printArt(artBook, fireColors...)
	c.Skills = append(c.Skills, SpellFireball)
	success("Les pages prennent feu… et vous apprenez Boule de Feu ! Le livre, lui, est fichu.")
}

func upgradeInventorySlot(c *Character) bool {
	if c.InventoryUpgrades >= 3 {
		fail("Votre sac a déjà été agrandi 3 fois. Au-delà, c'est une armoire.")
		return false
	}
	c.InventoryUpgrades++
	c.InventoryMax += 10
	success("Votre sac peut maintenant contenir %d objets ! Votre dos proteste.", c.InventoryMax)
	return true
}
