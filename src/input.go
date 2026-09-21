// ════════════════════════════════════════════════════════════════════════
//   input.go · [SAISIE]
//   Tout ce que le joueur tape au clavier passe par ici. readChoice
//   redemande tant que la réponse n'est pas valide : taper n'importe quoi
//   ne peut jamais faire planter le jeu.
// ════════════════════════════════════════════════════════════════════════

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Le scanner lit le clavier (os.Stdin) ligne par ligne. Il est créé une
// seule fois et partagé par tout le jeu.
var scanner = bufio.NewScanner(os.Stdin) // [SAISIE] sert à lire ce que le joueur tape, une ligne à la fois

func readLine() string { // [SAISIE] sert à afficher le curseur ► et à lire une ligne tapée par le joueur, sans les espaces autour
	fmt.Print("\n" + Gold + Bold + "  ► " + Reset)
	// Scan renvoie false quand il n'y a plus rien à lire (Ctrl+D, fin du
	// fichier) : on quitte le jeu proprement au lieu de boucler à l'infini.
	if !scanner.Scan() {
		os.Exit(0)
	}
	return strings.TrimSpace(scanner.Text())
}

func readChoice(low, high int) int { // [SAISIE] sert à lire un numéro entre low et high (redemande tant que la réponse est invalide) et le renvoie
	for {
		// strconv.Atoi transforme le texte "3" en nombre 3. Si le texte
		// n'est pas un nombre ("abc"), err n'est pas vide (nil).
		number, err := strconv.Atoi(readLine())
		if err == nil && number >= low && number <= high {
			return number
		}
		fail("Tapez un nombre entre %d et %d.", low, high)
	}
}

func askYesNo(question string) bool { // [SAISIE] sert à poser une question Oui / Non et renvoie true si le joueur répond Oui
	section(question)
	option(1, "Oui")
	option(2, "Non")
	return readChoice(1, 2) == 1
}

func pause() { // [SAISIE] sert à attendre que le joueur appuie sur Entrée (pour qu'il ait le temps de lire)
	fmt.Print("\n" + DarkGray + "  [ Entrée pour continuer ]" + Reset)
	if !scanner.Scan() {
		os.Exit(0)
	}
}
