package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Une sauvegarde est un fichier JSON qui contient la structure Character.
const saveSlots = 3

func slotFile(slot int) string {
	return fmt.Sprintf("dungo_save_%d.json", slot)
}

func saveGame(c *Character) {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		fail("Impossible de sauvegarder : %v", err)
		return
	}
	err = os.WriteFile(slotFile(c.SaveSlot), data, 0644)
	if err != nil {
		fail("Impossible d'écrire la sauvegarde : %v", err)
		return
	}
	success("Partie sauvegardée (emplacement %d).", c.SaveSlot)
}

func loadGame(slot int) (*Character, error) {
	data, err := os.ReadFile(slotFile(slot))
	if err != nil {
		return nil, err
	}
	var c Character
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	// Une vieille sauvegarde peut ne pas avoir ces champs : sans ces maps
	// vides, le jeu planterait au premier objet rangé.
	if c.Inventory == nil {
		c.Inventory = map[string]int{}
	}
	if c.Equipment == nil {
		c.Equipment = map[string]string{}
	}
	c.SaveSlot = slot
	return &c, nil
}

// chooseSlot affiche les 3 emplacements et renvoie celui choisi (0 = retour).
func chooseSlot(title string, mustExist bool) int {
	for {
		clearScreen()
		printArt(artScroll, campColors...)
		banner(title, Gold)
		fmt.Println()
		for slot := 1; slot <= saveSlots; slot++ {
			c, err := loadGame(slot)
			if err != nil {
				option(slot, "Emplacement vide")
			} else {
				option(slot, fmt.Sprintf("%s%s%s · %s · niveau %d · étage %d/3", Gold+Bold, c.Name, Reset, c.Class, c.Level, c.FloorsCleared))
			}
		}
		back("Retour")

		slot := readChoice(0, saveSlots)
		if slot == 0 {
			return 0
		}
		_, err := loadGame(slot)
		exists := err == nil

		if mustExist && !exists {
			fail("Cet emplacement est vide.")
			pause()
			continue
		}
		if !mustExist && exists {
			if !ask("Cet emplacement contient déjà une partie. L'écraser ?") {
				continue
			}
		}
		return slot
	}
}
