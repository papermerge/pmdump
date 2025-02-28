package commands

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

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
		fmt.Fprintf(os.Stderr, "Error: %v\n", configFile, err)
		os.Exit(1)
	}

	if _, err := validExportConfig(settings); err != nil {
		fmt.Fprintf(os.Stderr, "Validation Error: %v\n", err)
		os.Exit(1)
	}

	db, err := database.Open(settings.DatabaseURL, types.AppVersion(settings.AppVersion))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer db.DB.Close()

	users, err := database.GetUsers(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
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

func validExportConfig(settings *config.Config) (bool, error) {

	if !contains(types.AppVersionsForExport, types.AppVersion(settings.AppVersion)) {
		return false, fmt.Errorf("AppVersion %q not supported", settings.AppVersion)
	}

	parsedDBURL, err := url.Parse(settings.DatabaseURL)

	if err != nil {
		return false, fmt.Errorf("Error parsing dburl %s: %v", settings.DatabaseURL, err)
	}

	if strings.HasPrefix(parsedDBURL.Scheme, "postgres") && settings.AppVersion == string(types.V2_0) {
		errMsg := "Export of Papermerge DMS v2.0 and Postgres database is not implemented.\n" +
			"In case you need this feature please open a ticket https://github.com/ciur/papermerge\n" +
			"In the ticket post docker compose file with your configurations."
		return false, fmt.Errorf(errMsg)
	}

	if strings.HasPrefix(parsedDBURL.Scheme, "maria") && settings.AppVersion == string(types.V2_0) {
		errMsg := "Export of Papermerge DMS v2.0 and MariaDB database is not implemented.\n" +
			"In case you need this feature please open a ticket https://github.com/ciur/papermerge\n" +
			"In the ticket post docker compose file with your configurations."
		return false, fmt.Errorf(errMsg)
	}

	if strings.HasPrefix(parsedDBURL.Scheme, "my") && settings.AppVersion == string(types.V2_0) {
		errMsg := "Export of Papermerge DMS v2.0 and MySQL database is not implemented.\n" +
			"In case you need this feature please open a ticket https://github.com/ciur/papermerge\n" +
			"In the ticket post docker compose file with your configurations."
		return false, fmt.Errorf(errMsg)
	}

	return true, nil
}

func contains(slice []types.AppVersion, value types.AppVersion) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
