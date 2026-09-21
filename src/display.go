// ════════════════════════════════════════════════════════════════════════
//   display.go · [AFFICHAGE]
//   Tout ce qui s'affiche à l'écran : couleurs, dessins, cadres, menus,
//   messages et jauges. Tous les écrans du jeu passent par ces fonctions :
//   c'est pour ça que le style est le même partout.
// ════════════════════════════════════════════════════════════════════════

package main

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

// ════════════════════════════════════════════════════════════════════════
//   LES COULEURS
// ════════════════════════════════════════════════════════════════════════

// Un code ANSI est un petit texte que le terminal comprend comme un ordre
// au lieu de l'afficher. "\033[38;5;196m" veut dire : « écris avec la
// couleur n° 196 » (le terminal en connaît 256, numérotées de 0 à 255).
const ( // [AFFICHAGE] sert à donner un nom aux couleurs du texte (toute couleur ouverte doit être refermée par Reset)
	Reset  = "\033[0m" // revient au style normal
	Bold   = "\033[1m" // gras
	Italic = "\033[3m" // italique

	Silver   = "\033[38;5;250m"
	Gray     = "\033[38;5;244m"
	DarkGray = "\033[38;5;238m"
	Gold     = "\033[38;5;220m"
	Yellow   = "\033[38;5;227m"
	Orange   = "\033[38;5;208m"
	Red      = "\033[38;5;196m"
	Green    = "\033[38;5;77m"
	Sky      = "\033[38;5;117m"
	Blue     = "\033[38;5;33m"
	Purple   = "\033[38;5;135m"
)

// Un dégradé = 6 nuances d'une même couleur, de la plus claire à la plus
// foncée. printArt les applique aux dessins, du haut vers le bas.
var ( // [AFFICHAGE] sert à définir les 5 dégradés utilisés pour colorer les dessins
	campColors  = shades(229, 228, 222, 220, 214, 172) // doré : le camp
	stoneColors = shades(255, 253, 251, 249, 247, 245) // gris pierre
	mossColors  = shades(190, 154, 118, 82, 40, 34)    // vert : étage 1
	iceColors   = shades(195, 159, 123, 81, 75, 33)    // bleu glacé : étage 2
	fireColors  = shades(227, 220, 214, 208, 202, 196) // feu : étage 3 et dragon
)

func shades(colorNumbers ...int) []string { // [AFFICHAGE] sert à transformer des numéros de couleur (0 à 255) en codes ANSI, et renvoie la liste
	colors := []string{}
	for _, number := range colorNumbers {
		colors = append(colors, fmt.Sprintf("\033[38;5;%dm", number))
	}
	return colors
}

// ════════════════════════════════════════════════════════════════════════
//   L'ÉCRAN ET LES DESSINS
// ════════════════════════════════════════════════════════════════════════

func clearScreen() { // [AFFICHAGE] sert à effacer tout l'écran (code ANSI : curseur en haut à gauche + effacement)
	fmt.Print("\033[H\033[2J")
}

func wait(milliseconds int) { // [AFFICHAGE] sert à mettre le jeu en pause quelques millisecondes (effet de suspense)
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)
}

func printArt(art string, colors ...string) { // [AFFICHAGE] sert à afficher un dessin, avec un dégradé de couleurs du haut vers le bas
	if len(colors) == 0 {
		colors = []string{Silver}
	}
	lines := strings.Split(strings.Trim(art, "\n"), "\n")
	for index, line := range lines {
		// Chaque couleur couvre une part égale des lignes :
		// 12 lignes et 6 couleurs → chaque couleur colore 2 lignes.
		color := colors[index*len(colors)/len(lines)]
		fmt.Println(Bold + color + line + Reset)
	}
}

// ════════════════════════════════════════════════════════════════════════
//   LE TEXTE
// ════════════════════════════════════════════════════════════════════════

func padRight(text string, width int) string { // [AFFICHAGE] sert à compléter un texte avec des espaces jusqu'à la largeur voulue (pour aligner des colonnes)
	// On compte les lettres (runes) et pas les octets : « é » compte pour 1.
	length := utf8.RuneCountInString(text)
	if length >= width {
		return text
	}
	return text + strings.Repeat(" ", width-length)
}

func wrapText(text string) []string { // [AFFICHAGE] sert à couper un long texte en lignes de 76 caractères maximum, sans couper les mots, et renvoie les lignes
	lines := []string{}
	line := ""
	for _, word := range strings.Fields(text) {
		if line == "" {
			line = word
		} else if utf8.RuneCountInString(line+" "+word) <= 76 {
			line = line + " " + word
		} else {
			lines = append(lines, line)
			line = word
		}
	}
	return append(lines, line)
}

func paragraph(color, text string) { // [AFFICHAGE] sert à afficher un paragraphe coloré, coupé proprement en lignes
	for _, line := range wrapText(text) {
		fmt.Println("  " + color + line + Reset)
	}
}

func say(name, color, text string) { // [AFFICHAGE] sert à afficher la réplique d'un personnage : Nom : « texte »
	paragraph(color+Italic, name+" : « "+text+" »")
}

// ════════════════════════════════════════════════════════════════════════
//   LES CADRES ET LES MENUS
// ════════════════════════════════════════════════════════════════════════

func banner(text, color string) { // [AFFICHAGE] sert à afficher le grand cadre ╔══╗ qui sert de titre à un écran
	line := strings.Repeat("═", utf8.RuneCountInString(text)+6)
	fmt.Println(color + "  ╔" + line + "╗" + Reset)
	fmt.Println(color + "  ║   " + Bold + text + Reset + color + "   ║" + Reset)
	fmt.Println(color + "  ╚" + line + "╝" + Reset)
}

func section(text string) { // [AFFICHAGE] sert à afficher un petit titre de partie : ── Titre ──────
	width := 50 - utf8.RuneCountInString(text)
	if width < 2 {
		width = 2
	}
	line := strings.Repeat("─", width)
	fmt.Println("\n  " + DarkGray + "── " + Reset + Gold + Bold + text + Reset + DarkGray + " " + line + Reset)
}

func option(number int, label string) { // [AFFICHAGE] sert à afficher une ligne de menu : [1] Texte
	fmt.Printf("   %s[%d]%s %s\n", Gold+Bold, number, Reset, label)
}

func optionHint(number int, label, hint string) { // [AFFICHAGE] sert à afficher une ligne de menu avec une explication grise à droite
	option(number, padRight(label, 20)+Gray+hint+Reset)
}

func backOption(label string) { // [AFFICHAGE] sert à afficher le choix [0] (retour, quitter…), en gris
	fmt.Printf("   %s[0] %s%s\n", Gray, label, Reset)
}

func itemLine(number int, name, detail string) { // [AFFICHAGE] sert à afficher un objet dans une liste : [ 1] Nom de l'objet   détail (prix, quantité…)
	fmt.Printf("   %s[%2d]%s %-30s %s\n", Gold+Bold, number, Reset, name, detail)
}

// ════════════════════════════════════════════════════════════════════════
//   LES MESSAGES  (√ réussite · × échec · • info · ! alerte)
// ════════════════════════════════════════════════════════════════════════

// Ces 4 fonctions s'utilisent comme fmt.Printf : un texte avec des %d / %s,
// puis les valeurs à mettre dedans. « args ...any » veut dire « autant de
// valeurs qu'on veut, de n'importe quel type ».

func success(format string, args ...any) { // [AFFICHAGE] sert à afficher un message de réussite, en vert : √ …
	fmt.Printf("  "+Green+"√ "+format+Reset+"\n", args...)
}

func fail(format string, args ...any) { // [AFFICHAGE] sert à afficher un message d'échec ou d'erreur, en rouge : × …
	fmt.Printf("  "+Red+"× "+format+Reset+"\n", args...)
}

func info(format string, args ...any) { // [AFFICHAGE] sert à afficher une information neutre, en gris : • …
	fmt.Printf("  "+Silver+"• "+format+Reset+"\n", args...)
}

func warn(format string, args ...any) { // [AFFICHAGE] sert à afficher une alerte, en orange : ! …
	fmt.Printf("  "+Orange+"! "+format+Reset+"\n", args...)
}

// ════════════════════════════════════════════════════════════════════════
//   LES JAUGES  (█████░░░░░)
// ════════════════════════════════════════════════════════════════════════

func bar(current, maximum, width int, color string) string { // [AFFICHAGE] sert à fabriquer une jauge █████░░░░░ de la largeur voulue, et la renvoie
	// On borne la valeur entre 0 et le maximum pour ne jamais déborder.
	if current < 0 {
		current = 0
	}
	if current > maximum {
		current = maximum
	}
	filled := current * width / maximum
	return color + strings.Repeat("█", filled) + DarkGray + strings.Repeat("░", width-filled) + Reset
}

func hpBar(current, maximum, width int) string { // [AFFICHAGE] sert à fabriquer une jauge de PV : verte, puis jaune sous 60 %, puis rouge sous 30 %
	color := Green
	if current*100 <= maximum*60 {
		color = Yellow
	}
	if current*100 <= maximum*30 {
		color = Red
	}
	return bar(current, maximum, width, color)
}
