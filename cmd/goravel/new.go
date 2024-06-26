package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/saalikmubeen/goravel"
)

var appURL string

func createNewGoravelApp(arg2 string) error {

	appName := strings.ToLower(arg2)
	appURL = appName

	// sanitize the application name (convert url to single word)
	// github.com/saalikmubeen/new-goravel-app -> new-goravel-app
	if strings.Contains(appName, "/") {
		stringSlice := strings.SplitAfter(appName, "/")
		appName = stringSlice[(len(stringSlice) - 1)]
	}

	log.Println("App name is", appName)

	// make a new directory
	color.Green("\tCreating new directory...")
	err := os.Mkdir(appName, 0755)
	if err != nil {
		return err
	}

	// change directory
	err = os.Chdir(appName)
	if err != nil {
		return err
	}

	// set the rootPath to the new created directory
	rootPath, err := os.Getwd()
	if err != nil {
		return err
	}
	// fmt.Println("Root path is", rootPath)
	gor.RootPath = rootPath

	paths := goravel.InitPaths{
		RootPath:    gor.RootPath,
		FolderNames: []string{"handlers", "migrations", "views", "mail", "models", "public", "tmp", "logs", "middleware"},
	}

	// ** create the necessary folders
	err = gor.CreateFolderStructure(paths)
	if err != nil {
		return err
	}

	// ** Check if the .env file exists, if not create it
	// to file
	toFile := gor.RootPath + "/" + ".env"
	err = handleCopyDataToFile("templates/new/env.txt", toFile, ReplaceDataMap{
		"${APP_NAME}":   appName,
		"${GO_APP_URL}": appURL,
		"${KEY}":        gor.RandomString(32),
	})
	if err != nil {
		return err
	}

	// Copy main.go file
	toFile = gor.RootPath + "/" + "main.go"
	err = handleCopyDataToFile("templates/new/main.go.txt", toFile, ReplaceDataMap{
		"${APP_URL}": appURL,
	})
	if err != nil {
		return err
	}

	// Copy routes.go file
	toFile = gor.RootPath + "/" + "routes.go"
	err = handleCopyDataToFile("templates/new/routes.go.txt", toFile, ReplaceDataMap{})
	if err != nil {
		return err
	}

	// Copy routes_api.go file
	toFile = gor.RootPath + "/" + "routes_api.go"
	err = handleCopyDataToFile("templates/new/routes_api.go.txt", toFile, ReplaceDataMap{})
	if err != nil {
		return err
	}

	// Copy init-goravel.go file
	toFile = gor.RootPath + "/" + "init-goravel.go"
	err = handleCopyDataToFile("templates/new/init-goravel.go.txt", toFile, ReplaceDataMap{
		"${APP_URL}": appURL,
	})
	if err != nil {
		return err
	}

	// Copy middlewares.go file
	toFile = gor.RootPath + "/middleware" + "/" + "middleware.go"
	err = handleCopyDataToFile("templates/middleware/middleware.go.txt", toFile, ReplaceDataMap{
		"${APP_URL}": appURL,
	})
	if err != nil {
		return err
	}

	// Copy models.go file
	toFile = gor.RootPath + "/models" + "/" + "models.go"
	err = handleCopyDataToFile("templates/models/models.go.txt", toFile, ReplaceDataMap{
		"${APP_URL}": appURL,
	})
	if err != nil {
		return err
	}

	// Copy handlers.go file
	toFile = gor.RootPath + "/handlers" + "/" + "handlers.go"
	err = handleCopyDataToFile("templates/handlers/handlers.go.txt", toFile, ReplaceDataMap{
		"${APP_URL}": appURL,
	})
	if err != nil {
		return err
	}

	// copy the base layout file
	err = gor.CreateDirIfNotExists(gor.RootPath + "/views/layouts") // Create the layouts folder
	if err != nil {
		return err
	}
	toFile = gor.RootPath + "/views" + "/layouts/" + "base.jet"
	err = handleCopyDataToFile("templates/views/layouts/base.jet", toFile, ReplaceDataMap{})
	if err != nil {
		return err
	}

	// copy the home view file
	toFile = gor.RootPath + "/views" + "/" + "home.jet"
	err = handleCopyDataToFile("templates/views/home.jet", toFile, ReplaceDataMap{})
	if err != nil {
		return err
	}

	// Copy go.mod file
	toFile = gor.RootPath + "/" + "go.mod"
	err = handleCopyDataToFile("templates/new/go.mod.txt", toFile, ReplaceDataMap{
		"${APP_URL}": appURL,
	})
	if err != nil {
		return err
	}

	// Copy .gitignore file
	toFile = gor.RootPath + "/" + ".gitignore"
	err = handleCopyDataToFile("templates/new/gitignore.txt", toFile, ReplaceDataMap{})
	if err != nil {
		return err
	}

	color.Green("✓ Successfully created new Goravel project: %s", arg2)
	color.Yellow("")

	// run go mod tidy in the project directory
	color.Yellow("\tRunning go mod tidy...")
	cmd := exec.Command("go", "mod", "tidy")
	err = cmd.Start()
	if err != nil {
		return err
	}

	color.Green("Done building " + appURL)
	color.Green("Go build something awesome")

	return nil
}

type ReplaceDataMap map[string]string

// handleCopyDataToFile copies data from a file in the templateFS, replaces the placeholders
// and writes it to a new file in the rootPath
func handleCopyDataToFile(from string, to string, replace ReplaceDataMap) error {

	if fileExists(to) {
		return errors.New(to + " already exists!")
	}

	data, err := templateFS.ReadFile(from)
	if err != nil {
		return err
	}

	processedFile := string(data)

	for key, value := range replace {
		processedFile = strings.ReplaceAll(processedFile, key, value)
	}

	err = os.WriteFile(to, []byte(processedFile), 0644)
	if err != nil {
		return err
	}

	return nil
}
