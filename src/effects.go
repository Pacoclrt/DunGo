// ════════════════════════════════════════════════════════════════════════
//   effects.go · [EFFETS]
//   Les effets de combat : saignement, brûlure, étourdissement,
//   affaiblissement et Peau de Pierre.
//
//   Un effet n'est qu'un compteur (un nombre) rangé dans le héros
//   (Character.Bleeding…) ou dans le monstre (Monster.Burning) :
//     • inflictEffect   met le compteur en place (3 tours, 1 tour…)
//     • le compteur baisse de 1 à chaque fois qu'il sert
//     • à 0, l'effet est terminé
// ════════════════════════════════════════════════════════════════════════

package main

import "fmt"

const ( // [EFFETS] sert à nommer les 4 effets qu'un monstre peut infliger avec son attaque puissante (champ Monster.Effect)
	EffectBleed  = "saignement"
	EffectBurn   = "brûlure"
	EffectStun   = "étourdissement"
	EffectWeaken = "affaiblissement"
)

func inflictEffect(character *Character, effect string) { // [EFFETS] sert à infliger un effet au héros (met son compteur de tours en place) et à l'annoncer
	switch effect {
	case EffectBleed:
		character.Bleeding = 3
		warn("Vous saignez ! (-3 PV par tour, 3 tours)")
	case EffectBurn:
		character.Burning = 3
		warn("Vous brûlez ! (-5 PV par tour, 3 tours)")
	case EffectStun:
		character.Stunned = 1
		warn("Vous êtes étourdi ! (vous passez votre prochain tour)")
	case EffectWeaken:
		character.Weakened = 3
		warn("Vous êtes affaibli ! (vos 3 prochaines attaques ÷ 2)")
	}
	// Un monstre sans effet a Effect = "" : aucun case ne correspond, il ne se passe rien.
}

func endOfTurnEffects(character *Character, monster *Monster) { // [EFFETS] sert à appliquer les dégâts de fin de tour (saignement, brûlure) et à faire baisser leurs compteurs
	if character.Bleeding > 0 {
		character.Bleeding--
		character.HP = max(character.HP-3, 0)
		fmt.Println(Red + "  « Saignement : -3 PV" + Reset)
	}
	if character.Burning > 0 {
		character.Burning--
		character.HP = max(character.HP-5, 0)
		fmt.Println(Orange + "  « Brûlure : -5 PV" + Reset)
	}
	// Le monstre peut brûler lui aussi : c'est l'effet du sort Boule de Feu.
	if monster.Burning > 0 {
		monster.Burning--
		removeMonsterHP(monster, 5, " (brûlure)")
	}
}

func clearEffects(character *Character) { // [EFFETS] sert à remettre tous les effets du héros à 0 (au début de chaque combat)
	character.Bleeding = 0
	character.Burning = 0
	character.Stunned = 0
	character.Weakened = 0
	character.StoneSkin = 0
}

func effectBadges(character *Character) string { // [EFFETS] sert à fabriquer les étiquettes des effets actifs, ex : [Saignement 2] [Affaibli 3], et les renvoie
	badges := ""
	if character.Bleeding > 0 {
		badges += fmt.Sprintf(" %s[Saignement %d]%s", Red, character.Bleeding, Reset)
	}
	if character.Burning > 0 {
		badges += fmt.Sprintf(" %s[Brûlure %d]%s", Orange, character.Burning, Reset)
	}
	if character.Stunned > 0 {
		badges += fmt.Sprintf(" %s[Étourdi]%s", Yellow, Reset)
	}
	if character.Weakened > 0 {
		badges += fmt.Sprintf(" %s[Affaibli %d]%s", Purple, character.Weakened, Reset)
	}
	if character.StoneSkin > 0 {
		badges += fmt.Sprintf(" %s[Peau de Pierre %d]%s", Silver, character.StoneSkin, Reset)
	}
	return badges
}
