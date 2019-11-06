package api

import (
	"log"
	"net/http"
	"time"

	"github.com/amonaco/goboiler/database"
	"github.com/amonaco/goboiler/logging"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-pg/pg"
	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
)

type App struct {
	db     *pg.DB
	logger *logrus.Logger
}

func (app *App) ExampleEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Print(app.db.PoolStats())
	w.Write([]byte("ok"))
}

func New(config *toml.Tree) (*chi.Mux, error) {
	logger := logging.NewLogger(config)

	db, err := database.DBConn(config)
	if err != nil {
		logger.WithField("module", "database").Error(err)
		return nil, err
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.Timeout(15 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.WithValue("DBCONN", db))

	app := &App{
		db:     db,
		logger: logger,
		// eventually auth or redis
	}

	r.Get("/example_endpoint", app.ExampleEndpoint)
	return r, nil
}
