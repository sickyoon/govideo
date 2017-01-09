package govideo

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/burntsushi/toml"
	"github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
	"github.com/sickyoon/govideo/govideo/models"
)

// App is GoVideo Web Application
// TODO: add securecookie
type App struct {
	*httprouter.Router
	handlers http.Handler
	config   models.Config
	db       *MongoClient
}

// NewApp creates new web application
func NewApp(configFile string) *App {

	log.Printf("Initializing web application with " + configFile)

	// initialize app
	app := App{
		Router: httprouter.New(),
	}

	// load config file if exists
	if configFile != "" {
		if _, err := toml.DecodeFile(configFile, &app.config); err != nil {
			log.Panic(err)
		}
	}

	// establish db connection
	app.db = NewMongoClient(app.config.Database.URI, app.config.Database.DBName)

	// TODO: add handlers

	app.GET("/", app.index)

	// TODO: list returns json list of all available media
	// in paths specified in configuration file
	//app.GET("/list", app.list)

	// add static resources handler
	staticPath := "static"
	if app.config.Server.StaticPath != "" {
		staticPath = app.config.Server.StaticPath
	}
	app.ServeFiles("/static/*filepath", fileOnlyFs{http.Dir(staticPath)})

	// add middlewares
	h := handlers.LoggingHandler(os.Stdout, app)
	h = handlers.ProxyHeaders(h)
	h = handlers.CompressHandler(h)
	h = handlers.RecoveryHandler()(h)
	app.handlers = h

	return &app
}

// Run -
func (a *App) Run() {
	log.Printf("Running server at port " + strconv.Itoa(a.config.Server.Port))
	log.Fatal(
		http.ListenAndServe(":"+strconv.Itoa(a.config.Server.Port), a.handlers),
	)
}

type fileOnlyFs struct {
	fs http.FileSystem
}

func (fs fileOnlyFs) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if stat.IsDir() {
		return nil, os.ErrNotExist
	}
	return f, nil
}
