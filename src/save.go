package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Une partie se range dans l'un des trois emplacements, sous forme d'un
// fichier JSON qui contient exactement la structure Character.
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
	if err := os.WriteFile(slotFile(c.SaveSlot), data, 0644); err != nil {
		fail("Impossible d'écrire la sauvegarde : %v", err)
		return
	}
	printArt(artScroll, Brown)
	success("Partie sauvegardée dans l'emplacement %d.", c.SaveSlot)
}

func loadGame(slot int) (*Character, error) {
	data, err := os.ReadFile(slotFile(slot))
	if err != nil {
		return nil, err
	}
	var c Character
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	if c.Inventory == nil { // sac absent d'une très vieille sauvegarde
		c.Inventory = map[string]int{}
	}
	c.SaveSlot = slot
	return &c, nil
}

// slotSummary décrit une sauvegarde en une ligne ; le booléen dit si elle existe.
func slotSummary(slot int) (string, bool) {
	c, err := loadGame(slot)
	if err != nil {
		return DarkGray + "— vide —" + Reset, false
	}
	summary := fmt.Sprintf("%s%s%s · %s niv.%d · %s · étage %d/3",
		Bold+Gold, c.Name, Reset, c.Class, c.Level, difficultyOf(c.Difficulty).Name, c.FloorsCleared)
	if c.DragonSlain {
		summary += Gold + " · ★ Dragon vaincu" + Reset
	}
	return summary, true
}

func chooseSlot(title string, mustExist bool) int {
	for {
		clearScreen()
		printArt(artScroll, Brown)
		banner(title, Gold)
		fmt.Println()
		for slot := 1; slot <= saveSlots; slot++ {
			summary, _ := slotSummary(slot)
			option(slot, fmt.Sprintf("Emplacement %d : %s", slot, summary), White)
		}
		option(0, "Retour", Gray)

		slot := readChoice(0, saveSlots)
		if slot == 0 {
			return 0
		}
		_, exists := slotSummary(slot)
		switch {
		case mustExist && !exists:
			fail("Cet emplacement est vide.")
			pause()
		case !mustExist && exists:
			warn("Cet emplacement contient déjà une partie.")
			if ask("L'écraser ?") {
				return slot
			}
		default:
			return slot
		}
	}
}
