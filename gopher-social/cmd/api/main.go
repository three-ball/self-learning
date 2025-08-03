package main

import (
	"log"

	"github.com/three-ball/gopher-social/internal/db"
	"github.com/three-ball/gopher-social/internal/env"
	"github.com/three-ball/gopher-social/internal/store"
)

const version = "1.0.0"

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":3000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	log.Printf("Connected to database at %s", cfg.db.addr)

	store := store.NewStorage(db) // Replace nil with actual database connection

	app := &application{
		config: cfg,
		store:  store,
	}

	log.Fatal(app.run(app.mount()))
}
