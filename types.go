package goravel

import "database/sql"

type InitPaths struct {
	RootPath    string   // rootPath is the path that we are in when we start the goravel app
	FolderNames []string // folderNames is the names of the folders that we need to create in the rootPath
}

type Server struct {
	ServerName string
	Port       string
	Secure     bool
	URL        string
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
