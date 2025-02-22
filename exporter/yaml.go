package exporter

import (
	"os"

	"github.com/papermerge/pmdump/models"
	"gopkg.in/yaml.v3"
)

func CreateYAML(
	fileName string,
	users []models.User,
) error {

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	data := models.Data{
		Users: users,
	}

	yamlData, err := yaml.Marshal(&data)
	if err != nil {
		return err
	}

	_, err = file.Write(yamlData)

	return err
}
