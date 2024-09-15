package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"greenlight-further/internal/model"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type application struct {
	config config
	logger *log.Logger
	models model.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port,
		"port",
		4000,
		"API server port")

	flag.StringVar(&cfg.env,
		"env",
		"development",
		"Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn,
		"db-dsn",
		"postgres://psql_db:secret@localhost/psql_db?sslmode=disable",
		"PostgreSQL DSN")

	flag.IntVar(
		&cfg.db.maxOpenConns,
		"db-max-open-conns",
		25,
		"PostgreSQL max open connections")

	flag.IntVar(
		&cfg.db.maxIdleConns,
		"db-max-idle-conns",
		25,
		"PostgreSQL max idle connections")

	flag.StringVar(
		&cfg.db.maxIdleTime,
		"db-max-idle-time",
		"15m",
		"PostgreSQL max connection idle time")

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Printf("database connection pool established")

	app := &application{
		config: cfg,
		logger: logger,
		models: model.NewModels(db),
	}

	router := app.routes()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start the HTTP server.
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

// The openDB() function returns a sql.DB connection pool.
func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config
	// struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	// Set the maximum idle timeout.
	db.SetConnMaxIdleTime(duration)

	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	// Return the sql.DB connection pool.
	return db, nil
}
