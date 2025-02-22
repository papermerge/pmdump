package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/papermerge/pmg-dump/config"
	"github.com/papermerge/pmg-dump/database2"
	"github.com/papermerge/pmg-dump/exporter2"
	"github.com/papermerge/pmg-dump/importer2"
	"github.com/papermerge/pmg-dump/models2"

	_ "github.com/mattn/go-sqlite3"
)

var configFile = flag.String("c", "source.yaml", "path to config file")
var listConfigurations = flag.Bool("l", false, "List configurations and quit")
var targetFile = flag.String("f", "output.tar.gz", "Target file - zipped tar archive file name were to dump")

const exportYaml = "export.yaml"
const exportCommand = "export"
const importCommand = "import"

func main() {
	flag.Parse()

	args := flag.Args()

	if *listConfigurations {
		listConfigs()
		os.Exit(0)
	}

	if len(args) == 0 {
		fmt.Printf("Missing command: can be either %q or %q\n", exportCommand, importCommand)
		os.Exit(1)
	}

	if args[0] == exportCommand {
		performExport2()
	} else if args[0] == importCommand {
		performImport2()
	} else {
		fmt.Printf("Unknown command. can be either %q or %q\n", exportCommand, importCommand)
		os.Exit(1)
	}
}

func performExport2() {
	var filePaths []models2.FilePath
	settings, err := config.ReadConfig(*configFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	db, err := sql.Open("sqlite3", settings.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	users, err := database2.GetUsers(db)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(users); i++ {
		database2.GetUserNodes(db, &users[i])
		docPages, err := database2.GetDocumentPageRows(db, users[i].ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting GetDocumentPageRows: %v", err)
		}
		models2.ForEachDocument(
			users[i].Home,
			users[i].ID,
			docPages,
			settings.MediaRoot,
			models2.InsertDocVersionsAndPages,
		)
		models2.ForEachDocument(
			users[i].Inbox,
			users[i].ID,
			docPages,
			settings.MediaRoot,
			models2.InsertDocVersionsAndPages,
		)
	}

	for i := 0; i < len(users); i++ {
		var allDocs []models2.Node

		inbox := users[i].Inbox.GetUserDocuments()
		home := users[i].Home.GetUserDocuments()
		allDocs = append(allDocs, inbox...)
		allDocs = append(allDocs, home...)
		userFilePaths, err := models2.GetFilePaths(allDocs, users[i].ID, settings.MediaRoot)

		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting file paths: %v\n", err)
		}

		filePaths = append(filePaths, userFilePaths...)
	}

	err = exporter2.CreateYAML(
		exportYaml,
		users,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file:performExport2: %v", err)
		os.Exit(1)
	}

	filePaths = append(filePaths, models2.FilePath{Source: exportYaml, Dest: exportYaml})

	err = exporter2.CreateTarGz(*targetFile, filePaths)
	if err != nil {
		log.Fatalf("Error creating archive: %v", err)
		return
	}
	os.Remove(exportYaml)
}

func performImport2() {
	settings, err := config.ReadConfig(*configFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	err = importer2.ExtractTarGz(*targetFile, settings.MediaRoot)
	if err != nil {
		log.Fatalf("Error extracting archive: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Documents extracted into %q\n", settings.MediaRoot)

	yamlPath := settings.MediaRoot + "/" + exportYaml
	var data models2.Data
	err = importer2.ReadYAML(yamlPath, &data)

	if err != nil {
		fmt.Printf("Error:performImport: %s", err)
	}
}

func listConfigs() {
	settings, err := config.ReadConfig(*configFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Configuration file: %s\n", *configFile)
	fmt.Printf("Database URL: %s\n", settings.DatabaseURL)
	fmt.Printf("Media Root: %s\n", settings.MediaRoot)
	fmt.Printf("Target File: %s\n", *targetFile)
}
