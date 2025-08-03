package main

import (
	"log"
	"net/http"
	"time"
)

type application struct {
	config config
}

type config struct {
	addr string
}

func (a *application) mount() *http.ServeMux {
	mux := http.NewServeMux()
	// Here you would mount your routes, e.g.:
	// mux.HandleFunc("/api/v1/posts", a.handlePosts)
	mux.HandleFunc("GET /v1/health", a.healthCheckHandler)
	return mux
}

func (app *application) run(mux *http.ServeMux) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	log.Println("Starting server on", app.config.addr)

	return srv.ListenAndServe()
}
