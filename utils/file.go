package utils

import "os"

func IsReadableFile(path string) bool {
	// Check if file exists
	info, err := os.Stat(path)
	if err != nil {
		return false // File does not exist or other error
	}

	// Check if it's a regular file (not a directory)
	if !info.Mode().IsRegular() {
		return false
	}

	// Try opening the file to check readability
	file, err := os.Open(path)
	if err != nil {
		return false // Cannot open file (permission denied, etc.)
	}
	file.Close() // Close after checking

	return true // File exists and is readable
}
