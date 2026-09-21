package main

// Monster décrit un adversaire. Les valeurs de base vivent dans bestiary ;
// newMonster en renvoie une copie, si bien qu'un combat ne modifie jamais
// le bestiaire lui-même.
type Monster struct {
	Name       string
	MaxHP      int
	HP         int
	Attack     int
	XP         int
	Gold       int
	Drop       string
	DropChance int    // probabilité du butin, en %
	Effect     string // infligé par l'attaque puissante
	Hit        string // verbe du journal de combat : « Loup des cavernes vous mord »
	Art        string
	Cry        string
	IsBoss     bool
	// Pattern décide de l'attaque du tour. Nil = schéma commun (attaque
	// puissante tous les 3 tours).
	Pattern func(*Monster, *Character, int)
	// Les couleurs ne sont pas dans le bestiaire : un monstre prend celles
	// du lieu où il apparaît (voir floors).
	Color   string
	Colors  []string // dégradé de l'entrée en scène des boss
	Enraged bool
	Burning int
}

var bestiary = map[string]Monster{
	"training_goblin": {Name: "Gobelin d'entraînement", MaxHP: 40, Attack: 5, XP: 12, Gold: 2,
		Hit: "vous tapote", Art: artGoblin, Cry: "Sergent Grol a dit pas taper trop fort…"},
	"goblin": {Name: "Gobelin chapardeur", MaxHP: 30, Attack: 5, XP: 16, Gold: 6,
		Drop: ItemHealthPotion, DropChance: 20, Hit: "vous pique avec une fourchette", Art: artGoblin,
		Cry: "Toi donner Y-Coins ! Moi donner coups !"},
	"wolf": {Name: "Loup des cavernes", MaxHP: 28, Attack: 6, XP: 18, Gold: 3,
		Drop: ItemWolfFur, DropChance: 75, Effect: EffectBleed, Hit: "vous mord", Art: artWolf,
		Cry: "Grrrrr… (il n'a clairement pas envie de jouer à la baballe)"},
	"raven": {Name: "Corbeau funeste", MaxHP: 18, Attack: 4, XP: 11, Gold: 2,
		Drop: ItemRavenFeather, DropChance: 85, Hit: "vous picore le crâne", Art: artRaven,
		Cry: "CROÂ ! CROÂÂÂ ! (traduction : « donne ton goûter »)"},
	"boar": {Name: "Sanglier furieux", MaxHP: 42, Attack: 7, XP: 24, Gold: 4,
		Drop: ItemBoarLeather, DropChance: 75, Hit: "vous charge", Art: artBoar,
		Cry: "GROUIIIK ! (il gratte le sol et fonce)"},
	"skeleton": {Name: "Squelette maudit", MaxHP: 55, Attack: 9, XP: 32, Gold: 9,
		Drop: ItemManaPotion, DropChance: 30, Effect: EffectWeaken, Hit: "vous frappe avec son propre fémur", Art: artSkeleton,
		Cry: "Rejoins-nous… On a des gâteaux… Enfin, on avait."},
	"troll": {Name: "Troll des cavernes", MaxHP: 85, Attack: 12, XP: 55, Gold: 14,
		Drop: ItemTrollSkin, DropChance: 75, Effect: EffectStun, Hit: "vous écrase", Art: artTroll,
		Cry: "TROLL AVOIR FAIM. TOI AVOIR L'AIR CROUSTILLANT."},
	"krokmou": {Name: "Krokmou", MaxHP: 70, Attack: 13, XP: 60, Gold: 17,
		Drop: ItemHealthPotion, DropChance: 35, Effect: EffectBurn, Hit: "vous mordille", Art: artKrokmou,
		Cry: "Krrrr ! (le bébé dragon crache une petite flamme… qui brûle quand même)"},

	"goblin_king": {Name: "Grukk, le Roi Gobelin", MaxHP: 110, Attack: 9, XP: 120, Gold: 60,
		Drop: ItemFireballBook, DropChance: 100, Hit: "vous assomme avec son sceptre", Art: artGoblinKing,
		IsBoss: true, Pattern: goblinKingPattern,
		Cry: "QUI OSE ENTRER DANS MON ROYAUME ? GARDES ! GAAARDES ! … Ils sont en pause ?"},
	"lich": {Name: "Mor'Vath, la Liche", MaxHP: 150, Attack: 13, XP: 220, Gold: 100,
		Drop: ItemManaPotion, DropChance: 100, Effect: EffectWeaken, Hit: "vous glace les os", Art: artLich,
		IsBoss: true, Pattern: lichPattern,
		Cry: "Ton mana… Donne-le-moi… Le mien est périmé depuis trois siècles…"},
	"dragon": {Name: "Ignarok, le dragon", MaxHP: 220, Attack: 12, XP: 400, Gold: 150,
		Effect: EffectBurn, Hit: "vous griffe", Art: artDragon,
		IsBoss: true, Pattern: dragonPattern,
		Cry: "QUI OSE INTERROMPRE MA SIESTE ? JE VAIS TE TRANSFORMER EN BROCHETTE !"},
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
