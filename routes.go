package goravel

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (g *Goravel) initRoutes() http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Recoverer)
	if g.Debug {
		mux.Use(middleware.Logger)
	}

	// load the session
	mux.Use(g.SessionLoad)

	// add the CSRF protection
	mux.Use(g.NoSurf)

	mux.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	return mux

}
