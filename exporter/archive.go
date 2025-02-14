package exporter

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
)

func createTarGz(outputFilename string, files []string) error {
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
	for _, file := range files {
		if err := addFileToTar(tarWriter, file); err != nil {
			return err
		}
	}

	return nil
}

func addFileToTar(tw *tar.Writer, filename string) error {
	// Open the file
	file, err := os.Open(filename)
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
		Name:    filename,
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
