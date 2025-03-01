package exporter_app_v2_0

import (
	"fmt"
	"log"
	"os"

	"github.com/papermerge/pmdump/config"
	"github.com/papermerge/pmdump/database"
	"github.com/papermerge/pmdump/exporter"
	models "github.com/papermerge/pmdump/models/app_v2_0"
	"github.com/papermerge/pmdump/types"
)

func PerformExport(settings config.Config, targetFile, exportYaml string) {
	var filePaths []types.FilePath

	db, err := database.Open(settings.DatabaseURL, types.AppVersion(settings.AppVersion))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer db.DB.Close()

	results, err := database.GetUsers(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	users := results.(models.Users)

	for i := 0; i < len(users); i++ {
		database.GetUserNodes(db, &users[i])
		results, err := database.GetDocumentPageRows(db, users[i].ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting GetDocumentPageRows: %v", err)
		}

		docPages := results.(models.DocumentPageRows)

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
		types.V2_0,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file:performExport2: %v", err)
		os.Exit(1)
	}

	filePaths = append(filePaths, types.FilePath{Source: exportYaml, Dest: exportYaml})

	err = exporter.CreateTarGz(targetFile, filePaths)
	if err != nil {
		log.Fatalf("Error creating archive: %v", err)
		return
	}
	os.Remove(exportYaml)
}
