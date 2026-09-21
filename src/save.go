// ════════════════════════════════════════════════════════════════════════
//   save.go · [SAUVEGARDE]
//   Une sauvegarde est un fichier texte au format JSON qui contient la
//   structure Character (le héros), et rien d'autre. Il y a 3 emplacements :
//   dungo_save_1.json, dungo_save_2.json et dungo_save_3.json, écrits dans
//   le dossier depuis lequel on lance le jeu.
// ════════════════════════════════════════════════════════════════════════

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const saveSlots = 3 // [SAUVEGARDE] sert à fixer le nombre d'emplacements de sauvegarde

func saveFileName(slot int) string { // [SAUVEGARDE] sert à donner le nom du fichier d'un emplacement (2 → "dungo_save_2.json") et le renvoie
	return fmt.Sprintf("dungo_save_%d.json", slot)
}

func saveExists(slot int) bool { // [SAUVEGARDE] sert à savoir si un fichier de sauvegarde existe pour cet emplacement ; renvoie true si oui
	_, err := os.Stat(saveFileName(slot)) // os.Stat lit les infos d'un fichier : err n'est pas nil s'il n'existe pas
	return err == nil
}

func saveGame(character *Character) { // [SAUVEGARDE] sert à écrire le héros dans son fichier JSON (et prévient si ça échoue)
	// 1. Character → texte JSON (MarshalIndent = avec des retours à la ligne, lisible)
	data, err := json.MarshalIndent(character, "", "  ")
	if err != nil {
		fail("Impossible de sauvegarder : %v", err)
		return
	}
	// 2. Texte JSON → fichier (0644 = lisible par tous, modifiable par soi)
	err = os.WriteFile(saveFileName(character.SaveSlot), data, 0644)
	if err != nil {
		fail("Impossible d'écrire la sauvegarde : %v", err)
		return
	}
	success("Partie sauvegardée (emplacement %d).", character.SaveSlot)
}

func loadGame(slot int) (*Character, error) { // [SAUVEGARDE] sert à relire un fichier de sauvegarde ; renvoie le héros, ou une erreur si le fichier manque ou est illisible
	// 1. Fichier → texte JSON
	data, err := os.ReadFile(saveFileName(slot))
	if err != nil {
		return nil, err
	}
	// 2. Texte JSON → Character (les champs absents du fichier restent à 0 / vides)
	var character Character
	err = json.Unmarshal(data, &character)
	if err != nil {
		return nil, err
	}
	// 3. Sécurité : une map absente du fichier vaut nil, et ranger un objet
	//    dans une map nil fait planter Go. On les crée vides au besoin.
	if character.Inventory == nil {
		character.Inventory = map[string]int{}
	}
	if character.Equipment == nil {
		character.Equipment = map[string]string{}
	}
	character.SaveSlot = slot
	return &character, nil
}

func slotDescription(slot int) string { // [SAUVEGARDE] sert à décrire un emplacement pour le menu (« Emplacement vide », « Sauvegarde illisible » ou le héros sauvegardé) et renvoie le texte
	if !saveExists(slot) {
		return "Emplacement vide"
	}
	saved, err := loadGame(slot)
	if err != nil {
		return Red + "Sauvegarde illisible" + Reset
	}
	return fmt.Sprintf("%s%s%s · %s · niveau %d · étage %d/%d", Gold+Bold, saved.Name, Reset, saved.Class, saved.Level, saved.FloorsCleared, len(floors))
}

func chooseSlot(title string, mustExist bool) int { // [SAUVEGARDE] sert à afficher les 3 emplacements et renvoie celui choisi (0 = retour) ; mustExist = true pour « Continuer »
	for {
		clearScreen()
		printArt(artScroll, campColors...)
		banner(title, Gold)
		fmt.Println()
		for slot := 1; slot <= saveSlots; slot++ {
			option(slot, slotDescription(slot))
		}
		backOption("Retour")

		slot := readChoice(0, saveSlots)
		if slot == 0 {
			return 0
		}
		exists := saveExists(slot)

		// « Continuer » sur un emplacement vide : impossible
		if mustExist && !exists {
			fail("Cet emplacement est vide.")
			pause()
			continue
		}
		// « Nouvelle partie » sur un emplacement occupé : on demande avant d'écraser
		if !mustExist && exists {
			if !askYesNo("Cet emplacement contient déjà une partie. L'écraser ?") {
				continue
			}
		}
		return slot
	}
}
