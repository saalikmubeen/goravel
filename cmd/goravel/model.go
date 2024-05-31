package main

import (
	"errors"
	"strings"

	"github.com/fatih/color"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

func handleModel(arg3 string) error {
	if arg3 == "" {
		exitGracefully(errors.New("you must give your model a name"), "Usage: goravel make model <model_name>")
	}

	modelName := arg3
	tableName := arg3

	plural := pluralize.NewClient()

	if plural.IsPlural(arg3) {
		tableName = strings.ToLower(arg3)
		modelName = plural.Singular(arg3)
	} else {
		tableName = strings.ToLower(plural.Plural(arg3))
		modelName = arg3
	}

	// to file
	fileName := gor.RootPath + "/models/" + strings.ToLower(modelName) + ".go"

	if fileExists(fileName) {
		return errors.New(fileName + " already exists!")
	}

	data, err := templateFS.ReadFile("templates/models/model.go.txt")
	if err != nil {
		exitGracefully(err)
	}

	model := string(data)

	model = strings.ReplaceAll(model, "$MODELNAME$", strcase.ToCamel(modelName))
	model = strings.ReplaceAll(model, "$TABLENAME$", tableName)

	err = copyDataToFile(fileName, []byte(model))
	if err != nil {
		return err
	}

	color.Green("The following model file has been added: %s", "/models/"+strings.ToLower(modelName)+".go")
	color.Yellow("Don't forget to register this model in the models/models.go file.")

	return nil
}
