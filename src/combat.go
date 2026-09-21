package main

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

// Issues possibles d'un combat.
const (
	Victory = "victoire"
	Defeat  = "défaite"
	Fled    = "fuite"
)

// Spell décrit un sort : son coût, ses dégâts et l'effet qu'il ajoute.
type Spell struct {
	Cost        int
	Description string
	Damage      int // 0 pour un sort de soutien
	Art         string
	Colors      []string
	Extra       func(*Character, *Monster) // effet en plus des dégâts
}

var spells = map[string]Spell{
	SpellPunch: {Cost: 5, Description: "8 dégâts", Damage: 8,
		Art: artFist, Colors: []string{Orange}},
	SpellFireball: {Cost: 15, Description: "18 dégâts + brûlure", Damage: 18,
		Art: artFireball, Colors: fireColors, Extra: castBurn},
	SpellSilverArrow: {Cost: 8, Description: "13 dégâts", Damage: 13,
		Art: artArrow, Colors: stoneColors},
	SpellSecondWind: {Cost: 10, Description: "rend 25 PV et soigne les saignements et brûlures",
		Art: artHeal, Colors: []string{Green}, Extra: castSecondWind},
	SpellStoneSkin: {Cost: 8, Description: "dégâts reçus ÷ 2 pendant 2 attaques",
		Art: artShield, Colors: stoneColors, Extra: castStoneSkin},
}

func castBurn(c *Character, m *Monster) {
	m.Burning = 3
	warn("%s prend feu ! (-5 PV par tour pendant 3 tours)", m.Name)
}

func castSecondWind(c *Character, m *Monster) {
	c.HP = min(c.HP+25, c.MaxHP)
	c.Bleeding, c.Burning = 0, 0
	success("+25 PV, et plus aucun bobo. Ça va mieux !")
}

func castStoneSkin(c *Character, m *Monster) {
	c.StoneSkin = 2
	success("Votre peau devient pierre : les 2 prochains coups font moitié moins mal.")
}

// ------------------------------------------------------------------ arène

func trainingFight(c *Character) {
	clearScreen()
	printArt(artArena, campColors...)
	banner("L'ENTRAÎNEMENT", Gold)
	say("Sergent Grol", Gold, "Mon gobelin cogne mou, mais il cogne. Ici, pas de risque… et pas de butin !")
	pause()
	m := newMonster("training_goblin", "")
	m.Color = Gold
	// Un combat pour de faux : on ressort de l'arène comme on y est entré.
	hp, mana := c.HP, c.Mana
	fight(c, m, true)
	c.HP, c.Mana = hp, mana
}

// ----------------------------------------------------------------- combat

// fight enchaîne les tours jusqu'à la victoire, la défaite ou la fuite.
// Le héros joue toujours en premier, puis le monstre riposte.
func fight(c *Character, m *Monster, training bool) string {
	clearEffects(c)
	defer clearEffects(c)

	for turn := 1; ; turn++ {
		clearScreen()
		showFight(c, m, turn)
		if characterTurn(c, m, training) {
			return Fled
		}
		if result := fightOver(c, m, training); result != "" {
			return result
		}
		monsterTurn(m, c, turn)
		if result := fightOver(c, m, training); result != "" {
			return result
		}
		updateEffects(c, m)
		if result := fightOver(c, m, training); result != "" {
			return result
		}
		pause()
	}
}

func fightOver(c *Character, m *Monster, training bool) string {
	switch {
	case m.HP <= 0:
		return winFight(c, m, training)
	case c.HP <= 0:
		return loseFight(c, training)
	}
	return ""
}

// showFight redessine l'écran de combat : monstre, jauges et effets.
func showFight(c *Character, m *Monster, turn int) {
	fmt.Printf("  %s── Tour %d %s%s\n", DarkGray, turn, strings.Repeat("─", 50), Reset)
	printArt(m.Art, m.Color)

	fmt.Printf("  %s%s%s", Bold+m.Color, m.Name, Reset)
	if m.IsBoss {
		fmt.Print(Red + Bold + "  ☠ BOSS" + Reset)
	}
	if m.Enraged {
		fmt.Print(Red + Bold + "  [EN RAGE]" + Reset)
	}
	if m.Burning > 0 {
		fmt.Printf(Orange+"  [Brûlure %d]"+Reset, m.Burning)
	}
	fmt.Printf("\n  %s♥%s %s %d / %d\n", Red, Reset, hpBar(m.HP, m.MaxHP, 40), m.HP, m.MaxHP)

	fmt.Println(DarkGray + "  " + strings.Repeat("─", 60) + Reset)
	fmt.Printf("  %s%s%s%s\n", Bold+Gold, c.Name, Reset, effectBadges(c))
	fmt.Printf("  %s♥%s %s %d / %d    %s♦%s %s %d / %d\n",
		Red, Reset, hpBar(c.HP, c.MaxHP, 20), c.HP, c.MaxHP,
		Blue, Reset, bar(c.Mana, c.MaxMana, 12, Blue), c.Mana, c.MaxMana)
}

// characterTurn renvoie true si le joueur a quitté le combat.
func characterTurn(c *Character, m *Monster, training bool) bool {
	if c.Stunned > 0 {
		c.Stunned--
		warn("Vous êtes étourdi. Vous voyez des étoiles et passez votre tour.")
		return false
	}

	leave := "Fuir"
	if training {
		leave = "Abandonner"
	}
	for {
		fmt.Println()
		optionHint(1, "Attaquer", fmt.Sprintf("%d dégâts", c.Attack))
		optionHint(2, "Sorts", fmt.Sprintf("%d mana", c.Mana))
		optionHint(3, "Objets", "potions, poison…")
		option(4, leave)

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

// flee renvoie true si le joueur quitte réellement le combat.
func flee(c *Character, m *Monster, training bool) bool {
	if training {
		info("Vous abandonnez. Le gobelin d'entraînement fait une danse de la victoire.")
		pause()
		return true
	}
	if m.IsBoss {
		fail("%s bloque la sortie. On ne fuit pas un boss !", m.Name)
		return false
	}
	dots("Vous tentez de fuir", Gray)
	if rand.IntN(100) < 50 {
		success("Vous courez jusqu'au camp. Avec beaucoup de dignité, évidemment.")
		pause()
		return true
	}
	fail("%s vous rattrape par le col. Raté !", m.Name)
	return false
}

// -------------------------------------------------------------------- dégâts

// damageMonster applique l'affaiblissement et les coups critiques du héros.
func damageMonster(c *Character, m *Monster, damage int) {
	note := ""
	if c.Weakened > 0 {
		damage = max(1, damage/2)
		c.Weakened--
		note += Purple + "  (affaibli : ÷ 2)" + Reset
	}
	if rand.IntN(100) < 10 {
		damage *= 2
		note += Gold + Bold + "  ★ CRITIQUE !" + Reset
	}
	hurtMonster(m, damage, note)
}

// hurtMonster retire des PV au monstre, sans bonus ni malus.
func hurtMonster(m *Monster, damage int, note string) {
	m.HP = max(m.HP-damage, 0)
	fmt.Printf("  %s» %s perd %d PV%s%s\n", Gold, m.Name, damage, Reset, note)
}

func damagePlayer(m *Monster, c *Character, damage int) {
	hitPlayer(m, c, m.Hit, damage)
}

// hitPlayer blesse le héros ; verb raconte le coup dans le journal de combat.
func hitPlayer(m *Monster, c *Character, verb string, damage int) {
	if c.TestMode {
		fmt.Printf("  %s« %s %s… et rebondit sur le mode test.%s\n", Gray, m.Name, verb, Reset)
		return
	}
	note := ""
	if c.StoneSkin > 0 {
		damage = max(1, damage/2)
		c.StoneSkin--
		note = Silver + "  (Peau de Pierre : ÷ 2)" + Reset
	}
	c.HP = max(c.HP-damage, 0)
	fmt.Printf("  %s« %s %s : -%d PV%s%s\n", Red, m.Name, verb, damage, Reset, note)
}

// --------------------------------------------------------------------- sorts

// spellMenu renvoie true si un sort a été lancé.
func spellMenu(c *Character, m *Monster) bool {
	section("Sorts")
	for i, name := range c.Skills {
		spell := spells[name]
		label := fmt.Sprintf("%-16s %s%2d mana%s · %s", name, Blue, spell.Cost, Reset+Gray, spell.Description)
		if !c.canCast(name) {
			label = DarkGray + fmt.Sprintf("%-16s %2d mana · %s", name, spell.Cost, spell.Description)
		}
		option(i+1, label+Reset)
	}
	back("Retour")

	choice := readChoice(0, len(c.Skills))
	if choice == 0 {
		return false
	}
	name := c.Skills[choice-1]
	if !c.canCast(name) {
		fail("Pas assez de mana pour %s (%d / %d). Il faudra taper.", name, c.Mana, spells[name].Cost)
		return false
	}

	spell := spells[name]
	if !c.TestMode {
		c.Mana -= spell.Cost
	}
	fmt.Printf("\n  %s✦ %s !%s\n", Purple+Bold, name, Reset)
	printArt(spell.Art, spell.Colors...)
	if spell.Damage > 0 {
		damageMonster(c, m, spell.Damage)
	}
	if spell.Extra != nil {
		spell.Extra(c, m)
	}
	return true
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
		fail("Vous fouillez vos poches : une miette et un bouton. Rien d'utile.")
		return false
	}

	section("Objets")
	for i, item := range usable {
		option(i+1, fmt.Sprintf("%s%s%s x%d", itemColor(item), item, Reset, c.Inventory[item]))
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

// ----------------------------------------------------------- tours du monstre

func monsterTurn(m *Monster, c *Character, turn int) {
	wait(400)
	if m.Pattern != nil {
		m.Pattern(m, c, turn)
		return
	}
	basicPattern(m, c, turn)
}

// basicPattern : attaque normale, doublée tous les 3 tours.
func basicPattern(m *Monster, c *Character, turn int) {
	if turn%3 != 0 {
		damagePlayer(m, c, m.Attack)
		return
	}
	fmt.Println(Orange + Bold + "  ATTAQUE PUISSANTE !" + Reset)
	damagePlayer(m, c, m.Attack*2)
	applyEffect(c, m.Effect)
}

// goblinKingPattern : Grukk appelle un garde qui le soigne tous les 3 tours.
func goblinKingPattern(king *Monster, c *Character, turn int) {
	switch turn % 3 {
	case 0:
		king.HP = min(king.HP+20, king.MaxHP)
		warn("Un garde accourt et soigne son roi : Grukk +20 PV !")
		damagePlayer(king, c, king.Attack*3/2)
	case 2:
		damagePlayer(king, c, king.Attack)
		warn("Grukk sort sa corne… Un garde arrive au prochain tour !")
	default:
		damagePlayer(king, c, king.Attack)
	}
}

// lichPattern : Mor'Vath vole le mana du héros pour se soigner.
func lichPattern(lich *Monster, c *Character, turn int) {
	switch turn % 3 {
	case 0:
		stolen := 0
		if !c.TestMode {
			stolen = min(c.Mana, 15)
			c.Mana -= stolen
		}
		lich.HP = min(lich.HP+stolen*2, lich.MaxHP)
		fmt.Printf("  %sDRAIN D'ÂME ! Vous perdez %d mana, la Liche regagne %d PV.%s\n", Purple+Bold, stolen, stolen*2, Reset)
		damagePlayer(lich, c, lich.Attack)
		applyEffect(c, lich.Effect)
	case 2:
		damagePlayer(lich, c, lich.Attack)
		warn("Les yeux de Mor'Vath s'illuminent… Drain d'âme au prochain tour !")
	default:
		damagePlayer(lich, c, lich.Attack)
	}
}

// dragonPattern : Ignarok s'enrage à mi-vie et souffle tous les 3 tours.
func dragonPattern(dragon *Monster, c *Character, turn int) {
	if !dragon.Enraged && dragon.HP <= dragon.MaxHP/2 {
		dragon.Enraged = true
		dragon.Attack += dragon.Attack * 3 / 10
		printArt(artDragon, bloodColors...)
		fmt.Println(Red + Bold + "  IGNAROK EST FURIEUX ! (attaque +30 %) On a dû le vexer." + Reset)
	}

	switch turn % 3 {
	case 0:
		printArt(artBreath, fireColors...)
		fmt.Println(Red + Bold + "  SOUFFLE INFERNAL !" + Reset)
		hitPlayer(dragon, c, "vous fait rôtir", dragon.Attack*5/2)
		applyEffect(c, dragon.Effect)
	case 2:
		damagePlayer(dragon, c, dragon.Attack)
		warn("Ignarok inspire très fort… Souffle Infernal au prochain tour !")
	default:
		damagePlayer(dragon, c, dragon.Attack)
	}
}

// ------------------------------------------------------------------- issues

func winFight(c *Character, m *Monster, training bool) string {
	wait(500)
	clearScreen()
	banner("VICTOIRE CONTRE "+strings.ToUpper(m.Name)+" !", Gold)

	// L'arène ne rapporte rien : sinon on y farmerait l'XP sans risque.
	if training {
		fmt.Println()
		say("Sergent Grol", Gold, "Propre ! Mais l'entraînement ne paie pas : le butin, c'est au donjon.")
		pause()
		return Victory
	}

	c.Victories++
	c.Gold += m.Gold
	c.GoldEarned += m.Gold
	printArt(artChest, campColors...)
	section("Butin")
	fmt.Printf("  %s¤ +%d Y-Coins%s\n", Gold+Bold, m.Gold, Reset)
	if m.Drop != "" && rand.IntN(100) < m.DropChance && addInventory(c, m.Drop, 1) {
		fmt.Printf("  %s■ %s%s\n", itemColor(m.Drop)+Bold, m.Drop, Reset)
	}
	gainXP(c, m.XP)
	updateQuest(c, m.Name)
	pause()
	return Victory
}

func loseFight(c *Character, training bool) string {
	wait(1200)
	// À l'entraînement on ne meurt pas vraiment : ni stèle, ni pénalité.
	if training {
		fmt.Println()
		say("Sergent Grol", Gold, "Battu par le gobelin d'entraînement… Relève-toi, personne n'a rien vu.")
		pause()
		return Defeat
	}
	isDead(c)
	lost := c.Gold / 5
	c.Gold -= lost
	fail("Quelqu'un a fait les poches de votre héros endormi : -%d Y-Coins.", lost)
	pause()
	return Defeat
}
