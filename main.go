package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/papermerge/pmg-dump/config"
	"github.com/papermerge/pmg-dump/exporter"

	_ "github.com/mattn/go-sqlite3"
)

var configFile = flag.String("c", "config.yaml", "path to config file")

func main() {
	flag.Parse()

	settings, err := config.ReadConfig(*configFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
		return
	}

	fmt.Printf("Configuration file: %s\n", *configFile)
	fmt.Printf("Database URL: %s\n", settings.DatabaseURL)
	fmt.Printf("Media Root: %s\n", settings.MediaRoot)
	fmt.Printf("Target File: %s\n", settings.TargetFile)

	db, err := sql.Open("sqlite3", settings.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	users, err := exporter.GetUsers(db)
	if err != nil {
		log.Fatal(err)
	}

	nodes, err := exporter.GetNodes(db)
	if err != nil {
		log.Fatal(err)
	}

	err = exporter.CreateYAML("export.yaml", users, nodes)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
		return
	}

	paths, err := exporter.GetFilePaths(users, nodes, settings.MediaRoot)

	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}

	err = exporter.CreateTarGz(settings.TargetFile, paths)
	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}

	fmt.Printf("Success!\n")
}
