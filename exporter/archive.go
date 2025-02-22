package exporter

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"github.com/papermerge/pmdump/models"
)

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
			fmt.Fprintf(os.Stderr, "Warning:CreateTarGz:skipping %q\n", path.Source)
			continue
		}
	}

	return nil
}

func addFileToTar(tw *tar.Writer, path models.FilePath) error {
	// Open the file
	file, err := os.Open(path.Source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error:addFileToTar: could not open source %q (with destination %q)\n", path.Source, path.Dest)
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
