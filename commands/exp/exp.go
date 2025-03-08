package exporter

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	exporter_app_v2_0 "github.com/papermerge/pmdump/commands/exp/app_v2_0"
	exporter_app_v3_2 "github.com/papermerge/pmdump/commands/exp/app_v3_2"
	exporter_app_v3_3 "github.com/papermerge/pmdump/commands/exp/app_v3_3"
	"github.com/papermerge/pmdump/config"
	"github.com/papermerge/pmdump/exporter"
	"github.com/papermerge/pmdump/types"
)

func PerformExport(settings config.Config, targetFile, exportYaml string) {
	var filePaths []types.FilePath

	if _, err := validExportConfig(settings); err != nil {
		fmt.Fprintf(os.Stderr, "Validation Error: %v\n", err)
		os.Exit(1)
	}

	switch settings.AppVersion {
	case string(types.V2_0):
		filePaths = exporter_app_v2_0.PerformExport(settings, targetFile, exportYaml)
	case string(types.V3_3):
		filePaths = exporter_app_v3_3.PerformExport(settings, targetFile, exportYaml)
	case string(types.V3_2):
		filePaths = exporter_app_v3_2.PerformExport(settings, targetFile, exportYaml)
	default:
		supported_versions := []types.AppVersion{
			types.V2_0, types.V2_1,
		}
		fmt.Fprintf(os.Stderr, "Export for version %q not supported\n", settings.AppVersion)
		fmt.Fprintf(os.Stderr, "Supported versions are %v\n", supported_versions)
		os.Exit(1)
	}

	filePaths = append(filePaths, types.FilePath{Source: exportYaml, Dest: exportYaml})

	err := exporter.CreateTarGz(targetFile, filePaths)
	if err != nil {
		log.Fatalf("Error creating archive: %v", err)
		return
	}
	os.Remove(exportYaml)
	fmt.Println("Export complete")
}

func validExportConfig(settings config.Config) (bool, error) {

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
