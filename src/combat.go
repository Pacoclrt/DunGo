package main

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

const (
	Victory = "victoire"
	Defeat  = "défaite"
	Fled    = "fuite"
)

type Spell struct {
	Cost        int
	Description string
	Damage      int
	Art         string
	Colors      []string
}

var spells = map[string]Spell{
	SpellPunch:       {5, "8 dégâts", 8, artFist, []string{Orange}},
	SpellFireball:    {15, "18 dégâts + brûlure", 18, artFireball, fireColors},
	SpellSilverArrow: {8, "13 dégâts", 13, artArrow, stoneColors},
	SpellSecondWind:  {10, "rend 25 PV et soigne saignement et brûlure", 0, artHeal, []string{Green}},
	SpellStoneSkin:   {8, "dégâts reçus ÷ 2 pendant 2 coups", 0, artShield, stoneColors},
}

const (
	EffectBleed  = "saignement"
	EffectBurn   = "brûlure"
	EffectStun   = "étourdissement"
	EffectWeaken = "affaiblissement"
)

func trainingFight(c *Character) {
	clearScreen()
	printArt(artArena, campColors...)
	banner("L'ENTRAÎNEMENT", Gold)
	say("Sergent Grol", Gold, "Mon gobelin cogne mou, mais il cogne. Ici, pas de risque… et pas de butin !")
	pause()

	m := newMonster("training_goblin", "")
	m.Color = Gold
	m.Colors = campColors
	hpBefore := c.HP
	manaBefore := c.Mana
	fight(c, m, true)
	c.HP = hpBefore
	c.Mana = manaBefore
}

// fight enchaîne les tours : le héros joue, puis le monstre, puis les effets.
func fight(c *Character, m *Monster, training bool) string {
	clearEffects(c)
	for turn := 1; ; turn++ {
		clearScreen()
		showFight(c, m, turn)

		fled := characterTurn(c, m, training)
		if fled {
			return Fled
		}
		if m.HP > 0 {
			monsterTurn(m, c, turn)
		}
		if m.HP > 0 && c.HP > 0 {
			updateEffects(c, m)
		}

		if m.HP <= 0 {
			return winFight(c, m, training)
		}
		if c.HP <= 0 {
			return loseFight(c, training)
		}
		pause()
	}
}

func showFight(c *Character, m *Monster, turn int) {
	fmt.Printf("  %s── Tour %d %s%s\n", DarkGray, turn, strings.Repeat("─", 50), Reset)
	printArt(m.Art, m.Colors...)

	fmt.Print("  " + Bold + m.Color + m.Name + Reset)
	if m.IsBoss {
		fmt.Print(Red + "  ☠ BOSS" + Reset)
	}
	if m.Enraged {
		fmt.Print(Red + "  [EN RAGE]" + Reset)
	}
	if m.Burning > 0 {
		fmt.Printf("%s  [Brûlure %d]%s", Orange, m.Burning, Reset)
	}
	fmt.Printf("\n  %s♥%s %s %d / %d\n", Red, Reset, hpBar(m.HP, m.MaxHP, 40), m.HP, m.MaxHP)

	fmt.Println(DarkGray + "  " + strings.Repeat("─", 60) + Reset)
	fmt.Println("  " + Gold + Bold + c.Name + Reset + effectBadges(c))
	fmt.Printf("  %s♥%s %s %d / %d    %s♦%s %s %d / %d\n",
		Red, Reset, hpBar(c.HP, c.MaxHP, 20), c.HP, c.MaxHP,
		Blue, Reset, bar(c.Mana, c.MaxMana, 12, Blue), c.Mana, c.MaxMana)
}

// characterTurn fait jouer le héros. Renvoie true s'il quitte le combat.
func characterTurn(c *Character, m *Monster, training bool) bool {
	if c.Stunned > 0 {
		c.Stunned--
		warn("Vous êtes étourdi et passez votre tour.")
		return false
	}

	for {
		fmt.Println()
		optionHint(1, "Attaquer", fmt.Sprintf("%d dégâts", c.Attack))
		optionHint(2, "Sorts", fmt.Sprintf("%d mana", c.Mana))
		option(3, "Objets")
		if training {
			option(4, "Abandonner")
		} else {
			option(4, "Fuir")
		}

		switch readChoice(1, 4) {
		case 1:
			fmt.Println()
			damageMonster(c, m, c.Attack)
			return false
		case 2:
			if spellMenu(c, m) {
				return false
			}
		case 3:
			if combatInventory(c, m) {
				return false
			}
		case 4:
			return flee(c, m, training)
		}
	}
}

// flee renvoie true si le héros quitte le combat.
func flee(c *Character, m *Monster, training bool) bool {
	if training {
		info("Vous abandonnez l'entraînement.")
		pause()
		return true
	}
	if m.IsBoss {
		fail("On ne fuit pas un boss !")
		return false
	}
	if rand.IntN(100) < 50 {
		success("Vous courez jusqu'au camp. Avec beaucoup de dignité.")
		pause()
		return true
	}
	fail("%s vous rattrape. La fuite échoue !", m.Name)
	return false
}

// damageMonster applique un coup du héros : ÷ 2 s'il est affaibli, et 10 %
// de chance de coup critique (× 2).
func damageMonster(c *Character, m *Monster, damage int) {
	note := ""
	if c.Weakened > 0 {
		c.Weakened--
		damage = max(1, damage/2)
		note = note + " (affaibli)"
	}
	if rand.IntN(100) < 10 {
		damage = damage * 2
		note = note + " ★ CRITIQUE !"
	}
	hurtMonster(m, damage, note)
}

func hurtMonster(m *Monster, damage int, note string) {
	m.HP = max(m.HP-damage, 0)
	fmt.Printf("  %s» %s perd %d PV%s%s\n", Gold, m.Name, damage, note, Reset)
}

// hitPlayer : le monstre frappe le héros. Peau de Pierre divise les dégâts par 2.
func hitPlayer(m *Monster, c *Character, verb string, damage int) {
	note := ""
	if c.StoneSkin > 0 {
		c.StoneSkin--
		damage = max(1, damage/2)
		note = " (Peau de Pierre)"
	}
	c.HP = max(c.HP-damage, 0)
	fmt.Printf("  %s« %s %s : -%d PV%s%s\n", Red, m.Name, verb, damage, note, Reset)
}

// spellMenu renvoie true si un sort a été lancé.
func spellMenu(c *Character, m *Monster) bool {
	section("Sorts")
	for i, name := range c.Skills {
		spell := spells[name]
		option(i+1, fmt.Sprintf("%-16s %s%2d mana%s  %s%s%s", name, Blue, spell.Cost, Reset, Gray, spell.Description, Reset))
	}
	back("Retour")

	choice := readChoice(0, len(c.Skills))
	if choice == 0 {
		return false
	}
	name := c.Skills[choice-1]
	if c.Mana < spells[name].Cost {
		fail("Pas assez de mana (%d / %d).", c.Mana, spells[name].Cost)
		return false
	}
	castSpell(c, m, name)
	return true
}

func castSpell(c *Character, m *Monster, name string) {
	spell := spells[name]
	c.Mana -= spell.Cost
	fmt.Printf("\n  %s✦ %s !%s\n", Purple+Bold, name, Reset)
	printArt(spell.Art, spell.Colors...)

	if spell.Damage > 0 {
		damageMonster(c, m, spell.Damage)
	}

	switch name {
	case SpellFireball:
		m.Burning = 3
		warn("%s prend feu ! (-5 PV par tour pendant 3 tours)", m.Name)
	case SpellSecondWind:
		c.HP = min(c.HP+25, c.MaxHP)
		c.Bleeding = 0
		c.Burning = 0
		success("+25 PV, saignement et brûlure soignés.")
	case SpellStoneSkin:
		c.StoneSkin = 2
		success("Votre peau devient pierre : les 2 prochains coups font moitié moins mal.")
	}
}

// combatInventory renvoie true si un objet a été utilisé.
func combatInventory(c *Character, m *Monster) bool {
	usable := []string{}
	for _, item := range sortedItems(c) {
		if items[item].InFight {
			usable = append(usable, item)
		}
	}
	if len(usable) == 0 {
		fail("Aucun objet utilisable en combat.")
		return false
	}

	section("Objets")
	for i, item := range usable {
		option(i+1, fmt.Sprintf("%s x%d", item, c.Inventory[item]))
	}
	back("Retour")

	choice := readChoice(0, len(usable))
	if choice == 0 {
		return false
	}
	fmt.Println()
	useItem(c, m, usable[choice-1])
	return true
}

// monsterTurn : chaque boss a sa propre attaque, les autres monstres
// utilisent l'attaque de base. Tout suit un rythme de 3 tours.
func monsterTurn(m *Monster, c *Character, turn int) {
	switch m.Kind {
	case "goblin_king":
		goblinKingTurn(m, c, turn)
	case "lich":
		lichTurn(m, c, turn)
	case "dragon":
		dragonTurn(m, c, turn)
	default:
		basicTurn(m, c, turn)
	}
}

// basicTurn : une attaque normale, doublée tous les 3 tours.
func basicTurn(m *Monster, c *Character, turn int) {
	if turn%3 == 0 {
		fmt.Println(Orange + "  ATTAQUE PUISSANTE !" + Reset)
		hitPlayer(m, c, m.Hit, m.Attack*2)
		applyEffect(c, m.Effect)
	} else {
		hitPlayer(m, c, m.Hit, m.Attack)
	}
}

// goblinKingTurn : tous les 3 tours, un garde vient soigner Grukk.
func goblinKingTurn(m *Monster, c *Character, turn int) {
	if turn%3 == 0 {
		m.HP = min(m.HP+20, m.MaxHP)
		warn("Un garde accourt et soigne son roi : Grukk +20 PV !")
		hitPlayer(m, c, m.Hit, m.Attack*3/2)
	} else {
		hitPlayer(m, c, m.Hit, m.Attack)
	}
}

// lichTurn : tous les 3 tours, Mor'Vath vole du mana pour se soigner.
func lichTurn(m *Monster, c *Character, turn int) {
	if turn%3 == 0 {
		stolen := min(c.Mana, 15)
		c.Mana -= stolen
		m.HP = min(m.HP+stolen*2, m.MaxHP)
		fmt.Printf("  %sDRAIN D'ÂME ! Vous perdez %d mana, la Liche regagne %d PV.%s\n", Purple, stolen, stolen*2, Reset)
		hitPlayer(m, c, m.Hit, m.Attack)
		applyEffect(c, m.Effect)
	} else {
		hitPlayer(m, c, m.Hit, m.Attack)
	}
}

// dragonTurn : Ignarok s'enrage sous la moitié de ses PV (+30 % d'attaque)
// et crache du feu tous les 3 tours.
func dragonTurn(m *Monster, c *Character, turn int) {
	if !m.Enraged && m.HP <= m.MaxHP/2 {
		m.Enraged = true
		m.Attack = m.Attack * 13 / 10
		fmt.Println(Red + Bold + "  IGNAROK EST FURIEUX ! (attaque +30 %)" + Reset)
	}

	if turn%3 == 0 {
		printArt(artBreath, fireColors...)
		fmt.Println(Red + Bold + "  SOUFFLE INFERNAL !" + Reset)
		hitPlayer(m, c, "vous fait rôtir", m.Attack*5/2)
		applyEffect(c, m.Effect)
	} else {
		hitPlayer(m, c, m.Hit, m.Attack)
	}
}

func applyEffect(c *Character, effect string) {
	switch effect {
	case EffectBleed:
		c.Bleeding = 3
		warn("Vous saignez ! (-3 PV par tour, 3 tours)")
	case EffectBurn:
		c.Burning = 3
		warn("Vous brûlez ! (-5 PV par tour, 3 tours)")
	case EffectStun:
		c.Stunned = 1
		warn("Vous êtes étourdi ! (vous passez votre prochain tour)")
	case EffectWeaken:
		c.Weakened = 3
		warn("Vous êtes affaibli ! (vos 3 prochaines attaques ÷ 2)")
	}
}

func updateEffects(c *Character, m *Monster) {
	if c.Bleeding > 0 {
		c.Bleeding--
		c.HP = max(c.HP-3, 0)
		fmt.Println(Red + "  « Saignement : -3 PV" + Reset)
	}
	if c.Burning > 0 {
		c.Burning--
		c.HP = max(c.HP-5, 0)
		fmt.Println(Orange + "  « Brûlure : -5 PV" + Reset)
	}
	if m.Burning > 0 {
		m.Burning--
		hurtMonster(m, 5, " (brûlure)")
	}
}

func clearEffects(c *Character) {
	c.Bleeding = 0
	c.Burning = 0
	c.Stunned = 0
	c.Weakened = 0
	c.StoneSkin = 0
}

func effectBadges(c *Character) string {
	badges := ""
	if c.Bleeding > 0 {
		badges += fmt.Sprintf(" %s[Saignement %d]%s", Red, c.Bleeding, Reset)
	}
	if c.Burning > 0 {
		badges += fmt.Sprintf(" %s[Brûlure %d]%s", Orange, c.Burning, Reset)
	}
	if c.Stunned > 0 {
		badges += fmt.Sprintf(" %s[Étourdi]%s", Yellow, Reset)
	}
	if c.Weakened > 0 {
		badges += fmt.Sprintf(" %s[Affaibli %d]%s", Purple, c.Weakened, Reset)
	}
	if c.StoneSkin > 0 {
		badges += fmt.Sprintf(" %s[Peau de Pierre %d]%s", Silver, c.StoneSkin, Reset)
	}
	return badges
}

func winFight(c *Character, m *Monster, training bool) string {
	clearScreen()
	banner("VICTOIRE CONTRE "+strings.ToUpper(m.Name)+" !", Gold)

	// L'entraînement ne rapporte rien : sinon on y gagnerait des niveaux sans risque.
	if training {
		fmt.Println()
		say("Sergent Grol", Gold, "Propre ! Mais le butin, c'est au donjon.")
		pause()
		return Victory
	}

	c.Victories++
	c.Gold += m.Gold
	printArt(artChest, campColors...)
	section("Butin")
	fmt.Printf("  %s¤ +%d Y-Coins%s\n", Gold, m.Gold, Reset)
	if m.Drop != "" && rand.IntN(100) < m.DropChance {
		if addInventory(c, m.Drop, 1) {
			fmt.Println("  ■ " + m.Drop)
		}
	}
	gainXP(c, m.XP)
	updateQuest(c, m.Name)
	pause()
	return Victory
}

func loseFight(c *Character, training bool) string {
	if training {
		fmt.Println()
		say("Sergent Grol", Gold, "Battu par le gobelin d'entraînement… Relève-toi, personne n'a rien vu.")
		pause()
		return Defeat
	}
	isDead(c)
	lost := c.Gold / 5
	c.Gold -= lost
	fail("Vous perdez %d Y-Coins en chemin.", lost)
	pause()
	return Defeat
}
