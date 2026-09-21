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

// Codes couleur ANSI : le terminal les interprète au lieu de les afficher.
const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Italic = "\033[3m"

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

// Dégradés appliqués aux dessins, du haut vers le bas : 6 nuances d'une même
// couleur, du plus clair au plus foncé. Les nombres sont des couleurs du
// terminal (de 0 à 255).
var (
	campColors  = shades(229, 228, 222, 220, 214, 172) // doré
	stoneColors = shades(255, 253, 251, 249, 247, 245) // gris pierre
	mossColors  = shades(190, 154, 118, 82, 40, 34)    // vert
	iceColors   = shades(195, 159, 123, 81, 75, 33)    // bleu glacé
	fireColors  = shades(227, 220, 214, 208, 202, 196) // feu
)

// shades transforme des numéros de couleur en codes ANSI.
func shades(codes ...int) []string {
	colors := []string{}
	for _, code := range codes {
		colors = append(colors, fmt.Sprintf("\033[38;5;%dm", code))
	}
	return colors
}

var scanner = bufio.NewScanner(os.Stdin)

func readLine() string {
	fmt.Print("\n" + Gold + Bold + "  ► " + Reset)
	if !scanner.Scan() {
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
		fail("Tapez un nombre entre %d et %d.", low, high)
	}
}

func ask(question string) bool {
	section(question)
	option(1, "Oui")
	option(2, "Non")
	return readChoice(1, 2) == 1
}

func pause() {
	fmt.Print("\n" + DarkGray + "  [ Entrée pour continuer ]" + Reset)
	if !scanner.Scan() {
		os.Exit(0)
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func wait(milliseconds int) {
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)
}

// printArt affiche un dessin. Avec plusieurs couleurs, elles forment un
// dégradé du haut vers le bas.
func printArt(art string, colors ...string) {
	if len(colors) == 0 {
		colors = []string{Silver}
	}
	lines := strings.Split(strings.Trim(art, "\n"), "\n")
	for i, line := range lines {
		color := colors[i*len(colors)/len(lines)]
		fmt.Println(Bold + color + line + Reset)
	}
}

func padRight(text string, width int) string {
	length := utf8.RuneCountInString(text)
	if length >= width {
		return text
	}
	return text + strings.Repeat(" ", width-length)
}

// wrap coupe un texte en lignes de 76 caractères maximum, sans couper les mots.
func wrap(text string) []string {
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

func paragraph(color, text string) {
	for _, line := range wrap(text) {
		fmt.Println("  " + color + line + Reset)
	}
}

func say(name, color, text string) {
	paragraph(color+Italic, name+" : « "+text+" »")
}

func banner(text, color string) {
	line := strings.Repeat("═", utf8.RuneCountInString(text)+6)
	fmt.Println(color + "  ╔" + line + "╗" + Reset)
	fmt.Println(color + "  ║   " + Bold + text + Reset + color + "   ║" + Reset)
	fmt.Println(color + "  ╚" + line + "╝" + Reset)
}

func section(text string) {
	width := 50 - utf8.RuneCountInString(text)
	if width < 2 {
		width = 2
	}
	line := strings.Repeat("─", width)
	fmt.Println("\n  " + DarkGray + "── " + Reset + Gold + Bold + text + Reset + DarkGray + " " + line + Reset)
}

func option(number int, label string) {
	fmt.Printf("   %s[%d]%s %s\n", Gold+Bold, number, Reset, label)
}

func optionHint(number int, label, hint string) {
	option(number, padRight(label, 20)+Gray+hint+Reset)
}

func back(label string) {
	fmt.Printf("   %s[0] %s%s\n", Gray, label, Reset)
}

func itemLine(number int, name, detail string) {
	fmt.Printf("   %s[%2d]%s %-30s %s\n", Gold+Bold, number, Reset, name, detail)
}

func success(format string, args ...any) {
	fmt.Printf("  "+Green+"√ "+format+Reset+"\n", args...)
}

func fail(format string, args ...any) {
	fmt.Printf("  "+Red+"× "+format+Reset+"\n", args...)
}

func info(format string, args ...any) {
	fmt.Printf("  "+Silver+"• "+format+Reset+"\n", args...)
}

func warn(format string, args ...any) {
	fmt.Printf("  "+Orange+"! "+format+Reset+"\n", args...)
}

func bar(current, maximum, width int, color string) string {
	if current < 0 {
		current = 0
	}
	if current > maximum {
		current = maximum
	}
	filled := current * width / maximum
	return color + strings.Repeat("█", filled) + DarkGray + strings.Repeat("░", width-filled) + Reset
}

// hpBar est une jauge de PV : verte, puis jaune sous 60 %, puis rouge sous 30 %.
func hpBar(current, maximum, width int) string {
	color := Green
	if current*100 <= maximum*60 {
		color = Yellow
	}
	if current*100 <= maximum*30 {
		color = Red
	}
	return bar(current, maximum, width, color)
}
