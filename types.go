package goravel

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
