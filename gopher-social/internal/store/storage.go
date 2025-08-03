package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	// ErrNotFound is returned when a requested resource is not found.
	ErrNotFound = errors.New("resource not found")
)

type Storage struct {
	Posts interface {
		Create(ctx context.Context, post *Post) error
		GetByID(ctx context.Context, id int64) (*Post, error)
		Update(ctx context.Context, post *Post) error
		Delete(ctx context.Context, id int64) error
	}

	Users interface {
		Create(ctx context.Context, user *User) error
	}

	Comments interface {
		GetByPostID(ctx context.Context, postID int64) ([]*Comment, error)
		Create(ctx context.Context, comment *Comment) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostsStore{db: db},
		Users:    &UsersStore{db: db},
		Comments: &CommentStore{db: db},
	}
}
