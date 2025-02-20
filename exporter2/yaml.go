package exporter2

import (
	"os"

	"github.com/papermerge/pmg-dump/models2"
	"gopkg.in/yaml.v3"
)

func CreateYAML(
	fileName string,
	users []models2.User,
) error {

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	data := models2.Data{
		Users: users,
	}

	yamlData, err := yaml.Marshal(&data)
	if err != nil {
		return err
	}

	_, err = file.Write(yamlData)

	return err
}
