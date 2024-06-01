package goravel

import "database/sql"

type initPaths struct {
	rootPath    string   // rootPath is the path that we are in when we start the goravel app
	folderNames []string // folderNames is the names of the folders that we need to create in the rootPath
}

type cookieConfig struct {
	name     string
	lifetime string
	persist  string
	secure   string
	domain   string
}

type databaseConfig struct {
	dsn          string
	databaseType string
}

type Database struct {
	DatabaseType string
	Pool         *sql.DB
}

type redisConfig struct {
	host     string
	password string
	prefix   string
}
