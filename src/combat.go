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
		Art: artArrow, Colors: steelColors},
	SpellSecondWind: {Cost: 10, Description: "rend 25 PV et soigne les saignements et brûlures",
		Art: artHeal, Colors: []string{Green}, Extra: castSecondWind},
	SpellStoneSkin: {Cost: 8, Description: "dégâts reçus ÷ 2 pendant 2 attaques",
		Art: artShield, Colors: steelColors, Extra: castStoneSkin},
}

func castBurn(c *Character, m *Monster) {
	m.Burning = 3
	warn("%s prend feu ! (-5 PV par tour pendant 3 tours)", m.Name)
}

func castSecondWind(c *Character, m *Monster) {
	c.HP = min(c.HP+25, c.MaxHP)
	c.Bleeding, c.Burning = 0, 0
	success("Un souffle de courage : +25 PV, saignements et brûlures soignés.")
	showHP("", c.HP, c.MaxHP)
}

func castStoneSkin(c *Character, m *Monster) {
	c.StoneSkin = 2
	success("Votre peau devient granit : dégâts divisés par 2 pendant 2 attaques.")
}

// ------------------------------------------------------------------ arène

func trainingFight(c *Character) {
	clearScreen()
	printArt(artArena, bloodColors...)
	banner("L'ARÈNE DU SERGENT GROL", Red)
	say("Sergent Grol", Red, "Alors, la recrue ! Mon gobelin cogne mou, mais il cogne.")
	say("Sergent Grol", Red, "Ici, pas de pénalité si tu tombes. Mais pas de butin non plus : on apprend !")
	pause()
	fight(c, newMonster("training_goblin", ""), true)
}

// ----------------------------------------------------------------- combat

// fight enchaîne les tours jusqu'à la victoire, la défaite ou la fuite.
func fight(c *Character, m *Monster, training bool) string {
	clearEffects(c)
	defer clearEffects(c)
	heroFirst := c.Initiative >= m.Initiative

	for turn := 1; ; turn++ {
		clearScreen()
		showFight(c, m, turn)
		if turn == 1 {
			if heroFirst {
				info("Initiative %d contre %d : vous attaquez en premier !", c.Initiative, m.Initiative)
			} else {
				warn("Initiative %d contre %d : %s est plus rapide !", c.Initiative, m.Initiative, m.Name)
			}
		}

		for _, heroTurn := range []bool{heroFirst, !heroFirst} {
			if !heroTurn {
				monsterTurn(m, c, turn)
			} else if characterTurn(c, m, training) {
				return Fled
			}
			if result := fightOver(c, m, training); result != "" {
				return result
			}
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
	fmt.Printf("%s  ⚔══════════════════  TOUR %d  ══════════════════⚔%s\n", Gold+Bold, turn, Reset)
	printArt(m.Art, m.Color)

	fmt.Printf("  %s%s%s", Bold+m.Color, m.Name, Reset)
	if m.IsBoss {
		fmt.Print(Red + Bold + "  ☠ BOSS" + Reset)
	}
	if m.Enraged {
		fmt.Print(Red + Bold + "  [ EN RAGE ]" + Reset)
	}
	if m.Burning > 0 {
		fmt.Printf(Orange+"  [Brûlure %d]"+Reset, m.Burning)
	}
	fmt.Printf("\n  %s♥%s %s %d / %d\n", Red, Reset, hpBar(m.HP, m.MaxHP, 40), m.HP, m.MaxHP)

	fmt.Println(DarkGray + "  " + strings.Repeat("─", 60) + Reset)
	fmt.Printf("  %s%s%s %s%s niv.%d%s%s\n", Bold+Gold, c.Name, Reset, classOf(c.Class).Color, c.Class, c.Level, Reset, effectBadges(c))
	fmt.Printf("  %s♥%s %s %d / %d    %s♦%s %s %d / %d\n",
		Red, Reset, hpBar(c.HP, c.MaxHP, 20), c.HP, c.MaxHP,
		Blue, Reset, bar(c.Mana, c.MaxMana, 12, Blue), c.Mana, c.MaxMana)
}

// characterTurn renvoie true si le joueur a quitté le combat.
func characterTurn(c *Character, m *Monster, training bool) bool {
	if c.Stunned > 0 {
		c.Stunned--
		warn("Vous êtes étourdi et ne pouvez pas agir ce tour-ci !")
		return false
	}

	leave := "Fuir"
	if training {
		leave = "Abandonner"
	}
	for {
		section("À vous de jouer !", Gold)
		option(1, fmt.Sprintf("Attaquer     %s(Attaque basique · %d dégâts)%s", Gray, c.Attack, Reset), Orange)
		option(2, fmt.Sprintf("Sorts        %s(%d mana disponible)%s", Gray, c.Mana, Reset), Purple)
		option(3, "Inventaire", Sky)
		option(4, leave, Gray)

		switch readChoice(1, 4) {
		case 1:
			fmt.Printf("\n  %s%s utilise Attaque basique !%s\n", Orange+Bold, c.Name, Reset)
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
		info("Vous abandonnez l'entraînement. Le gobelin a l'air soulagé.")
		return true
	}
	if m.IsBoss {
		fail("%s vous barre la route. Impossible de fuir un boss !", m.Name)
		return false
	}
	dots("Vous tentez de fuir", Gray)
	if rand.IntN(100) < 50 {
		success("Vous prenez vos jambes à votre cou et remontez au camp !")
		pause()
		return true
	}
	fail("%s vous rattrape ! La fuite échoue.", m.Name)
	return false
}

// -------------------------------------------------------------------- dégâts

// damageMonster applique l'affaiblissement et les coups critiques du héros.
func damageMonster(c *Character, m *Monster, damage int) {
	if c.Weakened > 0 {
		damage = max(1, damage/2)
		c.Weakened--
		fmt.Println(Purple + "  Affaibli : vos dégâts sont divisés par 2." + Reset)
	}
	if rand.IntN(100) < 10 {
		damage *= 2
		fmt.Println(Gold + Bold + "  ★ COUP CRITIQUE ! ★" + Reset)
	}
	hurtMonster(c, m, damage)
}

// hurtMonster retire des PV au monstre, sans bonus ni malus.
func hurtMonster(c *Character, m *Monster, damage int) {
	m.HP = max(m.HP-damage, 0)
	fmt.Printf("  %s%s inflige %d dégâts à %s%s\n", Yellow, c.Name, damage, m.Name, Reset)
	showHP(m.Name, m.HP, m.MaxHP)
}

func damagePlayer(m *Monster, c *Character, damage int) {
	if c.TestMode {
		fmt.Printf("  %s%s vous frappe… mais le mode test vous rend invincible !%s\n", Cyan, m.Name, Reset)
		return
	}
	if c.StoneSkin > 0 {
		damage = max(1, damage/2)
		c.StoneSkin--
		fmt.Println(Silver + "  ■ Peau de Pierre absorbe la moitié du coup !" + Reset)
	}
	c.HP = max(c.HP-damage, 0)
	fmt.Printf("  %s%s inflige à %s %d de dégâts%s\n", Red, m.Name, c.Name, damage, Reset)
	showHP(c.Name, c.HP, c.MaxHP)
}

// --------------------------------------------------------------------- sorts

// spellMenu renvoie true si un sort a été lancé.
func spellMenu(c *Character, m *Monster) bool {
	section("Grimoire", Purple)
	for i, name := range c.Skills {
		spell := spells[name]
		color := Purple
		if !c.canCast(name) {
			color = DarkGray
		}
		option(i+1, fmt.Sprintf("%-16s %s%2d mana%s · %s", name, Blue, spell.Cost, Reset+Gray, spell.Description), color)
	}
	option(0, "Retour", Gray)

	choice := readChoice(0, len(c.Skills))
	if choice == 0 {
		return false
	}
	name := c.Skills[choice-1]
	if !c.canCast(name) {
		fail("Mana insuffisant pour %s (%d / %d).", name, c.Mana, spells[name].Cost)
		return false
	}

	spell := spells[name]
	if !c.TestMode {
		c.Mana -= spell.Cost
	}
	fmt.Printf("\n  %s%s lance %s !%s\n", Purple+Bold, c.Name, name, Reset)
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
		fail("Aucun objet utilisable en combat.")
		return false
	}

	section("Objets utilisables", Sky)
	for i, item := range usable {
		option(i+1, fmt.Sprintf("%s x%d", item, c.Inventory[item]), itemColor(item))
	}
	option(0, "Retour", Gray)

	choice := readChoice(0, len(usable))
	if choice == 0 {
		return false
	}
	fmt.Printf("\n  %sVous utilisez %s%s\n", Sky+Bold, usable[choice-1], Reset)
	useItem(c, m, usable[choice-1])
	return true
}

// ----------------------------------------------------------- tours du monstre

func monsterTurn(m *Monster, c *Character, turn int) {
	section("Au tour de "+m.Name, m.Color)
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
	fmt.Printf("  %s%s prend son élan pour une ATTAQUE PUISSANTE !%s\n", Orange+Bold, m.Name, Reset)
	damagePlayer(m, c, m.Attack*2)
	applyEffect(c, m.Effect)
}

// goblinKingPattern : Grukk appelle un garde qui le soigne tous les 3 tours.
func goblinKingPattern(king *Monster, c *Character, turn int) {
	switch turn % 3 {
	case 0:
		fmt.Printf("  %sGrukk souffle dans sa corne : un garde gobelin accourt !%s\n", Gold+Bold, Reset)
		king.HP = min(king.HP+20, king.MaxHP)
		fmt.Printf("  %sLe garde soigne son roi : +20 PV%s\n", Green, Reset)
		showHP(king.Name, king.HP, king.MaxHP)
		damagePlayer(king, c, king.Attack*3/2)
	case 2:
		warn("Grukk porte sa corne à la bouche… (au prochain tour, il appellera un garde !)")
		damagePlayer(king, c, king.Attack)
	default:
		damagePlayer(king, c, king.Attack)
	}
}

// lichPattern : Mor'Vath vole le mana du héros pour se soigner.
func lichPattern(lich *Monster, c *Character, turn int) {
	switch turn % 3 {
	case 0:
		fmt.Printf("  %sMor'Vath lance DRAIN D'ÂME !%s\n", Purple+Bold, Reset)
		stolen := 0
		if !c.TestMode {
			stolen = min(c.Mana, 15)
			c.Mana -= stolen
		}
		lich.HP = min(lich.HP+stolen*2, lich.MaxHP)
		fmt.Printf("  %sVous perdez %d mana, la Liche récupère %d PV.%s\n", Purple, stolen, stolen*2, Reset)
		damagePlayer(lich, c, lich.Attack)
		applyEffect(c, lich.Effect)
	case 2:
		warn("Les yeux de Mor'Vath s'illuminent… (au prochain tour : Drain d'âme !)")
		damagePlayer(lich, c, lich.Attack)
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
		fmt.Println(Red + Bold + "  IGNAROK RUGIT DE RAGE ! Ses écailles virent au rouge vif. (Attaque +30 %)" + Reset)
	}

	switch turn % 3 {
	case 0:
		printArt(artBreath, fireColors...)
		fmt.Println(Red + Bold + "  Ignarok déchaîne son SOUFFLE INFERNAL !" + Reset)
		damagePlayer(dragon, c, dragon.Attack*5/2)
		applyEffect(c, dragon.Effect)
	case 2:
		fmt.Println(Orange + "  Ignarok vous lacère de ses griffes !" + Reset)
		damagePlayer(dragon, c, dragon.Attack)
		warn("Ignarok inspire profondément… Au prochain tour : SOUFFLE INFERNAL ! Protégez-vous !")
	default:
		fmt.Println(Orange + "  Ignarok balaie l'air de sa queue immense !" + Reset)
		damagePlayer(dragon, c, dragon.Attack)
	}
}

// ------------------------------------------------------------------- issues

func winFight(c *Character, m *Monster, training bool) string {
	wait(500)
	clearScreen()
	printArt(artSkull, Silver)
	banner("VICTOIRE CONTRE "+strings.ToUpper(m.Name)+" !", Green)

	// L'arène ne rapporte rien : sinon on y farmerait l'XP sans risque.
	if training {
		say("Sergent Grol", Red, "Du travail propre ! Mais l'entraînement ne paie pas : le butin, c'est au donjon.")
		pause()
		return Victory
	}

	c.Victories++
	c.Gold += m.Gold
	c.GoldEarned += m.Gold
	printArt(artChest, goldColors...)
	section("Butin", Gold)
	fmt.Printf("  %s¤ +%d Y-Coins%s   (bourse : %d Y-Coins)\n", Gold+Bold, m.Gold, Reset, c.Gold)
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
		c.HP = max(c.HP, c.MaxHP/2)
		say("Sergent Grol", Red, "Pas mal… pour un sac de patates. Relève-toi, ça ne compte pas !")
		pause()
		return Defeat
	}
	isDead(c)
	lost := c.Gold / 5
	c.Gold -= lost
	fail("Vous vous réveillez au camp, courbaturé. Il vous manque %d Y-Coins…", lost)
	pause()
	return Defeat
}
