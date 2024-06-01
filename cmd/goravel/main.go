package main

import (
	"errors"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/saalikmubeen/goravel"
)

const version = "1.1.0"

var gor *goravel.Goravel

func main() {

	arg1, arg2, arg3, err := readArgs()
	if err != nil {
		exitGracefully(err)
	}

	// ** Initialize Goravel
	err = initGoravel(arg1)
	if err != nil {
		exitGracefully(err)
	}

	switch arg1 {
	case "help":
		showHelp()

	case "version":
		color.Green("Version: %s", gor.Version)

	case "new":
		if arg2 == "" {
			exitGracefully(errors.New("new requires a project name"), "Usage: goravel new <project_name>")
		} else {
			err := createNewGoravelApp(arg2)
			if err != nil {
				exitGracefully(err)
			}

		}

	case "make":
		if arg2 == "" {
			exitGracefully(errors.New("make requires a subcommand: (migration|model|controller)"), "Usage: goravel make <command>")
		} else {
			err := handleMake(arg2, arg3)
			if err != nil {
				exitGracefully(err)
			}
		}

	case "migrate":
		if arg2 == "" {
			var dsn = buildDSN()
			gor.MigrateUp(dsn)
		} else {
			err := handleMigrate(arg2, arg3)
			if err != nil {
				exitGracefully(err)
			}
		}
	}

}

func initGoravel(arg1 string) error {

	gor = &goravel.Goravel{}

	rootPath, err := os.Getwd()
	if err != nil {
		return err
	}

	gor.RootPath = rootPath
	gor.Version = version

	if arg1 == "help" || arg1 == "version" || arg1 == "new" {
		return nil
	}

	// Load the .env file into the environment
	// read .env
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	// ** create the loggers
	infoLog, errorLog := gor.CreateLoggers()
	gor.InfoLog = infoLog
	gor.ErrorLog = errorLog

	gor.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))

	var databaseType = ""

	if os.Getenv("DATABASE_TYPE") == "postgres" || os.Getenv("DATABASE_TYPE") == "postgresql" {
		databaseType = "postgres"
	} else {
		databaseType = os.Getenv("DATABASE_TYPE")
	}

	gor.DB = goravel.Database{
		DatabaseType: databaseType,
	}

	return nil
}

func readArgs() (string, string, string, error) {
	var arg1, arg2, arg3 string

	if len(os.Args) > 1 { // goravel <arg1>
		arg1 = os.Args[1] // len(os.Args) >= 2

		if len(os.Args) >= 3 { // goravel <arg1> <arg2>
			arg2 = os.Args[2]
		}

		if len(os.Args) >= 4 { // goravel <arg1> <arg2> <arg3>
			arg3 = os.Args[3]
		}
	} else {
		color.Red("Error: command required")
		showHelp()
		return "", "", "", errors.New("command required")
	}

	return arg1, arg2, arg3, nil
}

func exitGracefully(err error, msg ...string) {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	}

	if err != nil {
		color.Red("Error: %v\n", err)
	}

	if len(message) > 0 {
		color.Yellow(message)
	} else {
		color.Green("Finished!")
	}

	os.Exit(0)
}
