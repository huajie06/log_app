package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Println(dir)

	res, err := FindLogDirectory(dir)
	if err != nil {
		panic(err)
	}
	// fmt.Println(res)
	// /home/hj/apps/log_app/journal/src/event-form.html
	fmt.Println(filepath.Join(res, "./journal/src/event-form.html"))

}

// "log_app" directory by traversing up the directory tree
func FindLogDirectory(startPath string) (string, error) {
	maxLevels := 10
	targetFolder := "log_app"

	// Get absolute path
	absPath, err := filepath.Abs(startPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	currentPath := absPath
	levelsUp := 0

	for levelsUp < maxLevels {
		// List directory contents
		entries, err := os.ReadDir(currentPath)
		if err != nil {
			return "", fmt.Errorf("failed to read directory %s: %w", currentPath, err)
		}

		// Check current directory for target folder
		for _, entry := range entries {
			if entry.IsDir() && strings.EqualFold(entry.Name(), targetFolder) {
				return filepath.Join(currentPath, entry.Name()), nil
			}
		}

		// Move up one directory level
		parentPath := filepath.Dir(currentPath)

		// Check if we've reached the root directory
		if parentPath == currentPath {
			return "", fmt.Errorf("reached root directory without finding %s", targetFolder)
		}

		currentPath = parentPath
		levelsUp++
	}

	return "", fmt.Errorf("directory %s not found within %d levels up from %s",
		targetFolder, maxLevels, startPath)
}
