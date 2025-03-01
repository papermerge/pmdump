package exporter_app_v3_3

import (
	"os"

	models "github.com/papermerge/pmdump/models/app_v2_0"
	"gopkg.in/yaml.v3"
)

func CreateYAML(
	fileName string,
	u interface{},
) error {

	users := u.(models.Users)

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
