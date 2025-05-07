package main

import (
	"fmt"
	"os"

	"github.com/blang/semver/v4"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

const (
	currentVersion = "0.1.0"                              // The current version hardcoded in your app
	repoSlug       = "robertgouveia/rkw-software-support" // GitHub repository slug
)

func main() {
	// Start the self-update process
	updateSoftware()
}

func updateSoftware() {
	// Parse the current version using github.com/blang/semver
	v, err := semver.Parse(currentVersion)
	if err != nil {
		fmt.Printf("Error parsing current version: %s\n", err.Error())
		return
	}

	// Create the updater
	updater, err := selfupdate.NewUpdater(selfupdate.Config{})
	if err != nil {
		fmt.Printf("Error creating updater: %s\n", err.Error())
		return
	}

	// Detect the latest release on GitHub
	latest, found, err := updater.DetectLatest(repoSlug)
	if err != nil {
		fmt.Printf("Error checking for updates: %s\n", err.Error())
		return
	}

	// Check if we found the latest version or it returned a nil
	if !found || latest == nil {
		fmt.Println("No latest version found or failed to fetch the version.")
		return
	}

	// Ensure that latest.Version is comparable to v (both should be *semver.Version)
	latestVersion, err := semver.Parse(latest.Version.String())
	if err != nil {
		fmt.Printf("Error parsing latest version: %s\n", err.Error())
		return
	}

	// If the latest version is not found or is not newer, exit
	if latestVersion.LTE(v) {
		fmt.Println("You're up-to-date!")
		return
	}

	// Display the new version found and start updating
	fmt.Printf("New version found: %s\nUpdating...\n", latestVersion)
	if err := updater.UpdateTo(latest, os.Args[0]); err != nil {
		fmt.Printf("Update failed: %s\n", err.Error())
		return
	}

	// Success message
	fmt.Printf("Successfully updated to version: %s\n", latestVersion)
}
