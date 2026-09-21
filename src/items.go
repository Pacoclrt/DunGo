package main

import (
	"fmt"
	"slices"
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

type Item struct {
	Price   int
	InFight bool // utilisable pendant un combat
}

var items = map[string]Item{
	ItemHealthPotion:     {8, true},
	ItemManaPotion:       {10, true},
	ItemPoisonPotion:     {15, true},
	ItemFireballBook:     {60, true},
	ItemWolfFur:          {10, false},
	ItemTrollSkin:        {18, false},
	ItemBoarLeather:      {8, false},
	ItemRavenFeather:     {3, false},
	ItemInventoryUpgrade: {75, false},
}

const (
	ItemAdventurerHat   = "Chapeau de l'aventurier"
	ItemAdventurerTunic = "Tunique de l'aventurier"
	ItemAdventurerBoots = "Bottes de l'aventurier"
	ItemTrollHelm       = "Heaume du Chasseur de Trolls"
	ItemBoarCuirass     = "Cuirasse du Sanglier Noir"
	ItemWerewolfBoots   = "Bottes du Loup-Garou"

	ItemAdventurerSword  = "Épée de l'aventurier"
	ItemAdventurerBow    = "Arc de l'aventurier"
	ItemAdventurerHammer = "Marteau de l'aventurier"
	ItemKnightBlade      = "Lame du Chevalier"
	ItemSylvanBow        = "Arc Sylvestre"
	ItemRunicHammer      = "Marteau Runique"
)

// Gear décrit un objet qui se porte : une armure donne des PV, une arme de
// l'attaque. Class est la classe pour laquelle l'arme est faite.
type Gear struct {
	Slot   string
	HP     int
	Attack int
	Class  string
}

var gear = map[string]Gear{
	ItemAdventurerHat:   {"Tête", 10, 0, ""},
	ItemAdventurerTunic: {"Torse", 25, 0, ""},
	ItemAdventurerBoots: {"Pieds", 15, 0, ""},
	ItemTrollHelm:       {"Tête", 25, 0, ""},
	ItemBoarCuirass:     {"Torse", 45, 0, ""},
	ItemWerewolfBoots:   {"Pieds", 30, 0, ""},

	ItemAdventurerSword:  {"Arme", 0, 3, "Humain"},
	ItemAdventurerBow:    {"Arme", 0, 3, "Elfe"},
	ItemAdventurerHammer: {"Arme", 0, 3, "Nain"},
	ItemKnightBlade:      {"Arme", 0, 7, "Humain"},
	ItemSylvanBow:        {"Arme", 0, 7, "Elfe"},
	ItemRunicHammer:      {"Arme", 0, 7, "Nain"},
}

var equipmentSlots = []string{"Tête", "Torse", "Pieds", "Arme"}

func isWeapon(item string) bool {
	return gear[item].Slot == "Arme"
}

// weaponBonus donne +2 attaque si l'arme est faite pour la classe du héros.
func weaponBonus(c *Character, item string) int {
	bonus := gear[item].Attack
	if isWeapon(item) && gear[item].Class == c.Class {
		bonus += 2
	}
	return bonus
}

func gearText(item string) string {
	if isWeapon(item) {
		return fmt.Sprintf("+%d attaque · pour %s", gear[item].Attack, gear[item].Class)
	}
	return fmt.Sprintf("+%d PV · %s", gear[item].HP, gear[item].Slot)
}

func inventoryCount(c *Character) int {
	total := 0
	for _, quantity := range c.Inventory {
		total += quantity
	}
	return total
}

func addInventory(c *Character, item string, quantity int) bool {
	if inventoryCount(c)+quantity > c.InventoryMax {
		fail("Votre sac est plein (%d / %d) !", inventoryCount(c), c.InventoryMax)
		return false
	}
	c.Inventory[item] += quantity
	return true
}

func removeInventory(c *Character, item string, quantity int) {
	c.Inventory[item] -= quantity
	if c.Inventory[item] <= 0 {
		delete(c.Inventory, item)
	}
}

// sortedItems renvoie les objets du sac par ordre alphabétique
// (une map Go n'a pas d'ordre).
func sortedItems(c *Character) []string {
	names := []string{}
	for name := range c.Inventory {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func accessInventory(c *Character) {
	for {
		clearScreen()
		printArt(artBag, campColors...)
		banner("INVENTAIRE", Gold)
		showStatus(c)
		fmt.Println()

		list := sortedItems(c)
		if len(list) == 0 {
			info("Votre sac est vide. Même les mites sont parties.")
		}
		for i, item := range list {
			itemLine(i+1, item, fmt.Sprintf("x%d", c.Inventory[item]))
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

// useItem utilise un objet du sac. m vaut nil hors combat : la Potion de
// poison est alors bue au lieu d'être lancée sur l'ennemi.
func useItem(c *Character, m *Monster, item string) {
	if _, isGear := gear[item]; isGear {
		equipItem(c, item)
		return
	}

	switch item {
	case ItemHealthPotion:
		c.HP = min(c.HP+50, c.MaxHP)
		success("Glou glou… +50 PV (%d / %d).", c.HP, c.MaxHP)

	case ItemManaPotion:
		c.Mana = min(c.Mana+30, c.MaxMana)
		success("Ça pétille : +30 mana (%d / %d).", c.Mana, c.MaxMana)

	case ItemPoisonPotion:
		if m != nil {
			fmt.Println(Green + "  Vous lancez la Potion de poison. Splash !" + Reset)
			for second := 1; second <= 3; second++ {
				wait(1000)
				hurtMonster(m, 10, " (poison)")
			}
		} else {
			warn("Vous buvez la Potion de poison. Audacieux. Idiot, mais audacieux.")
			for second := 1; second <= 3; second++ {
				wait(1000)
				c.HP = max(c.HP-10, 0)
				fmt.Printf("  %sLe poison vous brûle : -10 PV (%d / %d)%s\n", Green, c.HP, c.MaxHP, Reset)
			}
			isDead(c)
		}

	case ItemFireballBook:
		if slices.Contains(c.Skills, SpellFireball) {
			fail("Vous connaissez déjà Boule de Feu.")
			return
		}
		c.Skills = append(c.Skills, SpellFireball)
		printArt(artBook, fireColors...)
		success("Vous apprenez le sort Boule de Feu ! Le livre, lui, part en fumée.")

	default:
		info("%s ne s'utilise pas : apportez-la au forgeron.", item)
		return
	}

	removeInventory(c, item, 1)
}

// equipItem porte un objet. L'objet déjà porté au même endroit retourne dans le sac.
func equipItem(c *Character, item string) {
	slot := gear[item].Slot
	removeInventory(c, item, 1)

	old := c.Equipment[slot]
	if old != "" {
		c.MaxHP -= gear[old].HP
		c.Attack -= weaponBonus(c, old)
		addInventory(c, old, 1)
		info("%s retourne dans le sac.", old)
	}

	c.Equipment[slot] = item
	c.MaxHP += gear[item].HP
	c.HP = min(c.HP, c.MaxHP)
	c.Attack += weaponBonus(c, item)

	if isWeapon(item) {
		printArt(artSword, stoneColors...)
		success("Vous équipez %s : +%d attaque (%d en tout).", item, weaponBonus(c, item), c.Attack)
		if gear[item].Class != c.Class {
			warn("Cette arme est faite pour la classe %s : pas de bonus.", gear[item].Class)
		}
	} else {
		printArt(artShield, stoneColors...)
		success("Vous enfilez %s : +%d PV max (%d en tout).", item, gear[item].HP, c.MaxHP)
	}
}
