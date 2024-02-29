package incompetenthostingprovider

import (
	"io"

	"os"
	"path/filepath"

	log "github.com/rs/zerolog/log"
)

var terraformTemplateFilenames = []string{"main.tf", "terraform.tfvars.json", "creds.tfvars.json"}

func EnsureTerraformUserDirectory(path string, userId string) error {
	// Concat path and userId
	fullPath := filepath.Join(path, userId)

	// Check if directory exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// Create directory
		err := os.Mkdir(fullPath, 0755)
		if err != nil {
			return err
		}

	}

	// Check if directory is missing any necessary terraform file
	for _, filename := range terraformTemplateFilenames {
		if _, err := os.Stat(filepath.Join(fullPath, filename)); os.IsNotExist(err) {
			log.Debug().Msgf("File %s does not exist in directory %s, copying from template directory", filename, fullPath)
			// Copy file from template directory
			err := copyFile(filepath.Join(path, "template", filename), filepath.Join(fullPath, filename))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src string, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}

	// Copy source file to destination file
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// Close destination file
	err = dstFile.Close()
	if err != nil {
		return err
	}

	return nil
}
