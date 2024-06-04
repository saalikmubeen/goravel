package main

import (
	"embed"
	"errors"
	"os"
)

//go:embed "all:templates"
var templateFS embed.FS

func copyFilefromTemplate(from string, to string) error {
	// from -> templatePath

	if fileExists(to) {
		return errors.New(to + " already exists!")
	}

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
