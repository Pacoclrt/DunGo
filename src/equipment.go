package main

import "fmt"

// Equipment retient ce que le héros porte sur lui.
type Equipment struct {
	Head   string
	Torso  string
	Feet   string
	Weapon string
}

// Emplacements, dans l'ordre d'affichage de la fiche.
var equipmentSlots = []string{"Tête", "Torse", "Pieds", "Arme"}

// slot renvoie l'emplacement correspondant, pour le lire ou le remplacer.
func (c *Character) slot(name string) *string {
	switch name {
	case "Tête":
		return &c.Equipment.Head
	case "Torse":
		return &c.Equipment.Torso
	case "Pieds":
		return &c.Equipment.Feet
	}
	return &c.Equipment.Weapon
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

// Gear décrit un objet qui se porte : une armure donne des PV, une arme
// de l'attaque. Class indique la lignée à laquelle l'arme est destinée.
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

func isEquipment(item string) bool {
	_, ok := gear[item]
	return ok
}

func isWeapon(item string) bool {
	return gear[item].Slot == "Arme"
}

// weaponBonus ajoute +2 attaque si l'arme est faite pour la lignée du héros.
func weaponBonus(c *Character, item string) int {
	if isWeapon(item) && gear[item].Class == c.Class {
		return gear[item].Attack + 2
	}
	return gear[item].Attack
}

func itemBonusText(item string) string {
	if isWeapon(item) {
		return fmt.Sprintf("+%d attaque", gear[item].Attack)
	}
	return fmt.Sprintf("+%d PV", gear[item].HP)
}

func itemSlotText(item string) string {
	if isWeapon(item) {
		return "Arme · " + gear[item].Class
	}
	return gear[item].Slot
}

// equipItem porte un objet ; celui qu'il remplace retourne dans le sac.
func equipItem(c *Character, item string) {
	removeInventory(c, item, 1)
	worn := c.slot(gear[item].Slot)
	if old := *worn; old != "" {
		c.MaxHP -= gear[old].HP
		c.Attack -= weaponBonus(c, old)
		addInventory(c, old, 1) // la place libérée juste avant suffit
		info("%s retourne dans votre sac.", old)
	}
	*worn = item
	c.MaxHP += gear[item].HP
	c.HP = min(c.HP, c.MaxHP)
	c.Attack += weaponBonus(c, item)

	if !isWeapon(item) {
		printArt(artShield, steelColors...)
		success("Vous équipez %s (%s, +%d PV max).", item, gear[item].Slot, gear[item].HP)
		showHP("", c.HP, c.MaxHP)
		return
	}
	printArt(artSword, steelColors...)
	success("Vous équipez %s : +%d attaque (attaque totale : %d).", item, weaponBonus(c, item), c.Attack)
	if gear[item].Class == c.Class {
		info("Arme de votre lignée : +2 attaque bonus !")
	} else {
		warn("Cette arme est faite pour un %s : vous ne profitez pas du bonus de lignée.", gear[item].Class)
	}
}
