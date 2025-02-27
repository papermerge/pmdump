package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/papermerge/pmdump/config"
	"github.com/papermerge/pmdump/database"
	"github.com/papermerge/pmdump/exporter"
	"github.com/papermerge/pmdump/models"
	"github.com/papermerge/pmdump/types"
)

func PerformExport(configFile, targetFile, exportYaml string) {
	var filePaths []models.FilePath
	settings, err := config.ReadConfig(configFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	db, err := database.Open(settings.DatabaseURL, types.AppVersion(settings.AppVersion))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening dburl %s: %v", settings.DatabaseURL, err)
		os.Exit(1)
	}
	defer db.DB.Close()

	users, err := database.GetUsers(db)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(users); i++ {
		database.GetUserNodes(db, &users[i])
		docPages, err := database.GetDocumentPageRows(db, users[i].ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting GetDocumentPageRows: %v", err)
		}
		models.ForEachDocument(
			users[i].Home,
			users[i].ID,
			docPages,
			settings.MediaRoot,
			models.InsertDocVersionsAndPages,
		)
		models.ForEachDocument(
			users[i].Inbox,
			users[i].ID,
			docPages,
			settings.MediaRoot,
			models.InsertDocVersionsAndPages,
		)
	}

	for i := 0; i < len(users); i++ {
		var allDocs []models.Node

		inbox := users[i].Inbox.GetUserDocuments()
		home := users[i].Home.GetUserDocuments()
		allDocs = append(allDocs, inbox...)
		allDocs = append(allDocs, home...)
		userFilePaths, err := models.GetFilePaths(allDocs, users[i].ID, settings.MediaRoot)

		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting file paths: %v\n", err)
		}

		filePaths = append(filePaths, userFilePaths...)
	}

	err = exporter.CreateYAML(
		exportYaml,
		users,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file:performExport2: %v", err)
		os.Exit(1)
	}

	filePaths = append(filePaths, models.FilePath{Source: exportYaml, Dest: exportYaml})

	err = exporter.CreateTarGz(targetFile, filePaths)
	if err != nil {
		log.Fatalf("Error creating archive: %v", err)
		return
	}
	os.Remove(exportYaml)
}
