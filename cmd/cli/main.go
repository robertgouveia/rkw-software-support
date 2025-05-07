package main

import (
	"fmt"
	"log"
	"os"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/robertgouveia/do-my-job/menu"
	"github.com/robertgouveia/do-my-job/tea"
)

const (
	currentVersion = "0.9.9"                              // The current version hardcoded in your app
	repoSlug       = "robertgouveia/rkw-software-support" // GitHub repository slug
)

func main() {
	mainMenu := tea.Create("Server Configuration Tool")
	menu.ScriptMenu(mainMenu)
	menu.ServerMenu(mainMenu, updateSoftware)

	_, err := mainMenu.Run()
	if err != nil {
		log.Fatalf("Error running menu: %v", err)
	}
}

func updateSoftware() string {
	// Parse the current version using github.com/blang/semver
	v, err := semver.Parse(currentVersion)
	if err != nil {
		return fmt.Sprintf("Error parsing current version: %s\n", err.Error())
	}

	// Create the updater
	updater, err := selfupdate.NewUpdater(selfupdate.Config{})
	if err != nil {
		return fmt.Sprintf("Error creating updater: %s\n", err.Error())
	}

	// Detect the latest release on GitHub
	latest, found, err := updater.DetectLatest(repoSlug)
	if err != nil {
		return fmt.Sprintf("Error checking for updates: %s\n", err.Error())
	}

	// Check if we found the latest version or it returned a nil
	if !found || latest == nil {
		return fmt.Sprintf("No latest version found or failed to fetch the version.")
	}

	// Ensure that latest.Version is comparable to v (both should be *semver.Version)
	latestVersion, err := semver.Parse(latest.Version.String())
	if err != nil {
		return fmt.Sprintf("Error parsing latest version: %s\n", err.Error())
	}

	// If the latest version is not found or is not newer, exit
	if latestVersion.LTE(v) {
		return fmt.Sprint("You're up-to-date!")
	}

	// Display the new version found and start updating
	fmt.Printf("New version found: %s\nUpdating...\n", latestVersion)
	if _, err := updater.UpdateSelf(v, repoSlug); err != nil {
		fmt.Printf("Update failed: %s\n", err.Error())
		os.Exit(0)
	}

	// Success message
	return fmt.Sprintf("Successfully updated to version: %s\n", latestVersion)
}
