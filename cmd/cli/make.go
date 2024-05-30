package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/fatih/color"
)

func handleMake(arg2 string, arg3 string) error {

	switch arg2 {
	case "migration":
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
				exitGracefully(err)
			}

			err = copyFilefromTemplate("templates/migrations/migration."+dbType+".down.sql", downFile)
			if err != nil {
				exitGracefully(err)
			}

			color.Green("Migration files created successfully!")
			color.Yellow(`The following migration files have been added:
			 1. %s
			 2. %s
			 `, "/migrations/"+fileName+"."+dbType+".up.sql", "/migrations/"+fileName+"."+dbType+".down.sql")
		}
	}

	return nil
}
