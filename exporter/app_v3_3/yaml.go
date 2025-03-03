package exporter_app_v3_3

import (
	"os"

	models "github.com/papermerge/pmdump/models/app_v3_3"
	"gopkg.in/yaml.v3"
)

func CreateYAML(
	fileName string,
	payload any,
) error {

	data := payload.(models.Data)

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	yamlData, err := yaml.Marshal(&data)
	if err != nil {
		return err
	}

	_, err = file.Write(yamlData)

	return err
}
