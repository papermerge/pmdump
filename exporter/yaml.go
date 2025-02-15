package exporter

import (
	"os"

	"github.com/papermerge/pmg-dump/models"
	"gopkg.in/yaml.v3"
)

func CreateYAML(fileName string, users []models.User, nodes []models.Node) error {

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	var folders []models.Folder
	var documents []models.Document

	for _, node := range nodes {
		if node.Model == models.DocumentModelName {
			document := models.Document{
				ID:        node.ID,
				Title:     node.Title,
				UserID:    node.UserID,
				ParentID:  node.ParentID,
				Version:   node.Version,
				FileName:  node.FileName,
				PageCount: node.PageCount,
				UUID:      node.UUID,
			}
			documents = append(documents, document)
		} else {
			folder := models.Folder{
				ID:       node.ID,
				Title:    node.Title,
				UserID:   node.UserID,
				ParentID: node.ParentID,
				UUID:     node.UUID,
			}
			folders = append(folders, folder)
		}
	}

	data := models.Data{
		Users:     users,
		Documents: documents,
		Folders:   folders,
	}

	yamlData, err := yaml.Marshal(&data)
	if err != nil {
		return err
	}

	_, err = file.Write(yamlData)

	return err
}
