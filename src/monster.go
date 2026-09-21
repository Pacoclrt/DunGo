// ════════════════════════════════════════════════════════════════════════
//   monster.go · [MONSTRES]
//   Tout ce qui concerne les monstres : leur structure, le bestiaire (la
//   liste de tous les monstres), la création d'un monstre prêt à combattre,
//   et la façon dont chaque monstre attaque (les boss ont leur propre tour).
// ════════════════════════════════════════════════════════════════════════

package main

import "fmt"

// ════════════════════════════════════════════════════════════════════════
//   LE MONSTRE ET LE BESTIAIRE
// ════════════════════════════════════════════════════════════════════════

type Monster struct { // [MONSTRES] sert à décrire un monstre : ses stats, son butin, son effet, son dessin et ses répliques
	Key        string   // sa clé dans le bestiaire : "wolf", "dragon"… (remplie par newMonster)
	Name       string   // son nom affiché : « Loup des cavernes »
	HP         int      // ses PV actuels
	MaxHP      int      // ses PV au départ (et au maximum)
	Attack     int      // les dégâts de son attaque normale
	XP         int      // l'XP qu'il rapporte
	YCoins     int      // les Y-Coins qu'il rapporte
	Drop       string   // l'objet qu'il peut laisser après la victoire ("" = aucun)
	DropChance int      // la chance de laisser cet objet, en % (100 = toujours)
	Effect     string   // l'effet infligé par son attaque puissante ("" = aucun)
	HitText    string   // la fin de la phrase quand il frappe : « vous mord », « vous écrase »…
	Art        string   // son dessin (art.go)
	Cry        string   // ce qu'il dit en apparaissant
	IsBoss     bool     // un boss ne laisse pas fuir
	Color      string   // la couleur de son nom : celle de l'étage où il apparaît
	ArtColors  []string // le dégradé de son dessin : celui de l'étage
	Enraged    bool     // pour Ignarok : true une fois passé sous la moitié de ses PV
	Burning    int      // les tours de brûlure qui lui restent (sort Boule de Feu)
}

// Les stats du bestiaire sont celles de la difficulté Normal. HP est laissé
// vide ici : newMonster le met à MaxHP (le monstre commence à pleine vie).
var bestiary = map[string]Monster{ // [MONSTRES] sert à ranger tous les monstres du jeu par leur clé : bestiary["wolf"] donne le loup
	"training_goblin": {Name: "Gobelin d'entraînement", MaxHP: 40, Attack: 5,
		HitText: "vous tapote", Art: artGoblin},

	// Les monstres ordinaires
	"goblin": {Name: "Gobelin chapardeur", MaxHP: 30, Attack: 5, XP: 16, YCoins: 6,
		Drop: ItemHealthPotion, DropChance: 20, HitText: "vous pique avec une fourchette", Art: artGoblin,
		Cry: "Toi donner Y-Coins ! Moi donner coups !"},
	"wolf": {Name: "Loup des cavernes", MaxHP: 28, Attack: 6, XP: 18, YCoins: 3,
		Drop: ItemWolfFur, DropChance: 75, Effect: EffectBleed, HitText: "vous mord", Art: artWolf,
		Cry: "Grrrrr… (il n'a pas envie de jouer à la baballe)"},
	"raven": {Name: "Corbeau funeste", MaxHP: 18, Attack: 4, XP: 11, YCoins: 2,
		Drop: ItemRavenFeather, DropChance: 85, HitText: "vous picore le crâne", Art: artRaven,
		Cry: "CROÂ ! CROÂÂÂ !"},
	"boar": {Name: "Sanglier furieux", MaxHP: 42, Attack: 7, XP: 24, YCoins: 4,
		Drop: ItemBoarLeather, DropChance: 75, HitText: "vous charge", Art: artBoar,
		Cry: "GROUIIIK ! (il gratte le sol et fonce)"},
	"skeleton": {Name: "Squelette maudit", MaxHP: 55, Attack: 9, XP: 32, YCoins: 9,
		Drop: ItemManaPotion, DropChance: 30, Effect: EffectWeaken, HitText: "vous frappe avec son propre fémur", Art: artSkeleton,
		Cry: "Rejoins-nous… On a des gâteaux… Enfin, on avait."},
	"troll": {Name: "Troll des cavernes", MaxHP: 85, Attack: 12, XP: 55, YCoins: 14,
		Drop: ItemTrollSkin, DropChance: 75, Effect: EffectStun, HitText: "vous écrase", Art: artTroll,
		Cry: "TROLL AVOIR FAIM. TOI AVOIR L'AIR CROUSTILLANT."},
	"krokmou": {Name: "Krokmou", MaxHP: 70, Attack: 13, XP: 60, YCoins: 17,
		Drop: ItemHealthPotion, DropChance: 35, Effect: EffectBurn, HitText: "vous mordille", Art: artKrokmou,
		Cry: "Krrrr ! (le bébé dragon crache une petite flamme… qui brûle quand même)"},

	// Les boss : un par étage
	"goblin_king": {Name: "Grukk, le Roi Gobelin", MaxHP: 110, Attack: 9, XP: 120, YCoins: 60,
		Drop: ItemFireballBook, DropChance: 100, HitText: "vous assomme avec son sceptre", Art: artGoblinKing,
		IsBoss: true, Cry: "QUI OSE ENTRER DANS MON ROYAUME ? GARDES ! GAAARDES !"},
	"lich": {Name: "Mor'Vath, la Liche", MaxHP: 150, Attack: 13, XP: 220, YCoins: 100,
		Drop: ItemManaPotion, DropChance: 100, Effect: EffectWeaken, HitText: "vous glace les os", Art: artLich,
		IsBoss: true, Cry: "Ton mana… Donne-le-moi… Le mien est périmé depuis trois siècles…"},
	"dragon": {Name: "Ignarok, le dragon", MaxHP: 220, Attack: 12, XP: 400, YCoins: 150,
		Effect: EffectBurn, HitText: "vous griffe", Art: artDragon,
		IsBoss: true, Cry: "QUI OSE INTERROMPRE MA SIESTE ? JE VAIS TE TRANSFORMER EN BROCHETTE !"},
}

func newMonster(key string, difficultyName string) *Monster { // [MONSTRES] sert à créer un monstre prêt à combattre (copie du bestiaire, à pleine vie, adaptée à la difficulté) et le renvoie
	// bestiary[key] donne une COPIE du monstre : on peut modifier cette
	// copie (PV, attaque…) sans jamais abîmer le bestiaire.
	monster := bestiary[key]
	difficulty := findDifficulty(difficultyName)

	monster.Key = key
	monster.MaxHP = monster.MaxHP * difficulty.MonsterPercent / 100
	monster.HP = monster.MaxHP
	monster.Attack = max(1, monster.Attack*difficulty.MonsterPercent/100)
	monster.XP = monster.XP * difficulty.RewardPercent / 100
	monster.YCoins = monster.YCoins * difficulty.RewardPercent / 100
	return &monster
}

// ════════════════════════════════════════════════════════════════════════
//   LE TOUR DU MONSTRE
//   Tous les monstres suivent un rythme de 3 tours : aux tours 3, 6, 9…
//   (turn%3 == 0), ils font leur attaque spéciale.
// ════════════════════════════════════════════════════════════════════════

func monsterTurn(monster *Monster, character *Character, turn int) { // [MONSTRES] sert à faire jouer le monstre : choisit sa façon d'attaquer selon sa clé (chaque boss a la sienne)
	switch monster.Key {
	case "goblin_king":
		goblinKingTurn(monster, character, turn)
	case "lich":
		lichTurn(monster, character, turn)
	case "dragon":
		dragonTurn(monster, character, turn)
	default:
		normalMonsterTurn(monster, character, turn)
	}
}

func normalMonsterTurn(monster *Monster, character *Character, turn int) { // [MONSTRES] sert à faire attaquer un monstre ordinaire : attaque normale, et attaque × 2 + effet tous les 3 tours
	if turn%3 == 0 {
		fmt.Println(Orange + "  ATTAQUE PUISSANTE !" + Reset)
		monsterHitsCharacter(monster, character, monster.HitText, monster.Attack*2)
		inflictEffect(character, monster.Effect)
	} else {
		monsterHitsCharacter(monster, character, monster.HitText, monster.Attack)
	}
}

func goblinKingTurn(monster *Monster, character *Character, turn int) { // [MONSTRES] sert à faire attaquer Grukk : tous les 3 tours, un garde le soigne (+20 PV) et il frappe × 1,5
	if turn%3 == 0 {
		monster.HP = min(monster.HP+20, monster.MaxHP)
		warn("Un garde accourt et soigne son roi : Grukk +20 PV !")
		monsterHitsCharacter(monster, character, monster.HitText, monster.Attack*3/2)
	} else {
		monsterHitsCharacter(monster, character, monster.HitText, monster.Attack)
	}
}

func lichTurn(monster *Monster, character *Character, turn int) { // [MONSTRES] sert à faire attaquer Mor'Vath : tous les 3 tours, elle vole jusqu'à 15 mana, se soigne du double et affaiblit le héros
	if turn%3 == 0 {
		stolen := min(character.Mana, 15) // on ne peut pas voler plus de mana que le héros n'en a
		character.Mana -= stolen
		monster.HP = min(monster.HP+stolen*2, monster.MaxHP)
		fmt.Printf("  %sDRAIN D'ÂME ! Vous perdez %d mana, la Liche regagne %d PV.%s\n", Purple, stolen, stolen*2, Reset)
		monsterHitsCharacter(monster, character, monster.HitText, monster.Attack)
		inflictEffect(character, monster.Effect)
	} else {
		monsterHitsCharacter(monster, character, monster.HitText, monster.Attack)
	}
}

func dragonTurn(monster *Monster, character *Character, turn int) { // [MONSTRES] sert à faire attaquer Ignarok : rage sous 50 % de PV (+30 % d'attaque) et Souffle Infernal (× 2,5 + brûlure) tous les 3 tours
	// 1. La rage : une seule fois, dès qu'il passe sous la moitié de ses PV
	if !monster.Enraged && monster.HP <= monster.MaxHP/2 {
		monster.Enraged = true
		monster.Attack = monster.Attack * 13 / 10
		fmt.Println(Red + Bold + "  IGNAROK EST FURIEUX ! (attaque +30 %)" + Reset)
	}

	// 2. L'attaque
	if turn%3 == 0 {
		printArt(artBreath, fireColors...)
		fmt.Println(Red + Bold + "  SOUFFLE INFERNAL !" + Reset)
		monsterHitsCharacter(monster, character, "vous fait rôtir", monster.Attack*5/2)
		inflictEffect(character, monster.Effect)
	} else {
		monsterHitsCharacter(monster, character, monster.HitText, monster.Attack)
	}
}
