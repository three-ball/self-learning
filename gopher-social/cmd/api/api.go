package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/three-ball/gopher-social/internal/store"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
	db   dbConfig
	env  string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (a *application) mount() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", a.healthCheckHandler)

		// Posts routes
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", a.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				// middleware to load post data into context
				r.Use(a.postsContextMiddleware)

				r.Get("/", a.getPostHandler)
				r.Patch("/", a.patchPostHandler)
				r.Delete("/", a.deletePostHandler)
			})
		})

		// Comments routes
		r.Route("/comments", func(r chi.Router) {
			r.Post("/", a.createCommentHandler) // Assuming createCommentHandler is defined
		})

		// User routes
		r.Route("/users", func(r chi.Router) {
			r.Post("/", a.createUserHandler)
			r.Route("/{userID}", func(r chi.Router) {
				// middleware to load user data into context
				r.Use(a.usersContextMiddleware)

				r.Get("/", a.getUserHandler)
				r.Patch("/", a.patchUserHandler)
				r.Delete("/", a.deleteUserHandler)
				r.Put("/follow", a.followUserHandler)     // Using PUT instead of POST for idempotency: Following/unfollowing a user multiple times should have the same result
				r.Put("/unfollow", a.unfollowUserHandler) // Using PUT instead of POST for idempotency: Following/unfollowing a user multiple times should have the same result
				r.Get("/feed", a.getUserFeedHandler)      // just put feed here for simplicity, we will add auth later
			})

		})
	})

	return r
}

func (app *application) run(r *chi.Mux) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      r,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	log.Println("Starting server on", app.config.addr)

	return srv.ListenAndServe()
}
