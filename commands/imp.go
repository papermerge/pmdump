package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/papermerge/pmdump/config"
	"github.com/papermerge/pmdump/database"
	"github.com/papermerge/pmdump/importer"
	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/types"
)

func PerformImport(configFile, targetFile, exportYaml string) {
	settings, err := config.ReadConfig(configFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	err = importer.ExtractTarGz(targetFile, settings.MediaRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting archive: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Documents extracted into %q\n", settings.MediaRoot)

	yamlPath := settings.MediaRoot + "/" + exportYaml
	var sourceData models.Data
	err = importer.ReadYAML(yamlPath, &sourceData)

	if err != nil {
		fmt.Printf("Error:performImport: %s", err)
	}
	db, err := database.Open(settings.DatabaseURL, types.AppVersion(settings.AppVersion))
	if err != nil {
		log.Fatal(err)
	}
	defer db.DB.Close()

	targetUsers, err := database.GetTargetUsers(db)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading target users: %v\n", err)
		os.Exit(1)
	}

	database.InsertUsersData(db, sourceData.Users, targetUsers)
}
