package exporter_app_v2_0

import (
	"fmt"
	"os"

	"github.com/papermerge/pmdump/config"
	"github.com/papermerge/pmdump/database"
	"github.com/papermerge/pmdump/exporter"
	models "github.com/papermerge/pmdump/models/app_v2_0"
	"github.com/papermerge/pmdump/types"
)

func PerformExport(
	settings config.Config,
	targetFile,
	exportYaml string,
) []types.FilePath {
	var filePaths []types.FilePath

	db, err := database.Open(settings.DatabaseURL, types.AppVersion(settings.AppVersion))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer db.DB.Close()

	results, err := database.GetUsers(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: GetUsers: %v\n", err)
		os.Exit(1)
	}

	users := results.([]models.User)

	for i := 0; i < len(users); i++ {
		database.GetUserNodes(db, &users[i])
		results, err := database.GetDocumentPageRows(db, users[i].LegacyID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting GetDocumentPageRows: %v", err)
		}

		docPages := results.([]models.DocumentPageRow)

		models.ForEachDocument(
			users[i].Home,
			users[i].LegacyID,
			docPages,
			settings.MediaRoot,
			models.InsertDocVersionsAndPages,
		)
		models.ForEachDocument(
			users[i].Inbox,
			users[i].LegacyID,
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
		userFilePaths, err := models.GetFilePaths(allDocs, users[i].LegacyID, settings.MediaRoot)

		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting file paths: %v\n", err)
		}

		filePaths = append(filePaths, userFilePaths...)
	}

	payload := models.Data{
		Users: users,
	}

	err = exporter.CreateYAML(
		exportYaml,
		payload,
		types.V2_0,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file:performExport2: %v", err)
		os.Exit(1)
	}

	return filePaths
}
