package filemanagement

import (
	"fmt"
	"os"
	"path/filepath"

	"duplyzer/shared"
)

func deleteFile(path string) error {
	return os.Remove(path)
}

func moveFile(srcPath, destDir string) error {
	destPath := filepath.Join(destDir, filepath.Base(srcPath))
	return os.Rename(srcPath, destPath)
}

func hardLinkFile(srcPath, destDir string) error {
	destPath := filepath.Join(destDir, filepath.Base(srcPath))
	return os.Link(srcPath, destPath)
}

func ManageDuplicates(action, targetDir string, duplicates shared.Results) error {
	for _, files := range duplicates {
		for i := 1; i < len(files); i++ {
			file := files[i]
			var err error
			switch action {
			case "delete":
				err = deleteFile(file)
			case "move":
				err = moveFile(file, targetDir)
			case "hard-link":
				err = hardLinkFile(file, targetDir)
			default:
				return fmt.Errorf("unsupported action: %s", action)
			}
			if err != nil {
				return fmt.Errorf("error managing file %s: %v", file, err)
			}
		}
	}
	return nil
}
