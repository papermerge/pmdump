package exporter_app_v3_3

import (
	"fmt"
	"os"

	"github.com/papermerge/pmdump/config"
	"github.com/papermerge/pmdump/database"
	"github.com/papermerge/pmdump/exporter"
	models "github.com/papermerge/pmdump/models/app_v3_3"
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
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	users := results.(models.Users)

	results, err = database.GetGroups(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	groups := results.([]models.Group)

	results, err = database.GetPermissions(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	perms := results.([]models.Permission)

	results, err = database.GetGroupsPermissions(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	groupsPermissions := results.([]models.GroupsPermissions)

	results, err = database.GetDocumentTypes(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting entries from 'document_types' table: %v\n", err)
		os.Exit(1)
	}

	documentTypes := results.([]models.DocumentType)

	results, err = database.GetTags(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting entries from 'tags' table: %v\n", err)
		os.Exit(1)
	}

	tags := results.([]models.Tag)

	results, err = database.GetNodesTags(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting entries from 'nodes_tags' table: %v\n", err)
		os.Exit(1)
	}

	nodesTags := results.([]models.NodesTags)

	for i := 0; i < len(users); i++ {

		database.GetUserNodes(db, &users[i])
		models.ForEachDocument(
			db,
			users[i].Home,
			database.InsertDocVersionsAndPages,
		)
		models.ForEachDocument(
			db,
			users[i].Inbox,
			database.InsertDocVersionsAndPages,
		)
	}

	for i := 0; i < len(users); i++ {
		var allDocs []models.Node

		inbox := users[i].Inbox.GetUserDocuments()
		home := users[i].Home.GetUserDocuments()
		allDocs = append(allDocs, inbox...)
		allDocs = append(allDocs, home...)
		userFilePaths, err := models.GetFilePaths(allDocs, settings.MediaRoot)

		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting file paths: %v\n", err)
		}

		filePaths = append(filePaths, userFilePaths...)
	}

	payload := models.Data{
		Users:             users,
		Groups:            groups,
		Permissions:       perms,
		GroupsPermissions: groupsPermissions,
		DocumentTypes:     documentTypes,
		Tags:              tags,
		NodesTags:         nodesTags,
	}

	err = exporter.CreateYAML(
		exportYaml,
		payload,
		types.V3_3,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file: %v", err)
		os.Exit(1)
	}

	return filePaths
}
