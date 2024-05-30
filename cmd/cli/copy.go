package main

import (
	"embed"
	"os"
)

//go:embed templates
var templateFS embed.FS

func copyFilefromTemplate(from string, to string) error {
	// from -> templatePath

	data, err := templateFS.ReadFile(from)

	if err != nil {
		return err
	}

	err = copyDataToFile(to, data)

	if err != nil {
		return nil
	}

	return nil
}

func copyDataToFile(to string, data []byte) error {
	err := os.WriteFile(to, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func fileExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}
