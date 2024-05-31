package goravel

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/saalikmubeen/goravel/render"
	"github.com/saalikmubeen/goravel/session"
)

const (
	version = "1.0.0"
)

type Goravel struct {
	AppName  string
	Debug    bool // true for development mode
	Version  string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	RootPath string // rootPath is the path that we are in when we start the goravel app
	Render   *render.Render
	Routes   *chi.Mux
	JetViews *jet.Set
	Session  *scs.SessionManager
	DB       Database

	// not exported, used internally
	// contains mostly loaded environment variables.
	config config
}

type config struct {
	port        string // port that the server will listen on
	renderer    string // name of the rendering engine that we want to use ("go" or" jet")
	cookie      cookieConfig
	sessionType string
	database    databaseConfig
	// redis       redisConfig
	// uploads     uploadConfig
}

// CreateFolderStructure creates necessary folders for our Goravel application
func (g *Goravel) CreateFolderStructure(p initPaths) error {
	rootPath := p.rootPath // string that holds the full pathname to the root level of my web app

	for _, folderName := range p.folderNames {
		// create a folder in the rootPath if it doesn't exist
		err := g.CreateDirIfNotExists(fmt.Sprintf("%s/%s", rootPath, folderName))

		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Goravel) checkDotEnv(path string) error {
	err := g.CreateFileIfNotExists(fmt.Sprintf("%s/.env", path))
	if err != nil {
		return err
	}
	return nil
}

func (g *Goravel) CreateLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}

func (g *Goravel) createRenderer() {
	myRenderer := render.Render{
		Renderer: g.config.renderer,
		RootPath: g.RootPath,
		Port:     g.config.port,
		JetViews: g.JetViews,
		Session:  g.Session,
	}

	g.Render = &myRenderer
}

func (g *Goravel) BuildDSN() string {

	var dsn string = ""
	dbType := os.Getenv("DATABASE_TYPE")

	switch dbType {
	case "postgres", "postgresql":
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",
			os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_NAME"), os.Getenv("DATABASE_SSLMODE"))

		if os.Getenv("DATABASE_PASSWORD") != "" {
			dsn = fmt.Sprintf("%s password=%s", dsn, os.Getenv("DATABASE_PASSWORD"))
		}
	case "mysql":

	default:

	}

	return dsn
}

func (g *Goravel) New(rootPath string) error {

	paths := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "views", "mail", "models", "public", "tmp", "logs", "middleware", "screenshots"},
	}

	// ** create the necessary folders
	err := g.CreateFolderStructure(paths)
	if err != nil {
		return err
	}

	// ** Check if the .env file exists, if not create it
	err = g.checkDotEnv(rootPath)
	if err != nil {
		return err
	}

	// Load the .env file into the environment
	// read .env
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	// ** create the loggers
	infoLog, errorLog := g.CreateLoggers()
	g.InfoLog = infoLog
	g.ErrorLog = errorLog

	g.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	g.Version = version
	g.RootPath = rootPath

	g.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
		cookie: cookieConfig{
			name:     os.Getenv("COOKIE_NAME"),
			lifetime: os.Getenv("COOKIE_LIFETIME"),
			persist:  os.Getenv("COOKIE_PERSISTS"),
			secure:   os.Getenv("COOKIE_SECURE"),
			domain:   os.Getenv("COOKIE_DOMAIN"),
		},
		database: databaseConfig{
			dsn:          g.BuildDSN(),
			databaseType: os.Getenv("DATABASE_TYPE"),
		},
	}

	// ** Initialize and create the Jet views
	var views = jet.NewSet(
		jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rootPath)),
		jet.InDevelopmentMode(),
	)

	g.JetViews = views

	// ** create the renderer that renders our templates
	g.createRenderer()

	//**  create the routes
	g.Routes = g.initRoutes().(*chi.Mux)

	// ** connect to the database
	dbType := os.Getenv("DATABASE_TYPE")
	if dbType != "" {
		dsn := g.BuildDSN()

		db, err := g.ConnectToDatabase(dbType, dsn)

		if err != nil {
			errorLog.Panicln(err)
			os.Exit(1)
		}

		g.DB = Database{
			DatabaseType: dbType,
			Pool:         db,
		}
	}

	// ** Create and initialize the session
	session := session.Session{
		CookieLifetime: g.config.cookie.lifetime,
		CookiePersist:  g.config.cookie.persist,
		CookieName:     g.config.cookie.name,
		CookieDomain:   g.config.cookie.domain,
		CookieSecure:   g.config.cookie.secure,
		SessionType:    g.config.sessionType,
		DBPool:         g.DB.Pool,
	}
	g.Session = session.InitSession()

	return nil

}

// ListenAndServe starts the web server
func (g *Goravel) ListenAndServe() {
	port := os.Getenv("PORT")
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      g.Routes,
		ErrorLog:     g.ErrorLog,
		IdleTimeout:  time.Second * 30,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 600,
	}

	defer g.DB.Pool.Close() // close the database connection when the server stops

	g.InfoLog.Printf("Starting server on port %s", port)
	err := srv.ListenAndServe()
	g.ErrorLog.Fatal(err)
}
