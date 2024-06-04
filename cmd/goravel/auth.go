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
	upFile := gor.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
	downFile := gor.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

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

	// Copy the auth models
	err = copyFilefromTemplate("templates/models/user.go.txt", gor.RootPath+"/models/user.go") // Copy the user model
	if err != nil {
		exitGracefully(err)
	}

	err = copyFilefromTemplate("templates/models/token.go.txt", gor.RootPath+"/models/token.go") // Copy the token model
	if err != nil {
		exitGracefully(err)
	}

	err = copyFilefromTemplate("templates/models/remember-me-token.go.txt", gor.RootPath+"/models/remember_me_token.go") // Copy the remember_me_token model
	if err != nil {
		exitGracefully(err)
	}

	// Copy the auth middlewares
	err = copyFilefromTemplate("templates/middleware/auth.go.txt", gor.RootPath+"/middleware/auth.go") // Copy the auth middleware
	if err != nil {
		exitGracefully(err)
	}

	// Copy the remember_me middleware
	toFile := gor.RootPath + "/middleware/remember_me.go"
	err = handleCopyDataToFile("templates/middleware/remember_me.go.txt", toFile, ReplaceDataMap{ // Copy the remember_me middleware
		"${APP_URL}": gor.GoAppURL,
	})
	if err != nil {
		return err
	}

	// Copy the auth handlers
	toFile = gor.RootPath + "/handlers/auth-handlers.go"
	err = handleCopyDataToFile("templates/handlers/auth_handlers.go.txt", toFile, ReplaceDataMap{
		"${APP_URL}": gor.GoAppURL,
	})
	if err != nil {
		return err
	}

	// Copy the forgot reset mail templates
	err = copyFilefromTemplate("templates/mailer/password-reset.html.tmpl", gor.RootPath+"/mail/password-reset.html.tmpl")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFilefromTemplate("templates/mailer/password-reset.plain.tmpl", gor.RootPath+"/mail/password-reset.plain.tmpl")
	if err != nil {
		exitGracefully(err)
	}

	// Copy the authentication views
	err = copyFilefromTemplate("templates/views/login.jet", gor.RootPath+"/views/login.jet")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFilefromTemplate("templates/views/signup.jet", gor.RootPath+"/views/signup.jet")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFilefromTemplate("templates/views/forgot-password.jet", gor.RootPath+"/views/forgot.jet")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFilefromTemplate("templates/views/reset-password.jet", gor.RootPath+"/views/reset-password.jet")
	if err != nil {
		exitGracefully(err)
	}

	color.Green("✓ Successfully created and executed the migrations for users, tokens, and remember_me_tokens.")
	color.Green("✓ Successfully generated user and token models.")
	color.Green("✓ Successfully created authentication middlewares.")
	color.Yellow("")
	color.Cyan("Note: Ensure that the models are registered in models/models.go.")
	color.Cyan(`      - Register the User, Token, and RememberMeToken models in the models/models.go file.`)
	color.Cyan(`      - Also don't forget to register the generated auth middlewares in the routes.go file.`)

	return nil
}
