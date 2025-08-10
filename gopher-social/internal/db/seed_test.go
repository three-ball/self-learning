package db

import (
	"testing"

	"github.com/three-ball/gopher-social/internal/store"
)

func TestSeed(t *testing.T) {
	// Initialize the database connection
	db, err := New("postgres://admin:adminpassword@localhost:5432/socialnetwork?sslmode=disable", 10, 5, "30s")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create a new storage instance
	storage := store.NewStorage(db)

	// Seed the database with initial data
	if err := Seed(&storage); err != nil {
		t.Fatalf("failed to seed database: %v", err)
	}
}
