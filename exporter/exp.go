package exporter

import (
	"fmt"

	exporter_app_v2_0 "github.com/papermerge/pmdump/exporter/app_v2_0"
	exporter_app_v3_3 "github.com/papermerge/pmdump/exporter/app_v3_3"
	"github.com/papermerge/pmdump/types"
)

func CreateYAML(
	fileName string,
	data any,
	appVersion types.AppVersion,
) error {
	switch appVersion {
	case types.V2_0:
		return exporter_app_v2_0.CreateYAML(fileName, data)

	case types.V3_3:
		return exporter_app_v3_3.CreateYAML(fileName, data)
	}

	return fmt.Errorf("CreateYaml: app version %q not supported", appVersion)
}
