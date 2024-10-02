package model

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type MovieRepository interface {
	Insert(movie *Movie) error
	Get(id int64) (*Movie, error)
	Update(movie *Movie) error
	Delete(id int64) error
}

// Create a Models struct which wraps the MovieModel. We'll add other models to this,
// like a UserModel and PermissionModel, as our build progresses.
type Models struct {
	Movies MovieRepository
}

// For ease of use, we also add a New() method which returns a Models struct containing
// the initialized MovieModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
	}
}
