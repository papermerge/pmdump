package importer2

import (
	"os"

	"github.com/papermerge/pmg-dump/models2"
	"gopkg.in/yaml.v3"
)

func ReadYAML(
	fileName string,
	data *models2.Data,
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
