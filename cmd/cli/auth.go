package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func handleAuth() error {
	// create migration files
	dbType := gor.DB.DatabaseType
	fileName := fmt.Sprintf("%d_create_auth_tables", time.Now().UnixMicro())

	// to files:
	upFile := gor.RootPath + "/migrations/" + fileName + ".up.sql"
	downFile := gor.RootPath + "/migrations/" + fileName + ".down.sql"

	err := copyFilefromTemplate("templates/migrations/auth_tables."+dbType+".up.sql", upFile)
	if err != nil {
		return err
	}

	err = copyFilefromTemplate("templates/migrations/auth_tables."+dbType+".down.sql", downFile)
	if err != nil {
		return err
	}

	// run those migrations
	err = handleMigrate("up", "")
	if err != nil {
		exitGracefully(err)
	}

	// Copy the auth models and authentication code
	err = copyFilefromTemplate("templates/models/user.go.txt", gor.RootPath+"/models/user.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFilefromTemplate("templates/models/token.go.txt", gor.RootPath+"/models/token.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFilefromTemplate("templates/models/models.go.txt", gor.RootPath+"/models/models.go")
	if err != nil {
		exitGracefully(err)
	}

	color.Green("✓ Successfully created and executed the migrations for users, tokens, and remember_tokens.")
	color.Green("✓ Successfully generated user and token models.")
	color.Yellow("")
	color.Cyan("Note: Ensure that the user and token models are registered in models/models.go.")
	color.Cyan(`      - Register you custom models in modes/modesl.go for initialization and usage`)

	return nil
}
