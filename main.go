package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/papermerge/pmg-dump/config"
	"github.com/papermerge/pmg-dump/database"
	"github.com/papermerge/pmg-dump/exporter"
	"github.com/papermerge/pmg-dump/models"

	_ "github.com/mattn/go-sqlite3"
)

var configFile = flag.String("c", "config.yaml", "path to config file")
var listConfigurations = flag.Bool("l", false, "List configurations and quit")
var targetFile = flag.String("f", "output.tar.gz", "Target file - zipped tar archive file name were to dump")

const export_yaml = "export.yaml"
const exportCommand = "export"
const importCommand = "import"

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Missing command")
		os.Exit(1)
		return
	}

	if args[0] == exportCommand {
		doExport()
	} else if args[0] == importCommand {
		doImport()
	} else {
		fmt.Printf("Unknown command. Can be either %q or %q\n", exportCommand, importCommand)
		os.Exit(1)
		return
	}

}

func doExport() {
	settings, err := config.ReadConfig(*configFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
		return
	}

	if *listConfigurations {
		fmt.Printf("Configuration file: %s\n", *configFile)
		fmt.Printf("Database URL: %s\n", settings.DatabaseURL)
		fmt.Printf("Media Root: %s\n", settings.MediaRoot)
		fmt.Printf("Target File: %s\n", *targetFile)
		return
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

	userIDdict := models.MakeUserID2UIDMap(users)
	nodeIDdict := models.MakeNodeID2UIDMap(nodes)

	idsDict := models.IDDict{
		UserIDs: userIDdict,
		NodeIDs: nodeIDdict,
	}

	folders, err := models.GetFolders(nodes, idsDict)

	documents, err := models.GetDocuments(nodes, settings.MediaRoot, idsDict)

	err = exporter.CreateYAML(export_yaml, users, folders, documents)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
		return
	}

	paths, err := models.GetFilePaths(documents, settings.MediaRoot)

	if err != nil {
		log.Fatalf("Error getting files paths: %v", err)
		return
	}

	paths = append(paths, models.FilePath{Source: export_yaml, Dest: export_yaml})

	err = exporter.CreateTarGz(*targetFile, paths)
	if err != nil {
		log.Fatalf("Error creating archive: %v", err)
		return
	}
	os.Remove(export_yaml)
}

func doImport() {

}
