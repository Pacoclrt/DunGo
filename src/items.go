// ════════════════════════════════════════════════════════════════════════
//   items.go · [INVENTAIRE]
//   Les objets (potions, livre, ressources), les équipements (armes et
//   armures), le sac du héros, et ce qui se passe quand on utilise ou
//   qu'on équipe un objet.
// ════════════════════════════════════════════════════════════════════════

package main

import (
	"fmt"
	"slices"
	"sort"
)

// ════════════════════════════════════════════════════════════════════════
//   LES OBJETS
// ════════════════════════════════════════════════════════════════════════

const ( // [INVENTAIRE] sert à nommer les objets qu'on peut acheter ou ramasser (ces noms sont aussi ceux affichés dans le jeu)
	ItemHealthPotion     = "Potion de vie"
	ItemManaPotion       = "Potion de mana"
	ItemPoisonPotion     = "Potion de poison"
	ItemFireballBook     = "Livre de Sort : Boule de Feu"
	ItemInventoryUpgrade = "Augmentation d'inventaire"

	// Les ressources : elles ne servent qu'au forgeron
	ItemWolfFur      = "Fourrure de Loup"
	ItemTrollSkin    = "Peau de Troll"
	ItemBoarLeather  = "Cuir de Sanglier"
	ItemRavenFeather = "Plume de Corbeau"
)

type Item struct { // [INVENTAIRE] sert à décrire un objet : son prix chez le marchand et s'il peut servir en combat
	Price         int  // en Y-Coins
	UsableInFight bool // true = il apparaît dans le menu « Objets » du combat
}

var items = map[string]Item{ // [INVENTAIRE] sert à ranger les objets par leur nom : items["Potion de vie"].Price donne le prix
	ItemHealthPotion:     {Price: 8, UsableInFight: true},
	ItemManaPotion:       {Price: 10, UsableInFight: true},
	ItemPoisonPotion:     {Price: 15, UsableInFight: true},
	ItemFireballBook:     {Price: 60, UsableInFight: true},
	ItemWolfFur:          {Price: 10, UsableInFight: false},
	ItemTrollSkin:        {Price: 18, UsableInFight: false},
	ItemBoarLeather:      {Price: 8, UsableInFight: false},
	ItemRavenFeather:     {Price: 3, UsableInFight: false},
	ItemInventoryUpgrade: {Price: 75, UsableInFight: false},
}

// ════════════════════════════════════════════════════════════════════════
//   LES ÉQUIPEMENTS (armures et armes)
// ════════════════════════════════════════════════════════════════════════

const ( // [INVENTAIRE] sert à nommer les 12 équipements fabriqués par le forgeron
	// Les armures (+PV max)
	ItemAdventurerHat   = "Chapeau de l'aventurier"
	ItemAdventurerTunic = "Tunique de l'aventurier"
	ItemAdventurerBoots = "Bottes de l'aventurier"
	ItemTrollHelm       = "Heaume du Chasseur de Trolls"
	ItemBoarCuirass     = "Cuirasse du Sanglier Noir"
	ItemWerewolfBoots   = "Bottes du Loup-Garou"

	// Les armes (+attaque)
	ItemAdventurerSword  = "Épée de l'aventurier"
	ItemAdventurerBow    = "Arc de l'aventurier"
	ItemAdventurerHammer = "Marteau de l'aventurier"
	ItemKnightBlade      = "Lame du Chevalier"
	ItemSylvanBow        = "Arc Sylvestre"
	ItemRunicHammer      = "Marteau Runique"
)

type Gear struct { // [INVENTAIRE] sert à décrire un équipement : où il se porte et ce qu'il rapporte (PV pour une armure, attaque pour une arme)
	Slot        string // "Tête", "Torse", "Pieds" ou "Arme"
	BonusHP     int    // PV max en plus (armures)
	BonusAttack int    // attaque en plus (armes)
	ForClass    string // la classe pour qui l'arme est faite : elle lui donne +2 attaque en plus
}

var gear = map[string]Gear{ // [INVENTAIRE] sert à ranger les 12 équipements par leur nom : gear["Arc Sylvestre"] donne l'arme
	ItemAdventurerHat:   {Slot: "Tête", BonusHP: 10},
	ItemAdventurerTunic: {Slot: "Torse", BonusHP: 25},
	ItemAdventurerBoots: {Slot: "Pieds", BonusHP: 15},
	ItemTrollHelm:       {Slot: "Tête", BonusHP: 25},
	ItemBoarCuirass:     {Slot: "Torse", BonusHP: 45},
	ItemWerewolfBoots:   {Slot: "Pieds", BonusHP: 30},

	ItemAdventurerSword:  {Slot: "Arme", BonusAttack: 3, ForClass: "Humain"},
	ItemAdventurerBow:    {Slot: "Arme", BonusAttack: 3, ForClass: "Elfe"},
	ItemAdventurerHammer: {Slot: "Arme", BonusAttack: 3, ForClass: "Nain"},
	ItemKnightBlade:      {Slot: "Arme", BonusAttack: 7, ForClass: "Humain"},
	ItemSylvanBow:        {Slot: "Arme", BonusAttack: 7, ForClass: "Elfe"},
	ItemRunicHammer:      {Slot: "Arme", BonusAttack: 7, ForClass: "Nain"},
}

var equipmentSlots = []string{"Tête", "Torse", "Pieds", "Arme"} // [INVENTAIRE] sert à lister les 4 emplacements d'équipement, dans l'ordre d'affichage

func isWeapon(item string) bool { // [INVENTAIRE] sert à savoir si un objet est une arme ; renvoie true si oui
	return gear[item].Slot == "Arme"
}

func attackBonus(character *Character, item string) int { // [INVENTAIRE] sert à calculer l'attaque qu'un équipement donne au héros (+2 si l'arme est faite pour sa classe) et la renvoie
	bonus := gear[item].BonusAttack
	if isWeapon(item) && gear[item].ForClass == character.Class {
		bonus += 2
	}
	return bonus
}

func gearBonusText(item string) string { // [INVENTAIRE] sert à décrire le bonus d'un équipement en texte (« +7 attaque · pour Elfe » ou « +25 PV · Tête ») et le renvoie
	if isWeapon(item) {
		return fmt.Sprintf("+%d attaque · pour %s", gear[item].BonusAttack, gear[item].ForClass)
	}
	return fmt.Sprintf("+%d PV · %s", gear[item].BonusHP, gear[item].Slot)
}

// ════════════════════════════════════════════════════════════════════════
//   LE SAC
//   character.Inventory est une map « nom de l'objet → quantité ».
// ════════════════════════════════════════════════════════════════════════

func inventoryCount(character *Character) int { // [INVENTAIRE] sert à compter le nombre total d'objets dans le sac et le renvoie
	total := 0
	for _, quantity := range character.Inventory {
		total += quantity
	}
	return total
}

func addToInventory(character *Character, item string, quantity int) bool { // [INVENTAIRE] sert à ranger un objet dans le sac ; renvoie false (et prévient) si le sac est plein
	if inventoryCount(character)+quantity > character.InventoryMax {
		fail("Votre sac est plein (%d / %d) !", inventoryCount(character), character.InventoryMax)
		return false
	}
	character.Inventory[item] += quantity
	return true
}

func removeFromInventory(character *Character, item string, quantity int) { // [INVENTAIRE] sert à retirer un objet du sac (la ligne disparaît quand il n'en reste plus)
	character.Inventory[item] -= quantity
	if character.Inventory[item] <= 0 {
		delete(character.Inventory, item)
	}
}

func sortedInventory(character *Character) []string { // [INVENTAIRE] sert à renvoyer les noms des objets du sac par ordre alphabétique
	// Une map Go n'a pas d'ordre : sans ce tri, les objets s'afficheraient
	// dans un ordre différent à chaque fois.
	names := []string{}
	for name := range character.Inventory {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// ════════════════════════════════════════════════════════════════════════
//   L'ÉCRAN D'INVENTAIRE
// ════════════════════════════════════════════════════════════════════════

func inventoryMenu(character *Character) { // [INVENTAIRE] sert à afficher le sac et à utiliser l'objet choisi (boire, apprendre, équiper), jusqu'au Retour
	for {
		clearScreen()
		printArt(artBag, campColors...)
		banner("INVENTAIRE", Gold)
		showCharacterBar(character)
		fmt.Println()

		list := sortedInventory(character)
		if len(list) == 0 {
			info("Votre sac est vide. Même les mites sont parties.")
		}
		for index, item := range list {
			itemLine(index+1, item, fmt.Sprintf("x%d", character.Inventory[item]))
		}
		fmt.Println()
		backOption("Retour")

		choice := readChoice(0, len(list))
		if choice == 0 {
			return
		}
		fmt.Println()
		useItem(character, nil, list[choice-1]) // nil : pas de monstre, on est hors combat
		pause()
	}
}

func useItem(character *Character, monster *Monster, item string) bool { // [INVENTAIRE] sert à utiliser un objet (potion, livre, équipement) ; renvoie true s'il a servi (il est alors retiré du sac)
	// monster vaut nil hors combat : la Potion de poison est alors bue au
	// lieu d'être lancée sur l'ennemi.

	// Un équipement : on le porte
	if _, isGear := gear[item]; isGear {
		equipItem(character, item)
		return true
	}

	switch item {
	case ItemHealthPotion:
		if character.HP == character.MaxHP {
			fail("Vous avez déjà tous vos PV : gardez la potion pour plus tard !")
			return false
		}
		character.HP = min(character.HP+50, character.MaxHP)
		success("Glou glou… +50 PV (%d / %d).", character.HP, character.MaxHP)

	case ItemManaPotion:
		if character.Mana == character.MaxMana {
			fail("Votre mana est déjà au maximum : gardez la potion pour plus tard !")
			return false
		}
		character.Mana = min(character.Mana+30, character.MaxMana)
		success("Ça pétille : +30 mana (%d / %d).", character.Mana, character.MaxMana)

	case ItemPoisonPotion:
		if monster != nil {
			// En combat : lancée sur le monstre, 10 dégâts par seconde pendant 3 secondes
			fmt.Println(Green + "  Vous lancez la Potion de poison. Splash !" + Reset)
			for second := 1; second <= 3; second++ {
				wait(1000)
				removeMonsterHP(monster, 10, " (poison)")
			}
		} else {
			// Hors combat : bue… par le héros !
			warn("Vous buvez la Potion de poison. Audacieux. Idiot, mais audacieux.")
			for second := 1; second <= 3; second++ {
				wait(1000)
				character.HP = max(character.HP-10, 0)
				fmt.Printf("  %sLe poison vous brûle : -10 PV (%d / %d)%s\n", Green, character.HP, character.MaxHP, Reset)
			}
			if character.HP <= 0 {
				characterDies(character)
			}
		}

	case ItemFireballBook:
		if slices.Contains(character.Spells, SpellFireball) {
			fail("Vous connaissez déjà Boule de Feu.")
			return false
		}
		character.Spells = append(character.Spells, SpellFireball)
		printArt(artBook, fireColors...)
		success("Vous apprenez le sort Boule de Feu ! Le livre, lui, part en fumée.")

	default:
		// Les ressources (fourrure, peau…) ne s'utilisent pas.
		info("%s ne s'utilise pas : apportez cette ressource au forgeron.", item)
		return false
	}

	removeFromInventory(character, item, 1)
	return true
}

func equipItem(character *Character, item string) { // [INVENTAIRE] sert à porter un équipement : l'ancien objet au même emplacement retourne dans le sac, et les bonus sont mis à jour
	slot := gear[item].Slot
	removeFromInventory(character, item, 1)

	// 1. On enlève l'objet déjà porté à cet emplacement (et ses bonus)
	oldItem := character.Equipment[slot]
	if oldItem != "" {
		character.MaxHP -= gear[oldItem].BonusHP
		character.Attack -= attackBonus(character, oldItem)
		addToInventory(character, oldItem, 1) // il y a forcément la place : on vient de sortir le nouvel objet du sac
		info("%s retourne dans le sac.", oldItem)
	}

	// 2. On porte le nouvel objet (et on ajoute ses bonus)
	character.Equipment[slot] = item
	character.MaxHP += gear[item].BonusHP
	character.HP = min(character.HP, character.MaxHP) // si la nouvelle armure donne moins de PV max
	character.Attack += attackBonus(character, item)

	// 3. Le message
	if isWeapon(item) {
		printArt(artSword, stoneColors...)
		success("Vous équipez %s : +%d attaque (%d en tout).", item, attackBonus(character, item), character.Attack)
		if gear[item].ForClass != character.Class {
			warn("Cette arme est faite pour la classe %s : pas de bonus.", gear[item].ForClass)
		}
	} else {
		printArt(artShield, stoneColors...)
		success("Vous enfilez %s : +%d PV max (%d en tout).", item, gear[item].BonusHP, character.MaxHP)
	}
}
