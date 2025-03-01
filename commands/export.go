package commands

import (
	exporter "github.com/papermerge/pmdump/commands/exp"
	"github.com/papermerge/pmdump/config"
)

func PerformExport(settings config.Config, targetFile, exportYaml string) {
	exporter.PerformExport(settings, targetFile, exportYaml)
}
