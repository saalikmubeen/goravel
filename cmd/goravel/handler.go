package main

import (
	"errors"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/iancoleman/strcase"
)

func handleHandler(arg3 string) error {
	if arg3 == "" {
		exitGracefully(errors.New("you must give the handler a name"), "Usage: goravel make handler <handler_name>")
	}

	// to file
	fileName := gor.RootPath + "/handlers/" + strings.ToLower(arg3) + ".go"

	if fileExists(fileName) {
		return errors.New(fileName + " already exists!")
	}

	data, err := templateFS.ReadFile("templates/handlers/handler.go.txt")
	if err != nil {
		exitGracefully(err)
	}

	handler := string(data)
	handler = strings.ReplaceAll(handler, "$HANDLERNAME$", strcase.ToCamel(arg3))

	err = os.WriteFile(fileName, []byte(handler), 0644)
	if err != nil {
		return err
	}

	color.Green("The following handler file has been added: %s", "/handlers/"+strings.ToLower(arg3)+".go")
	color.Yellow("Don't forget to add your handler to the routes in routes.go")

	return nil
}
