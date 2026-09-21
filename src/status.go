package main

import "fmt"

// Effets de statut infligés par l'attaque puissante de certains monstres.
const (
	EffectBleed  = "saignement"
	EffectBurn   = "brûlure"
	EffectStun   = "étourdissement"
	EffectWeaken = "affaiblissement"
)

func applyEffect(c *Character, effect string) {
	if effect == "" || c.TestMode {
		return
	}
	switch effect {
	case EffectBleed:
		c.Bleeding = 3
		warn("Vous saignez ! -3 PV à la fin de chaque tour, pendant 3 tours.")
	case EffectBurn:
		c.Burning = 3
		warn("Vous brûlez ! -5 PV à la fin de chaque tour, pendant 3 tours.")
	case EffectStun:
		c.Stunned = 1
		warn("Vous êtes étourdi ! Vous passerez votre prochain tour.")
	case EffectWeaken:
		c.Weakened = 3
		warn("Vous êtes affaibli ! Vos 3 prochaines attaques font moitié moins de dégâts.")
	}
}

// updateEffects applique les dégâts de fin de tour, côté héros et côté monstre.
func updateEffects(c *Character, m *Monster) {
	if c.Bleeding > 0 {
		c.Bleeding--
		c.HP = max(c.HP-3, 0)
		fmt.Printf("  %s♦ Saignement : -3 PV%s\n", Red, Reset)
		showHP("", c.HP, c.MaxHP)
	}
	if c.Burning > 0 {
		c.Burning--
		c.HP = max(c.HP-5, 0)
		fmt.Printf("  %s♦ Brûlure : -5 PV%s\n", Orange, Reset)
		showHP("", c.HP, c.MaxHP)
	}
	if m.Burning > 0 {
		m.Burning--
		m.HP = max(m.HP-5, 0)
		fmt.Printf("  %s♦ %s brûle : -5 PV%s\n", Orange, m.Name, Reset)
		showHP(m.Name, m.HP, m.MaxHP)
	}
}

func clearEffects(c *Character) {
	c.Bleeding, c.Burning, c.Stunned, c.Weakened, c.StoneSkin = 0, 0, 0, 0, 0
}

// effectBadges résume les effets actifs en une ligne, affichée en combat.
func effectBadges(c *Character) string {
	badges := ""
	for _, badge := range []struct {
		Turns int
		Color string
		Label string
	}{
		{c.Bleeding, Red, "Saignement"},
		{c.Burning, Orange, "Brûlure"},
		{c.Stunned, Yellow, "Étourdi"},
		{c.Weakened, Purple, "Affaibli"},
		{c.StoneSkin, Silver, "Peau de Pierre"},
	} {
		if badge.Turns > 0 {
			badges += fmt.Sprintf(" %s[%s %d]%s", badge.Color, badge.Label, badge.Turns, Reset)
		}
	}
	return badges
}
