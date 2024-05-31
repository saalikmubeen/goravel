package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func showHelp() {
	color.Yellow(`Available commands:

	help                  - show the help commands
	version               - print application version
	migrate               - runs all up migrations that have not been run previously
	migrate down          - reverses the most recent migration
	migrate reset         - runs all down migrations in reverse order, and then all up migrations
	make migration <name> - creates two new migration files one up and one down in the migrations folder
	make auth             - creates and runs migrations for authentication tables, and creates models and middleware
	make handler <name>   - creates a stub handler in the handlers directory
	make model <name>     - creates a new model in the data directory
	make session          - creates a table in the database as a session store
	`)
}

// required because migrate package and go's SQL require different DSN formats
func buildDSN() string {
	dbType := gor.DB.DatabaseType

	if dbType == "postgres" || dbType == "postgresql" {
		var dsn string
		if os.Getenv("DATABASE_PASSWORD") != "" {
			dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_PASSWORD"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSLMODE"))
		} else {
			dsn = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSLMODE"))
		}
		return dsn
	}
	return "mysql://" + gor.BuildDSN()
}
