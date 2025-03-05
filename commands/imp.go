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

func PerformImport(settings config.Config, targetFile, exportYaml string) {

	if _, err := validateImportConfig(settings); err != nil {
		fmt.Fprintf(os.Stderr, "Validation Error: %v\n", err)
		os.Exit(1)
	}

	err := importer.ExtractTarGz(targetFile, settings.MediaRoot)

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

	result, err := database.InsertUsersData(db, sourceData.Users, targetUsers)
	updateUserIDs(sourceData, result)

	if err = database.InsertGroups(db, sourceData.Groups); err != nil {
		fmt.Fprintf(os.Stderr, "Error inserting groups: %v\n", err)
		os.Exit(1)
	}

	if err = database.InsertPermissions(db, sourceData.Permissions); err != nil {
		fmt.Fprintf(os.Stderr, "Error inserting permissions: %v\n", err)
		os.Exit(1)
	}

	if err = database.InsertGroupsPermissions(db, sourceData.GroupsPermissions); err != nil {
		fmt.Fprintf(os.Stderr, "Error inserting groups_permissions: %v\n", err)
		os.Exit(1)
	}

	if err = database.InsertDocumentTypes(db, sourceData.DocumentTypes); err != nil {
		fmt.Fprintf(os.Stderr, "Error inserting document_types: %v\n", err)
		os.Exit(1)
	}

	if err = database.InsertCustomFields(db, sourceData.CustomFields); err != nil {
		fmt.Fprintf(os.Stderr, "Error inserting custom_fields: %v\n", err)
		os.Exit(1)
	}

	if err = database.InsertCustomFieldValues(db, sourceData.CustomFieldValues); err != nil {
		fmt.Fprintf(os.Stderr, "Error inserting custom_field_values: %v\n", err)
		os.Exit(1)
	}

	if err = database.InsertDocumentTypesCustomFields(db, sourceData.DocumentTypesCustomFields); err != nil {
		fmt.Fprintf(os.Stderr, "Error inserting document_types_custom_fields: %v\n", err)
		os.Exit(1)
	}

	if err = database.InsertTags(db, sourceData.Tags); err != nil {
		fmt.Fprintf(os.Stderr, "Error inserting tags: %v\n", err)
		os.Exit(1)
	}

	if err = database.InsertNodesTags(db, sourceData.NodesTags); err != nil {
		fmt.Fprintf(os.Stderr, "Error inserting nodes_tags: %v\n", err)
		os.Exit(1)
	}

	if err = database.InsertUsersGroups(db, sourceData.UsersGroups); err != nil {
		fmt.Fprintf(os.Stderr, "Error inserting users_groups: %v\n", err)
		os.Exit(1)
	}

	if err = database.InsertUsersPermissions(db, sourceData.UsersPermissions); err != nil {
		fmt.Fprintf(os.Stderr, "Error inserting users_permissions: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Import complete.")
}

func updateUserIDs(data models.Data, userIDschange []types.UserIDChange) {
	for _, entry := range userIDschange {
		updateUserID(data, entry)
	}
}

func updateUserID(data models.Data, userIDchange types.UserIDChange) {

	for i := range len(data.DocumentTypes) {
		if data.DocumentTypes[i].UserID == userIDchange.SourceUserID {
			data.DocumentTypes[i].UserID = userIDchange.TargetUserID
		}
	}

	for i := range len(data.Tags) {
		if data.Tags[i].UserID == userIDchange.SourceUserID {
			data.Tags[i].UserID = userIDchange.TargetUserID
		}
	}

	for i := range len(data.UsersGroups) {
		if data.UsersGroups[i].UserID == userIDchange.SourceUserID {
			data.UsersGroups[i].UserID = userIDchange.TargetUserID
		}
	}

	for i := range len(data.UsersPermissions) {
		if data.UsersPermissions[i].UserID == userIDchange.SourceUserID {
			data.UsersPermissions[i].UserID = userIDchange.TargetUserID
		}
	}

	for i := range len(data.CustomFields) {
		if data.CustomFields[i].UserID == userIDchange.SourceUserID {
			data.CustomFields[i].UserID = userIDchange.TargetUserID
		}
	}
}

func validateImportConfig(settings config.Config) (bool, error) {

	if settings.AppVersion != "3.4" {
		return false, fmt.Errorf("AppVersion %q not supported", settings.AppVersion)
	}

	return true, nil
}
