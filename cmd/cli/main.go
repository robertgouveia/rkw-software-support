package main

import (
	"log"

	"github.com/robertgouveia/do-my-job/menu"
	"github.com/robertgouveia/do-my-job/tea"
)

func main() {
	mainMenu := tea.Create("Server Configuration Tool")
	menu.ScriptMenu(mainMenu)
	menu.ServerMenu(mainMenu)

	_, err := mainMenu.Run()
	if err != nil {
		log.Fatalf("Error running menu: %v", err)
	}
}
