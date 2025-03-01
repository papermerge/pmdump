package exporter

import (
	"fmt"

	exporter_app_v2_0 "github.com/papermerge/pmdump/exporter/app_v2_0"
	"github.com/papermerge/pmdump/types"
)

func CreateYAML(
	fileName string,
	users interface{},
	appVersion types.AppVersion,
) error {
	switch appVersion {
	case types.V2_0:
		return exporter_app_v2_0.CreateYAML(fileName, users)

	case types.V3_3:
		return exporter_app_v2_0.CreateYAML(fileName, users)
	}

	return fmt.Errorf("CreateYaml: app version %q not supported", appVersion)
}
