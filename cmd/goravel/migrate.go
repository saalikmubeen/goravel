package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/color"
)

func handleMigrate(arg2, arg3 string) error {
	dsn := buildDSN()

	// run the migration command
	switch arg2 {
	case "up":
		err := gor.MigrateUp(dsn)
		if err != nil {
			return err
		}

	case "down":
		if arg3 == "all" {
			err := gor.MigrateDownAll(dsn)
			if err != nil {
				return err
			}
		} else {
			err := gor.MigrateSteps(-1, dsn)
			if err != nil {
				return err
			}
		}

	case "reset":
		err := gor.MigrateDownAll(dsn)
		if err != nil {
			return err
		}
		err = gor.MigrateUp(dsn)
		if err != nil {
			return err
		}

	case "fix":
		err := gor.MigrateForce(dsn)
		if err != nil {
			return err
		}

	case "to":
		if arg3 == "" {
			exitGracefully(errors.New("migrate to requires a number"), "Usage: goravel migrate to <number> [number -> positive for up, negative for down]")
		} else {
			num, err := strconv.Atoi(arg3)
			if err != nil {
				return errors.New("goravel migrate to <number> | migration number must be an integer")
			}
			err = gor.MigrateSteps(num, dsn)
			if err != nil {
				return err
			}
		}

	default:
		showHelp()
	}
	return nil
}

func handleMakeMigration(arg3 string) error {
	dbType := gor.DB.DatabaseType
	if arg3 == "" {
		exitGracefully(errors.New("you must give the migration a name"), "goravel make migation <miration_name>")
	} else {
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixMicro(), arg3)

		// to files:
		upFile := gor.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
		downFile := gor.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

		err := copyFilefromTemplate("templates/migrations/migration."+dbType+".up.sql", upFile)
		if err != nil {
			return err
		}

		err = copyFilefromTemplate("templates/migrations/migration."+dbType+".down.sql", downFile)
		if err != nil {
			return err
		}

		color.Green("Migration files created successfully!")
		color.Yellow(`The following migration files have been added:
			 1. %s
			 2. %s
			 `, "/migrations/"+fileName+"."+dbType+".up.sql", "/migrations/"+fileName+"."+dbType+".down.sql")
	}

	return nil
}
