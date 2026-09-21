package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// Largeur de référence du terminal : les textes longs sont repliés dessus.
const textWidth = 76

// Direction artistique : « l'or du camp, la pierre du donjon ».
//   - l'or est la couleur signature : titres, numéros des menus, Y-Coins ;
//   - les neutres (blanc → gris foncé) portent le texte et les cadres ;
//   - chaque statistique garde sa couleur partout : ♥ rouge, ♦ bleu, ★ violet ;
//   - le camp est chaud (dégradé doré), chaque étage a son ambiance :
//     mousse (étage 1), glace (étage 2), feu (étage 3).
const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Italic = "\033[3m"

	// Neutres
	White    = "\033[38;5;255m"
	Silver   = "\033[38;5;250m"
	Gray     = "\033[38;5;244m"
	DarkGray = "\033[38;5;238m"

	// Signature
	Gold   = "\033[38;5;220m"
	Yellow = "\033[38;5;227m"
	Orange = "\033[38;5;208m"
	Brown  = "\033[38;5;130m"

	// Sens fixe : PV et danger, réussite, mana, expérience et magie
	Red     = "\033[38;5;196m"
	DarkRed = "\033[38;5;124m"
	Green   = "\033[38;5;77m"
	Lime    = "\033[38;5;118m"
	Sky     = "\033[38;5;117m"
	Blue    = "\033[38;5;33m"
	Purple  = "\033[38;5;135m"
)

// Dégradés appliqués aux dessins, du haut vers le bas. Il n'y en a qu'un par
// ambiance : un dessin prend celui du lieu où il apparaît.
var (
	campColors  = []string{White, Yellow, Gold, Gold, Orange, Brown}   // le camp et ses habitants
	stoneColors = []string{White, Silver, Silver, Gray}                // le donjon, le métal, la mort
	mossColors  = []string{Lime, Lime, Green, Green, Green, Gray}      // étage 1
	iceColors   = []string{White, Sky, Sky, Blue, Blue, Gray}          // étage 2
	fireColors  = []string{Yellow, Gold, Orange, Orange, Red, DarkRed} // étage 3 et tout ce qui brûle
	bloodColors = []string{Red, Red, DarkRed, DarkRed, Gray}           // la défaite
)

// ------------------------------------------------------------------ saisie

var scanner = bufio.NewScanner(os.Stdin)

func readLine() string {
	fmt.Print("\n" + Gold + Bold + "  ► " + Reset)
	if !scanner.Scan() {
		fmt.Println("\nÀ bientôt, aventurier !")
		os.Exit(0)
	}
	return strings.TrimSpace(scanner.Text())
}

// readChoice redemande tant que la réponse n'est pas un nombre entre low et high.
func readChoice(low, high int) int {
	for {
		number, err := strconv.Atoi(readLine())
		if err == nil && number >= low && number <= high {
			return number
		}
		fail("Hmm ? Tapez un nombre entre %d et %d.", low, high)
	}
}

// ask pose une question fermée et renvoie true si le joueur répond oui.
func ask(question string) bool {
	section(question)
	option(1, "Oui")
	option(2, "Non")
	return readChoice(1, 2) == 1
}

func pause() {
	fmt.Print("\n" + DarkGray + Italic + "  [ Entrée pour continuer ]" + Reset)
	scanner.Scan()
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func wait(milliseconds int) {
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)
}

// ----------------------------------------------------------------- dessins

// printArt affiche un dessin braille. Avec une seule couleur il est uni,
// avec plusieurs il reçoit un dégradé réparti du haut vers le bas.
func printArt(art string, colors ...string) {
	drawArt(art, colors, 0)
}

// printArtSlow fait apparaître le dessin ligne par ligne.
func printArtSlow(art string, colors ...string) {
	drawArt(art, colors, 45)
}

func drawArt(art string, colors []string, delay int) {
	if len(colors) == 0 {
		colors = []string{White}
	}
	lines := strings.Split(strings.Trim(art, "\n"), "\n")
	for i, line := range lines {
		fmt.Println(Bold + colors[i*len(colors)/len(lines)] + line + Reset)
		wait(delay)
	}
}

// sideBySide affiche plusieurs dessins en colonnes de même largeur.
func sideBySide(arts, colors []string, width int) {
	blocks := make([][]string, len(arts))
	height := 0
	for i, art := range arts {
		blocks[i] = strings.Split(strings.Trim(art, "\n"), "\n")
		height = max(height, len(blocks[i]))
	}
	for row := 0; row < height; row++ {
		for i, lines := range blocks {
			line := ""
			if row < len(lines) {
				line = lines[row]
			}
			fmt.Print(Bold + colors[i] + padRight(line, width) + Reset)
		}
		fmt.Println()
	}
}

// ------------------------------------------------------------------ textes

func padRight(text string, width int) string {
	if pad := width - utf8.RuneCountInString(text); pad > 0 {
		return text + strings.Repeat(" ", pad)
	}
	return text
}

func centerText(text string, width int) string {
	pad := width - utf8.RuneCountInString(text)
	if pad <= 0 {
		return text
	}
	return strings.Repeat(" ", pad/2) + text + strings.Repeat(" ", pad-pad/2)
}

// wrap coupe un texte en lignes d'au plus width colonnes, sans couper les mots.
func wrap(text string, width int) []string {
	lines, line := []string{}, ""
	for _, word := range strings.Fields(text) {
		switch {
		case line == "":
			line = word
		case utf8.RuneCountInString(line)+1+utf8.RuneCountInString(word) <= width:
			line += " " + word
		default:
			lines, line = append(lines, line), word
		}
	}
	return append(lines, line)
}

// typeLine écrit une ligne lettre par lettre, comme une machine à écrire.
func typeLine(color, line string) {
	fmt.Print(color)
	for _, letter := range line {
		fmt.Print(string(letter))
		wait(12)
	}
	fmt.Println(Reset)
}

// typewrite affiche un texte long replié, lettre par lettre.
func typewrite(color, text string) {
	for _, line := range wrap(text, textWidth) {
		fmt.Print("  ")
		typeLine(color, line)
	}
}

// paragraph affiche un texte long replié, sans animation.
func paragraph(color, text string) {
	for _, line := range wrap(text, textWidth) {
		fmt.Println("  " + color + line + Reset)
	}
}

// say fait parler un personnage : « Nom : « réplique » », replié et aligné.
func say(name, color, text string) {
	width := utf8.RuneCountInString(name) + 5
	indent := strings.Repeat(" ", width)
	for i, line := range wrap("« "+text+" »", textWidth-width) {
		if i == 0 {
			fmt.Print("  " + Bold + color + name + Reset + Gray + " : " + Reset)
		} else {
			fmt.Print("  " + indent)
		}
		typeLine(Italic+color, line)
	}
}

// ------------------------------------------------------------------- cadres

// banner encadre le titre d'un écran : en or au camp, à la couleur de l'étage
// dans le donjon, en rouge face à un boss.
func banner(text, color string) {
	width := utf8.RuneCountInString(text) + 6
	line := strings.Repeat("═", width)
	fmt.Println(color + "  ╔" + line + "╗" + Reset)
	fmt.Println(color + "  ║" + Bold + centerText(text, width) + Reset + color + "║" + Reset)
	fmt.Println(color + "  ╚" + line + "╝" + Reset)
}

// section sépare les parties d'un écran : « ── Titre ──────── ».
func section(text string) {
	line := strings.Repeat("─", max(2, 50-utf8.RuneCountInString(text)))
	fmt.Println("\n  " + DarkGray + "── " + Reset + Gold + Bold + text + Reset + DarkGray + " " + line + Reset)
}

// option affiche une entrée de menu « [n] libellé ».
func option(number int, label string) {
	fmt.Printf("   %s[%d]%s %s\n", Gold+Bold, number, Reset, label)
}

// optionHint affiche une entrée de menu suivie d'une courte explication.
func optionHint(number int, label, hint string) {
	if hint == "" {
		option(number, label)
		return
	}
	option(number, padRight(label, 20)+Gray+hint+Reset)
}

// back affiche l'entrée [0], toujours en dernier : retour ou sortie.
func back(label string) {
	fmt.Printf("   %s[0] %s%s\n", Gray, label, Reset)
}

// itemLine affiche une ligne de liste « [n] Objet          détail ».
func itemLine(number int, name, color, detail string) {
	fmt.Printf("   %s[%2d]%s %s%-30s%s %s\n", Gold+Bold, number, Reset, color, name, Reset, detail)
}

// ------------------------------------------------------------------ messages

func success(format string, args ...any) {
	fmt.Printf("  "+Green+Bold+"√ "+Reset+Green+format+Reset+"\n", args...)
}

func fail(format string, args ...any) {
	fmt.Printf("  "+Red+Bold+"× "+Reset+Red+format+Reset+"\n", args...)
}

func info(format string, args ...any) {
	fmt.Printf("  "+Gray+"• "+Reset+Silver+format+Reset+"\n", args...)
}

func warn(format string, args ...any) {
	fmt.Printf("  "+Orange+Bold+"! "+Reset+Orange+format+Reset+"\n", args...)
}

// dots affiche un texte suivi de trois points qui s'égrènent.
func dots(text, color string) {
	fmt.Print("  " + color + text)
	for i := 0; i < 3; i++ {
		wait(300)
		fmt.Print(".")
	}
	fmt.Println(Reset)
}

// -------------------------------------------------------------------- barres

// bar dessine une jauge remplie proportionnellement à current / maximum.
func bar(current, maximum, width int, color string) string {
	maximum = max(maximum, 1)
	current = max(0, min(current, maximum))
	filled := current * width / maximum
	return color + strings.Repeat("█", filled) + DarkGray + strings.Repeat("░", width-filled) + Reset
}

// hpBar verdit, jaunit puis rougit à mesure que les PV baissent.
func hpBar(current, maximum, width int) string {
	color := Green
	switch {
	case current*100 <= maximum*30:
		color = Red
	case current*100 <= maximum*60:
		color = Yellow
	}
	return bar(current, maximum, width, color)
}

// showHP affiche « nom ♥ [jauge] pv / max ». Le nom peut être vide.
func showHP(name string, current, maximum int) {
	if name != "" {
		name += " "
	}
	fmt.Printf("  %s%s♥%s %s %d / %d\n", name, Red, Reset, hpBar(current, maximum, 20), current, maximum)
}

// showMana affiche la jauge de mana du héros.
func showMana(current, maximum int) {
	fmt.Printf("  %s♦%s %s %d / %d\n", Blue, Reset, bar(current, maximum, 20, Blue), current, maximum)
}
