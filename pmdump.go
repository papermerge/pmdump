package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/papermerge/pmdump/commands"
	"github.com/papermerge/pmdump/config"
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

	settings, err := config.ReadConfig(*configFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", configFile, err)
		os.Exit(1)
	}

	if *listConfigurations {
		commands.ListConfigs(*configFile, *targetFile)
		os.Exit(0)
	}

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Missing command: can be either %q or %q\n", exportCommand, importCommand)
		os.Exit(1)
	}

	if args[0] == exportCommand {
		commands.PerformExport(*settings, *targetFile, exportYaml)
	} else if args[0] == importCommand {
		commands.PerformImport(*settings, *targetFile, exportYaml)
	} else {
		fmt.Fprintf(
			os.Stderr,
			"Unknown command %q. Can be either %q or %q\n",
			args[0],
			exportCommand,
			importCommand,
		)
		os.Exit(1)
	}
}
