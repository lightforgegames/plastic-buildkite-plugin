package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

var rx = regexp.MustCompile(`.private.\d$`)
var suffixesToCheck []string

func init() {
	suffixesToCheck = make([]string, 0, 12)
	for i := 0; i < 12; i += 1 {
		suffixesToCheck = append(suffixesToCheck, fmt.Sprintf(".private.%d", i))
	}
}

func findPrivateFiles(pathToCheck string) []string {
	filesToDelete := make([]string, 0)
	matchPrivateFile := func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if rx.MatchString(path) {
			relPath, err := filepath.Rel(pathToCheck, path)
			if err != nil {
				return nil
			}
			filesToDelete = append(filesToDelete, relPath)
		}
		return nil
	}
	_ = filepath.WalkDir(pathToCheck, matchPrivateFile)
	return filesToDelete
}

func DeletePrivateFiles(pathToCheck string) []string {
	filesToDelete := findPrivateFiles(pathToCheck)
	for _, f := range filesToDelete {
		os.Remove(f)
	}

	return filesToDelete
}
