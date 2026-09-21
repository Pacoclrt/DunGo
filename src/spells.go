// ════════════════════════════════════════════════════════════════════════
//   spells.go · [SORTS]
//   Les 5 sorts du jeu : leur table (coût, dégâts, dessin), le menu des
//   sorts en combat et le lancement d'un sort.
// ════════════════════════════════════════════════════════════════════════

package main

import "fmt"

// Les noms des sorts sont des constantes : si on se trompe dans un nom
// (SpellFirebal), le programme ne compile pas, au lieu de bugger en silence.
const ( // [SORTS] sert à nommer les 5 sorts (ces noms sont aussi ceux affichés dans le jeu)
	SpellPunch       = "Coup de poing"   // tout le monde le connaît
	SpellSecondWind  = "Second Souffle"  // sort de l'Humain
	SpellSilverArrow = "Flèche d'Argent" // sort de l'Elfe
	SpellStoneSkin   = "Peau de Pierre"  // sort du Nain
	SpellFireball    = "Boule de Feu"    // s'apprend avec le Livre de Sort
)

type Spell struct { // [SORTS] sert à décrire un sort : son coût en mana, ses dégâts et son dessin
	ManaCost    int
	Description string // le texte affiché dans le menu des sorts
	Damage      int    // 0 = le sort ne fait pas de dégâts (soin, protection)
	Art         string
	ArtColors   []string
}

var spells = map[string]Spell{ // [SORTS] sert à ranger les 5 sorts par leur nom : spells["Boule de Feu"] donne le sort
	SpellPunch:       {ManaCost: 5, Description: "8 dégâts", Damage: 8, Art: artFist, ArtColors: []string{Orange}},
	SpellSilverArrow: {ManaCost: 8, Description: "13 dégâts", Damage: 13, Art: artArrow, ArtColors: stoneColors},
	SpellSecondWind:  {ManaCost: 10, Description: "rend 25 PV et soigne saignement et brûlure", Damage: 0, Art: artHeal, ArtColors: []string{Green}},
	SpellStoneSkin:   {ManaCost: 8, Description: "dégâts reçus ÷ 2 pendant 2 coups", Damage: 0, Art: artShield, ArtColors: stoneColors},
	SpellFireball:    {ManaCost: 15, Description: "18 dégâts + brûlure", Damage: 18, Art: artFireball, ArtColors: fireColors},
}

func spellMenu(character *Character, monster *Monster) bool { // [SORTS] sert à afficher les sorts du héros et à lancer celui choisi ; renvoie true si un sort a été lancé (false = Retour ou pas assez de mana)
	section("Sorts")
	for index, name := range character.Spells {
		spell := spells[name]
		option(index+1, fmt.Sprintf("%-16s %s%2d mana%s  %s%s%s", name, Blue, spell.ManaCost, Reset, Gray, spell.Description, Reset))
	}
	backOption("Retour")

	choice := readChoice(0, len(character.Spells))
	if choice == 0 {
		return false
	}
	name := character.Spells[choice-1]
	if character.Mana < spells[name].ManaCost {
		fail("Pas assez de mana (%d / %d).", character.Mana, spells[name].ManaCost)
		return false
	}
	castSpell(character, monster, name)
	return true
}

func castSpell(character *Character, monster *Monster, name string) { // [SORTS] sert à lancer un sort : retire le mana, affiche le dessin, fait les dégâts puis l'effet spécial du sort
	spell := spells[name]

	// 1. Le coût et le dessin
	character.Mana -= spell.ManaCost
	fmt.Printf("\n  %s✦ %s !%s\n", Purple+Bold, name, Reset)
	printArt(spell.Art, spell.ArtColors...)

	// 2. Les dégâts : comme une attaque (critique possible, ÷ 2 si affaibli)
	if spell.Damage > 0 {
		characterHitsMonster(character, monster, spell.Damage)
	}

	// 3. L'effet spécial, propre à chaque sort
	switch name {
	case SpellFireball:
		monster.Burning = 3
		warn("%s prend feu ! (-5 PV par tour pendant 3 tours)", monster.Name)
	case SpellSecondWind:
		character.HP = min(character.HP+25, character.MaxHP)
		character.Bleeding = 0
		character.Burning = 0
		success("+25 PV, saignement et brûlure soignés.")
	case SpellStoneSkin:
		character.StoneSkin = 2
		success("Votre peau devient pierre : les 2 prochains coups font moitié moins mal.")
	}
}
