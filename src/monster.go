package main

// Monster décrit un adversaire. Les valeurs de base vivent dans bestiary ;
// newMonster en renvoie une copie, si bien qu'un combat ne modifie jamais
// le bestiaire lui-même.
type Monster struct {
	Name       string
	MaxHP      int
	HP         int
	Attack     int
	Initiative int
	XP         int
	Gold       int
	Drop       string
	DropChance int    // probabilité du butin, en %
	Effect     string // infligé par l'attaque puissante
	Art        string
	Color      string
	Colors     []string // dégradé de l'entrée en scène des boss
	Cry        string
	IsBoss     bool
	// Pattern décide de l'attaque du tour. Nil = schéma commun (attaque
	// puissante tous les 3 tours).
	Pattern func(*Monster, *Character, int)
	Enraged bool
	Burning int
}

var bestiary = map[string]Monster{
	"training_goblin": {Name: "Gobelin d'entrainement", MaxHP: 40, Attack: 5, Initiative: 8, XP: 12, Gold: 2,
		Art: artGoblin, Color: Green, Cry: "Sergent Grol a dit pas taper trop fort…"},
	"goblin": {Name: "Gobelin chapardeur", MaxHP: 30, Attack: 5, Initiative: 9, XP: 16, Gold: 6,
		Drop: ItemHealthPotion, DropChance: 20, Art: artGoblin, Color: Green,
		Cry: "Toi donner Y-Coins ! Moi donner coups !"},
	"wolf": {Name: "Loup des cavernes", MaxHP: 28, Attack: 6, Initiative: 15, XP: 18, Gold: 3,
		Drop: ItemWolfFur, DropChance: 75, Effect: EffectBleed, Art: artWolf, Color: Gray,
		Cry: "Grrrrr… (il montre les crocs)"},
	"raven": {Name: "Corbeau funeste", MaxHP: 18, Attack: 4, Initiative: 18, XP: 11, Gold: 2,
		Drop: ItemRavenFeather, DropChance: 85, Art: artRaven, Color: Purple,
		Cry: "CROÂ ! CROÂÂÂ !"},
	"boar": {Name: "Sanglier furieux", MaxHP: 42, Attack: 7, Initiative: 6, XP: 24, Gold: 4,
		Drop: ItemBoarLeather, DropChance: 75, Art: artBoar, Color: Yellow,
		Cry: "GROUIIIK ! (il gratte le sol et charge)"},
	"skeleton": {Name: "Squelette maudit", MaxHP: 55, Attack: 9, Initiative: 9, XP: 32, Gold: 9,
		Drop: ItemManaPotion, DropChance: 30, Effect: EffectWeaken, Art: artSkeleton, Color: Gray,
		Cry: "Rejoins-nous… dans la tombe…"},
	"troll": {Name: "Troll des cavernes", MaxHP: 85, Attack: 12, Initiative: 4, XP: 55, Gold: 14,
		Drop: ItemTrollSkin, DropChance: 75, Effect: EffectStun, Art: artTroll, Color: Green,
		Cry: "TROLL AVOIR FAIM. TOI AVOIR L'AIR CROUSTILLANT."},
	"krokmou": {Name: "Krokmou", MaxHP: 70, Attack: 13, Initiative: 11, XP: 60, Gold: 17,
		Drop: ItemHealthPotion, DropChance: 35, Effect: EffectBurn, Art: artKrokmou, Color: Red,
		Cry: "Krrrr ! (le dragonneau crache une petite flamme… qui brûle quand même)"},

	"goblin_king": {Name: "Grukk, le Roi Gobelin", MaxHP: 110, Attack: 9, Initiative: 10, XP: 120, Gold: 60,
		Drop: ItemFireballBook, DropChance: 100, Art: artGoblinKing, Color: Gold, IsBoss: true,
		Colors: goldColors, Pattern: goblinKingPattern,
		Cry:    "QUI OSE ENTRER DANS MON ROYAUME ? GARDES ! GAAARDES !"},
	"lich": {Name: "Mor'Vath, la Liche", MaxHP: 150, Attack: 13, Initiative: 13, XP: 220, Gold: 100,
		Drop: ItemManaPotion, DropChance: 100, Effect: EffectWeaken, Art: artLich, Color: Purple, IsBoss: true,
		Colors: spiritColors, Pattern: lichPattern,
		Cry:    "Ta vie… ton mana… tout m'appartiendra…"},
	"dragon": {Name: "Ignarok, le Fléau Écarlate", MaxHP: 220, Attack: 12, Initiative: 12, XP: 400, Gold: 150,
		Effect: EffectBurn, Art: artDragon, Color: Red, IsBoss: true,
		Colors: fireColors, Pattern: dragonPattern,
		Cry:    "UN INSECTE DE PLUS DANS MON ANTRE… JE VAIS TE RÉDUIRE EN CENDRES !"},
}

// newMonster fabrique un monstre neuf, à pleine vie, adapté à la difficulté.
func newMonster(kind, difficulty string) *Monster {
	m := bestiary[kind]
	level := difficultyOf(difficulty)
	m.MaxHP = m.MaxHP * level.Power / 100
	m.Attack = max(1, m.Attack*level.Power/100)
	m.Gold = m.Gold * level.Reward / 100
	m.XP = m.XP * level.Reward / 100
	m.HP = m.MaxHP
	return &m
}
