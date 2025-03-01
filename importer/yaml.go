package importer

import (
	"os"

	models "github.com/papermerge/pmdump/models/app_v3_3"
	"gopkg.in/yaml.v3"
)

func ReadYAML(
	fileName string,
	data *models.Data,
) error {

	content, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, data)
	if err != nil {
		return err
	}

	return err
}
