package goravel

import (
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

// ListenAndServe starts the web server
func (g *Goravel) ListenAndServe() error {
	port := os.Getenv("PORT")
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      g.Routes,
		ErrorLog:     g.ErrorLog,
		IdleTimeout:  time.Second * 30,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 600,
	}

	if g.DB.Pool != nil {
		defer g.DB.Pool.Close() // close the database connection when the server stops
	}

	if redisPool != nil {
		defer redisPool.Close() // close the redis connection when the server stops
	}

	if badgerConn != nil {
		defer badgerConn.Close()
	}

	color.Yellow(Banner, Version)
	color.Green("Starting server on port %s", port)
	err := srv.ListenAndServe()

	if err != nil {
		return err
	} else {
		return nil
	}
}
