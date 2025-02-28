package commands

import (
	"fmt"
	"os"

	"github.com/papermerge/pmdump/config"
)

func ListConfigs(configFile, targetFile string) {
	settings, err := config.ReadConfig(configFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Configuration file: %s\n", configFile)
	fmt.Printf("Database URL: %s\n", settings.DatabaseURL)
	fmt.Printf("Media Root: %s\n", settings.MediaRoot)
	fmt.Printf("Target File: %s\n", targetFile)
	fmt.Printf("App Version: %s\n", settings.AppVersion)
}
