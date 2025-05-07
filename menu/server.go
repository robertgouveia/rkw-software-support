package menu

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/robertgouveia/do-my-job/database"
	"github.com/robertgouveia/do-my-job/lib"
	"github.com/robertgouveia/do-my-job/storage"
	"github.com/robertgouveia/do-my-job/tea"
)

func ServerMenu(mainMenu *tea.TeaModel, update func() string) *tea.TeaModel {
	configureServerMenu := tea.Create("Configure Servers")
	mainMenu.AddSubmenu("Configure Servers", configureServerMenu)

	configureServerMenu.AddSubmenu("RKW Data Warehouse", serverTemplate("RKW Data Warehouse", "RKW Data Warehouse"))
	configureServerMenu.AddSubmenu("NAVSLAT02", serverTemplate("NAVSQLAT02", "NAVSQLAT02"))
	configureServerMenu.AddSubmenu("RKW Level 1", serverTemplate("RKW Level 1", "RKW Level 1"))
	configureServerMenu.AddSubmenu("RKW Level 3 VIC", serverTemplate("RKW Level 3 VIC", "RKW Level 3 VIC"))
	configureServerMenu.AddSubmenu("RKW Level 3 STONE", serverTemplate("RKW Level 3 STONE", "RKW Level 3 STONE"))

	mainMenu.AddMenuItem("Update", func() string {
		return update()
	})

	mainMenu.AddMenuItem("About", func() string {
		return `
RKW Support Tool v0.1
=============================
A simple tool to execute automations to make your life easier
Built with Golang
Robert Gouveia

Configurations are stored in: ` + storage.GetConfigDir()
	})

	return configureServerMenu
}

func serverTemplate(serverName, title string) *tea.TeaModel {
	config, err := storage.LoadServerConfig(serverName)
	if err != nil {
		log.Printf("Warning: Could not load saved config for %s: %v", serverName, err)
	}

	rkwServerMenu := tea.Create(title)

	rkwServerMenu.AddTextInput(
		"Set Host",
		"Enter server hostname or IP address:",
		fmt.Sprintf("The hostname or IP address of the RKW data warehouse server\nCurrent value: %s",
			lib.StringOrDefault(config.Host, "[Not Set]")),
		func(input string) {
			config.Host = input
			config.LastUpdated = time.Now()

			// Save after update
			if err := storage.SaveServerConfig(serverName, config); err != nil {
				log.Printf("Failed to save config: %v", err)
			} else {
				fmt.Printf("Host set to: %s and saved\n", input)
			}
		},
	)

	rkwServerMenu.AddTextInput(
		"Set Port",
		"Enter server port:",
		fmt.Sprintf("The port number to connect to (e.g., 5432 for PostgreSQL)\nCurrent value: %s",
			lib.StringOrDefault(config.Port, "[Not Set]")),
		func(input string) {
			config.Port = input
			config.LastUpdated = time.Now()

			if err := storage.SaveServerConfig(serverName, config); err != nil {
				log.Printf("Failed to save config: %v", err)
			} else {
				fmt.Printf("Port set to: %s and saved\n", input)
			}
		},
	)

	rkwServerMenu.AddTextInput(
		"Set Username",
		"Enter username:",
		fmt.Sprintf("Database username for authentication\nCurrent value: %s",
			lib.StringOrDefault(config.Username, "[Not Set]")),
		func(input string) {
			config.Username = input
			config.LastUpdated = time.Now()

			if err := storage.SaveServerConfig(serverName, config); err != nil {
				log.Printf("Failed to save config: %v", err)
			} else {
				fmt.Printf("Username set to: %s and saved\n", input)
			}
		},
	)

	rkwServerMenu.AddTextInput(
		"Set Password",
		"Enter password:",
		fmt.Sprintf("Database password for authentication\nCurrent value: %s",
			lib.StringOrDefault(config.Password, "[Not Set]")),
		func(input string) {
			config.Password = input
			config.LastUpdated = time.Now()

			if err := storage.SaveServerConfig(serverName, config); err != nil {
				log.Printf("Failed to save config: %v", err)
			} else {
				fmt.Printf("Password set to: %s and saved\n", input)
			}
		},
	)

	rkwServerMenu.AddTextInput(
		"Set Database",
		"Enter database name:",
		fmt.Sprintf("The name of the database to connect to\nCurrent value: %s",
			lib.StringOrDefault(config.Database, "[Not Set]")),
		func(input string) {
			config.Database = input
			config.LastUpdated = time.Now()

			if err := storage.SaveServerConfig(serverName, config); err != nil {
				log.Printf("Failed to save config: %v", err)
			} else {
				fmt.Printf("Database set to: %s and saved\n", input)
			}
		},
	)

	rkwServerMenu.AddMenuItem("Test Connection", func() string {
		if config.Host == "" || config.Port == "" || config.Username == "" || config.Database == "" {
			return "Error: Please configure all server settings first."
		}

		_, err := database.Connect(serverName)
		if err != nil {
			return fmt.Sprintf("Error: %s", err.Error())
		}

		return fmt.Sprintf(`
Connection Test Results:
=======================
Host: %s
Port: %s
Username: %s
Database: %s
Status: Simulated connection successful
		`, config.Host, config.Port, config.Username, config.Database)
	})

	rkwServerMenu.AddMenuItem("View Configuration", func() string {
		lastUpdated := "Never"
		if !config.LastUpdated.IsZero() {
			lastUpdated = config.LastUpdated.Format(time.RFC1123)
		}

		return fmt.Sprintf(`
Current Configuration:
=====================
Host: %s
Port: %s
Username: %s
Password: %s
Database: %s
Last Updated: %s
Config File: %s
			`,
			lib.StringOrDefault(config.Host, "[Not Set]"),
			lib.StringOrDefault(config.Port, "[Not Set]"),
			lib.StringOrDefault(config.Username, "[Not Set]"),
			lib.StringOrDefault(config.Password, "[Not Set]"),
			lib.StringOrDefault(config.Database, "[Not Set]"),
			lastUpdated,
			filepath.Join(storage.GetConfigDir(), serverName+".json"),
		)
	})

	rkwServerMenu.AddMenuItem("Delete Saved Configuration", func() string {
		configPath := filepath.Join(storage.GetConfigDir(), serverName+".json")

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			return "No saved configuration found."
		}

		err := os.Remove(configPath)
		if err != nil {
			return fmt.Sprintf("Error deleting configuration: %v", err)
		}

		config = storage.ServerConfig{}

		return "Configuration successfully deleted."
	})

	return rkwServerMenu
}
