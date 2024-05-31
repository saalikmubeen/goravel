package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func handleSession() error {
	// create migration files
	dbType := gor.DB.DatabaseType
	fileName := fmt.Sprintf("%d_create_session_tables", time.Now().UnixMicro())

	// to files:
	upFile := gor.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
	downFile := gor.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

	err := copyFilefromTemplate("templates/migrations/sessions."+dbType+".up.sql", upFile)
	if err != nil {
		return err
	}

	err = copyFilefromTemplate("templates/migrations/sessions."+dbType+".down.sql", downFile)
	if err != nil {
		return err
	}

	// run those migrations
	err = handleMigrate("up", "")
	if err != nil {
		exitGracefully(err)
	}

	color.Green("✓ Successfully created and executed the migrations for sessions.")
	color.Yellow("")
	return nil
}
