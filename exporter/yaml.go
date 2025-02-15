package exporter

import (
	"os"

	"github.com/papermerge/pmg-dump/models"
	"gopkg.in/yaml.v3"
)

func CreateYAML(fileName string, users []models.User, folders []models.Folder) error {

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	data := models.Data{
		Users:   users,
		Folders: folders,
	}

	yamlData, err := yaml.Marshal(&data)
	if err != nil {
		return err
	}

	_, err = file.Write(yamlData)

	return err
}
