package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/papermerge/pmg-dump/config"
	"github.com/papermerge/pmg-dump/database"
	"github.com/papermerge/pmg-dump/database2"
	"github.com/papermerge/pmg-dump/exporter"
	"github.com/papermerge/pmg-dump/exporter2"
	"github.com/papermerge/pmg-dump/importer"
	"github.com/papermerge/pmg-dump/models"
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
		performImport()
	} else {
		fmt.Printf("Unknown command. can be either %q or %q\n", exportCommand, importCommand)
		os.Exit(1)
	}
}

func performExport2() {
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

	for i, _ := range users {
		database2.GetUserNodes(db, &users[i])
		docPages, err := database2.GetDocumentPageRows(db, users[i].ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting GetDocumentPageRows: %v", err)
		}
		docs1 := users[i].Home.GetUserDocuments()
		docs2 := users[i].Inbox.GetUserDocuments()
		docs := append(docs1, docs2...)
		err = models2.InsertDocVersionsAndPages(docs, docPages, settings.MediaRoot)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while inserting document versions: %v", err)
		}
	}

	err = exporter2.CreateYAML(
		exportYaml,
		users,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file:performExport2: %v", err)
		os.Exit(1)
	}
}

func performExport() {
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

	users, err := database.GetUsers(db)
	if err != nil {
		log.Fatal(err)
	}

	nodes, err := database.GetNodes(db)
	if err != nil {
		log.Fatal(err)
	}

	tags, err := database.GetTags(db)
	if err != nil {
		log.Fatal(err)
	}

	docPages, err := database.GetDocumentPageRows(db)

	userIDdict := models.MakeUserID2UIDMap(users)
	nodeIDdict := models.MakeNodeID2UIDMap(nodes)

	idsDict := models.IDDict{
		UserIDs: userIDdict,
		NodeIDs: nodeIDdict,
	}

	folders, err := models.GetFolders(nodes, idsDict)

	documents, err := models.GetDocuments(nodes, settings.MediaRoot, idsDict, docPages)

	err = exporter.CreateYAML(
		exportYaml,
		users,
		folders,
		documents,
		tags,
	)

	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
		return
	}

	paths, err := models.GetFilePaths(documents, settings.MediaRoot)

	if err != nil {
		log.Fatalf("Error getting files paths: %v", err)
		return
	}

	paths = append(paths, models.FilePath{Source: exportYaml, Dest: exportYaml})

	err = exporter.CreateTarGz(*targetFile, paths)
	if err != nil {
		log.Fatalf("Error creating archive: %v", err)
		return
	}
	os.Remove(exportYaml)
}

func performImport() {
	settings, err := config.ReadConfig(*configFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	err = importer.ExtractTarGz(*targetFile, settings.MediaRoot)
	if err != nil {
		log.Fatalf("Error extracting archive: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Documents extracted into %q\n", settings.MediaRoot)
	yamlPath := settings.MediaRoot + "/" + exportYaml
	var data models.Data
	err = importer.ReadYAML(yamlPath, &data)

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
