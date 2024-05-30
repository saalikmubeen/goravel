package main

import (
	"errors"
	"strconv"
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
