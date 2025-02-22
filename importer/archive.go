package importer

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

// extractTarGz extracts a .tar.gz archive into a specific directory
func ExtractTarGz(filename, destination string) error {
	// Open the .tar.gz file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a gzip reader
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// Create a tar reader
	tarReader := tar.NewReader(gzipReader)

	// Ensure the destination folder exists
	if err := os.MkdirAll(destination, os.ModePerm); err != nil {
		return err
	}

	// Iterate over files in the archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			// End of archive
			break
		}
		if err != nil {
			return err
		}

		// Get the full file path
		targetPath := filepath.Join(destination, header.Name)

		switch header.Typeflag {
		case tar.TypeDir: // Handle directories
			if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
				return err
			}
		case tar.TypeReg: // Handle regular files
			if err := extractFile(tarReader, targetPath, header.Mode); err != nil {
				return err
			}
		}
	}

	return nil
}

// extractFile extracts a single file from the tar archive
func extractFile(reader io.Reader, filePath string, fileMode int64) error {
	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// Create the file
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Copy content from the archive to the file
	if _, err := io.Copy(outFile, reader); err != nil {
		return err
	}

	// Set file permissions
	return os.Chmod(filePath, os.FileMode(fileMode))
}
