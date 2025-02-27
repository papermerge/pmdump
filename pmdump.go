package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/papermerge/pmdump/commands"
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
		commands.ListConfigs(*configFile, *targetFile)
		os.Exit(0)
	}

	if len(args) == 0 {
		fmt.Printf("Missing command: can be either %q or %q\n", exportCommand, importCommand)
		os.Exit(1)
	}

	if args[0] == exportCommand {
		commands.PerformExport(*configFile, *targetFile, exportYaml)
	} else if args[0] == importCommand {
		commands.PerformImport(*configFile, *targetFile, exportYaml)
	} else {
		fmt.Printf("Unknown command. can be either %q or %q\n", exportCommand, importCommand)
		os.Exit(1)
	}
}
