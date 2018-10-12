package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

// ReadPathsFile read the paths file from cache directory
func ReadPathsFile(cacheDir string) ([]string, error) {
	pathsFilename := filepath.Join(cacheDir, "paths")
	file, err := os.Open(pathsFilename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// WritePathsFile write the paths file to cache directory
func WritePathsFile(cacheDir string, paths []string) error {
	pathsFilename := filepath.Join(cacheDir, "paths")
	file, err := os.OpenFile(pathsFilename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range paths {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
