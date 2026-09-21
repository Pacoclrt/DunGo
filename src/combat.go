// ════════════════════════════════════════════════════════════════════════
//   combat.go · [COMBAT]
//   La boucle de combat tour par tour, le tour du héros (attaque, sorts,
//   objets, fuite), les dégâts, et la fin du combat (victoire ou défaite).
//
//   Les autres morceaux du combat sont rangés ailleurs :
//     • le tour du monstre et des boss  → monster.go  [MONSTRES]
//     • les sorts                       → spells.go   [SORTS]
//     • saignement, brûlure…            → effects.go  [EFFETS]
// ════════════════════════════════════════════════════════════════════════

package main

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

const ( // [COMBAT] sert à nommer les 3 façons dont un combat peut finir (c'est ce que renvoie fight)
	Victory = "victoire"
	Defeat  = "défaite"
	Fled    = "fuite"
)

// ════════════════════════════════════════════════════════════════════════
//   LA BOUCLE DE COMBAT
// ════════════════════════════════════════════════════════════════════════

func fight(character *Character, monster *Monster, training bool) string { // [COMBAT] sert à faire tout le combat, tour par tour, et renvoie comment il finit : Victory, Defeat ou Fled
	clearEffects(character)
	for turn := 1; ; turn++ { // pas de condition : on boucle jusqu'à un « return »
		clearScreen()
		showFightScreen(character, monster, turn)

		// 1. Le héros joue
		fled := characterTurn(character, monster, training)
		if fled {
			return Fled
		}

		// 2. Le monstre riposte, s'il est encore debout
		if monster.HP > 0 {
			monsterTurn(monster, character, turn)
		}

		// 3. Les effets de fin de tour (saignement, brûlure), si les deux sont encore debout
		if monster.HP > 0 && character.HP > 0 {
			endOfTurnEffects(character, monster)
		}

		// 4. Quelqu'un est tombé ? On regarde le héros EN PREMIER : si les
		//    deux tombent en même temps (à cause des effets), c'est une
		//    défaite. Le héros ne peut donc jamais gagner avec 0 PV.
		if character.HP <= 0 {
			return loseFight(character, training)
		}
		if monster.HP <= 0 {
			return winFight(character, monster, training)
		}
		pause()
	}
}

func showFightScreen(character *Character, monster *Monster, turn int) { // [COMBAT] sert à afficher l'écran de combat : le numéro du tour, le monstre (dessin + PV) et le héros (PV, mana, effets)
	fmt.Printf("  %s── Tour %d %s%s\n", DarkGray, turn, strings.Repeat("─", 50), Reset)
	printArt(monster.Art, monster.ArtColors...)

	// Le monstre
	fmt.Print("  " + Bold + monster.Color + monster.Name + Reset)
	if monster.IsBoss {
		fmt.Print(Red + "  ☠ BOSS" + Reset)
	}
	if monster.Enraged {
		fmt.Print(Red + "  [EN RAGE]" + Reset)
	}
	if monster.Burning > 0 {
		fmt.Printf("%s  [Brûlure %d]%s", Orange, monster.Burning, Reset)
	}
	fmt.Printf("\n  %s♥%s %s %d / %d\n", Red, Reset, hpBar(monster.HP, monster.MaxHP, 40), monster.HP, monster.MaxHP)

	// Le héros
	fmt.Println(DarkGray + "  " + strings.Repeat("─", 60) + Reset)
	fmt.Println("  " + Gold + Bold + character.Name + Reset + effectBadges(character))
	fmt.Printf("  %s♥%s %s %d / %d    %s♦%s %s %d / %d\n",
		Red, Reset, hpBar(character.HP, character.MaxHP, 20), character.HP, character.MaxHP,
		Blue, Reset, bar(character.Mana, character.MaxMana, 12, Blue), character.Mana, character.MaxMana)
}

// ════════════════════════════════════════════════════════════════════════
//   LE TOUR DU HÉROS
// ════════════════════════════════════════════════════════════════════════

func characterTurn(character *Character, monster *Monster, training bool) bool { // [COMBAT] sert à faire jouer le héros (attaquer, sort, objet ou fuir) et renvoie true s'il quitte le combat
	// Étourdi : il perd son tour, sans rien pouvoir choisir.
	if character.Stunned > 0 {
		character.Stunned--
		warn("Vous êtes étourdi et passez votre tour.")
		return false
	}

	// On repropose le menu tant que le héros n'a pas VRAIMENT joué : un
	// « Retour » dans les sorts ou une fuite interdite ne coûtent pas le tour.
	for {
		fmt.Println()
		optionHint(1, "Attaquer", fmt.Sprintf("%d dégâts", character.Attack))
		optionHint(2, "Sorts", fmt.Sprintf("%d mana", character.Mana))
		option(3, "Objets")
		if training {
			option(4, "Abandonner")
		} else {
			option(4, "Fuir")
		}

		switch readChoice(1, 4) {
		case 1: // Attaquer
			fmt.Println()
			characterHitsMonster(character, monster, character.Attack)
			return false

		case 2: // Sorts
			if spellMenu(character, monster) {
				return false
			}

		case 3: // Objets (interdits à l'entraînement)
			if training {
				fail("Pas d'objets à l'entraînement : le sergent Grol veut voir vos poings !")
			} else if combatItemMenu(character, monster) {
				return false
			}

		case 4: // Fuir (ou abandonner l'entraînement)
			if training {
				info("Vous abandonnez l'entraînement.")
				pause()
				return true
			}
			if monster.IsBoss {
				fail("On ne fuit pas un boss !")
			} else {
				return tryToFlee(monster)
			}
		}
	}
}

func tryToFlee(monster *Monster) bool { // [COMBAT] sert à tenter de fuir (1 chance sur 2) et renvoie true si la fuite réussit (sinon, le tour est perdu)
	// rand.IntN(100) tire un nombre au hasard entre 0 et 99 :
	// « < 50 » arrive 50 fois sur 100.
	if rand.IntN(100) < 50 {
		success("Vous courez jusqu'au camp. Avec beaucoup de dignité.")
		pause()
		return true
	}
	fail("%s vous rattrape. La fuite échoue !", monster.Name)
	return false
}

func combatItemMenu(character *Character, monster *Monster) bool { // [COMBAT] sert à afficher les objets utilisables en combat et à utiliser celui choisi ; renvoie true si un objet a vraiment servi
	// On ne garde que les objets du sac utilisables en combat (potions, livre).
	usable := []string{}
	for _, item := range sortedInventory(character) {
		if items[item].UsableInFight {
			usable = append(usable, item)
		}
	}
	if len(usable) == 0 {
		fail("Aucun objet utilisable en combat.")
		return false
	}

	section("Objets")
	for index, item := range usable {
		option(index+1, fmt.Sprintf("%s x%d", item, character.Inventory[item]))
	}
	backOption("Retour")

	choice := readChoice(0, len(usable))
	if choice == 0 {
		return false
	}
	fmt.Println()
	return useItem(character, monster, usable[choice-1])
}

// ════════════════════════════════════════════════════════════════════════
//   LES DÉGÂTS
//     characterHitsMonster  le héros frappe  (critique, affaiblissement)
//     removeMonsterHP       dégâts bruts     (poison, brûlure)
//     monsterHitsCharacter  le monstre frappe (Peau de Pierre)
// ════════════════════════════════════════════════════════════════════════

func characterHitsMonster(character *Character, monster *Monster, damage int) { // [COMBAT] sert à appliquer un coup du héros (attaque ou sort) : ÷ 2 s'il est affaibli, puis 10 % de chance de coup critique (× 2)
	note := ""
	if character.Weakened > 0 {
		character.Weakened--
		damage = max(1, damage/2) // au moins 1 dégât
		note = note + " (affaibli)"
	}
	if rand.IntN(100) < 10 {
		damage = damage * 2
		note = note + " ★ CRITIQUE !"
	}
	removeMonsterHP(monster, damage, note)
}

func removeMonsterHP(monster *Monster, damage int, note string) { // [COMBAT] sert à retirer des PV au monstre (sans critique ni bonus) et à l'afficher ; les PV ne descendent jamais sous 0
	monster.HP = max(monster.HP-damage, 0)
	fmt.Printf("  %s» %s perd %d PV%s%s\n", Gold, monster.Name, damage, note, Reset)
}

func monsterHitsCharacter(monster *Monster, character *Character, hitText string, damage int) { // [COMBAT] sert à appliquer un coup du monstre sur le héros : ÷ 2 si Peau de Pierre est active
	note := ""
	if character.StoneSkin > 0 {
		character.StoneSkin--
		damage = max(1, damage/2)
		note = " (Peau de Pierre)"
	}
	character.HP = max(character.HP-damage, 0)
	fmt.Printf("  %s« %s %s : -%d PV%s%s\n", Red, monster.Name, hitText, damage, note, Reset)
}

// ════════════════════════════════════════════════════════════════════════
//   LA FIN DU COMBAT
// ════════════════════════════════════════════════════════════════════════

func winFight(character *Character, monster *Monster, training bool) string { // [COMBAT] sert à gérer la victoire : butin, XP et mission (rien à l'entraînement), et renvoie Victory
	clearScreen()
	banner("VICTOIRE CONTRE "+strings.ToUpper(monster.Name)+" !", Gold)

	// L'entraînement ne rapporte rien : sinon on y gagnerait des niveaux sans risque.
	if training {
		fmt.Println()
		say("Sergent Grol", Gold, "Propre ! Mais le butin, c'est au donjon.")
		pause()
		return Victory
	}

	// 1. Les Y-Coins
	character.Victories++
	character.YCoins += monster.YCoins
	printArt(artChest, campColors...)
	section("Butin")
	fmt.Printf("  %s¤ +%d Y-Coins%s\n", Gold, monster.YCoins, Reset)

	// 2. L'objet laissé, avec DropChance % de chance
	if monster.Drop != "" && rand.IntN(100) < monster.DropChance {
		if addToInventory(character, monster.Drop, 1) {
			fmt.Println("  ■ " + monster.Drop)
		}
	}

	// 3. L'XP, puis la mission en cours
	gainXP(character, monster.XP)
	updateQuest(character, monster.Name)
	pause()
	return Victory
}

func loseFight(character *Character, training bool) string { // [COMBAT] sert à gérer la défaite : la mort du héros (rien à l'entraînement), et renvoie Defeat
	if training {
		fmt.Println()
		say("Sergent Grol", Gold, "Battu par le gobelin d'entraînement… Relève-toi, personne n'a rien vu.")
		pause()
		return Defeat
	}
	characterDies(character)
	pause()
	return Defeat
}
