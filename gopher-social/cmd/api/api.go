package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2" // http-swagger middleware

	// Assuming you have a doc package for Swagger documentation
	"github.com/three-ball/gopher-social/docs"
	"github.com/three-ball/gopher-social/internal/auth"
	"github.com/three-ball/gopher-social/internal/env"
	"github.com/three-ball/gopher-social/internal/mailer"
	"github.com/three-ball/gopher-social/internal/store"
	"go.uber.org/zap"
)

type application struct {
	config        config
	store         store.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
}

type config struct {
	addr        string
	db          dbConfig
	env         string
	apiURL      string // URL for the API, used in Swagger documentation
	frontendURL string
	mail        mailConfig
	auth        authConfig
}

type authConfig struct {
	token tokenConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type mailConfig struct {
	sendGrid  sendGridConfig
	mailTrap  mailTrapConfig
	fromEmail string
	exp       time.Duration
}

type mailTrapConfig struct {
	apiKey string
}

type sendGridConfig struct {
	apiKey string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (a *application) mount() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{env.GetString("CORS_ALLOWED_ORIGIN", "http://localhost:5174")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", a.healthCheckHandler)

		docURL := fmt.Sprintf("%s/doc.json", a.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docURL)))

		// Posts routes
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", a.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				// middleware to load post data into context
				r.Use(a.postsContextMiddleware)

				r.Get("/", a.getPostHandler)
				r.Patch("/", a.checkPostOwnership("moderator", a.patchPostHandler)) // only if the user is the owner of the post or has moderator role
				r.Delete("/", a.checkPostOwnership("admin", a.deletePostHandler))   // only if the user is the owner of the post or has admin role

			})
		})

		// Comments routes
		r.Route("/comments", func(r chi.Router) {
			r.Post("/", a.createCommentHandler) // Assuming createCommentHandler is defined
		})

		// User routes
		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", a.activateUserHandler) // Activate user with token

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(a.AuthTokenMiddleware)

				r.Get("/", a.getUserHandler)
				r.Put("/follow", a.followUserHandler)     // Using PUT instead of POST for idempotency: Following/unfollowing a user multiple times should have the same result
				r.Put("/unfollow", a.unfollowUserHandler) // Using PUT instead of POST for idempotency: Following/unfollowing a user multiple times should have the same result
			})

			r.Group(func(r chi.Router) {
				r.Use(a.AuthTokenMiddleware)
				r.Get("/feed", a.getUserFeedHandler) // just put feed here for simplicity, we will add auth later
			})

			r.Route("/authentication", func(r chi.Router) {
				r.Post("/user", a.registerUserHandler) // Register a new user
				r.Post("/token", a.createTokenHandler)
			})
		})
	})

	return r
}

func (app *application) run(r *chi.Mux) error {
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      r,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Infow("server has started", "addr", app.config.addr, "env", app.config.env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Infow("server has stopped", "addr", app.config.addr, "env", app.config.env)

	return nil
}
