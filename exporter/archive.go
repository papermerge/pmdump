package exporter

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"github.com/papermerge/pmg-dump/models"
)

func GetFolders(nodes []models.Node) ([]models.Folder, error) {
	var folders []models.Folder

	for _, node := range nodes {
		if node.Model == models.FolderModelName {

			folder := models.Folder{
				ID:     node.ID,
				Title:  node.Title,
				UserID: node.UserID,
				UUID:   node.UUID,
			}

			folder.UserUUID = UserID2UUID[folder.UserID]

			if node.ParentID != nil {
				parentUUID := NodeID2UUID[*node.ParentID]
				folder.ParentUUID = &parentUUID
				folder.ParentID = node.ParentID
			}
			folders = append(folders, folder)
		}
	}

	return folders, nil
}

func GetFilePaths(users []models.User, nodes []models.Node, mediaRoot string) ([]models.FilePath, error) {
	var paths []models.FilePath

	for _, user := range users {
		for _, node := range nodes {
			if node.Model == models.DocumentModelName && node.FileName != nil {
				source := fmt.Sprintf(
					"%s/docs/user_%d/document_%d/%s",
					mediaRoot,
					user.ID,
					node.ID,
					*node.FileName,
				)

				uid := node.UUID.String()
				dest := fmt.Sprintf("docvers/%s/%s/%s", uid[0:2], uid[2:4], *node.FileName)
				path := models.FilePath{
					Source: source,
					Dest:   dest,
				}
				paths = append(paths, path)
			}
		}
	}

	return paths, nil
}

func CreateTarGz(outputFilename string, paths []models.FilePath) error {
	// Create output file
	outFile, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Create gzip writer
	gzipWriter := gzip.NewWriter(outFile)
	defer gzipWriter.Close()

	// Create tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Add files to archive
	for _, path := range paths {
		if err := addFileToTar(tarWriter, path); err != nil {
			return err
		}
	}

	return nil
}

func addFileToTar(tw *tar.Writer, path models.FilePath) error {
	// Open the file
	file, err := os.Open(path.Source)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file info
	info, err := file.Stat()
	if err != nil {
		return err
	}

	// Create tar header
	header := &tar.Header{
		Name:    path.Dest,
		Size:    info.Size(),
		Mode:    int64(info.Mode()),
		ModTime: info.ModTime(),
	}

	// Write header to tar
	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	// Copy file data to tar
	if _, err := io.Copy(tw, file); err != nil {
		return err
	}

	return nil
}
