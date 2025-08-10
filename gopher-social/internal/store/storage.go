package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	// ErrNotFound is returned when a requested resource is not found.
	ErrNotFound = errors.New("resource not found")
	// ErrEntityExists is returned when trying to create an entity that already exists.
	ErrEntityExists = errors.New("entity already exists")
)

type Storage struct {
	Posts interface {
		Create(ctx context.Context, post *Post) error
		GetByID(ctx context.Context, id int64) (*Post, error)
		Update(ctx context.Context, post *Post) error
		Delete(ctx context.Context, id int64) error
		GetUserFeed(ctx context.Context, userID int64) ([]PostWithMetadata, error)
	}

	Users interface {
		Create(ctx context.Context, user *User) error
		GetByID(ctx context.Context, id int64) (*User, error)
		Update(ctx context.Context, user *User) error
		Delete(ctx context.Context, id int64) error
	}

	Comments interface {
		GetByPostID(ctx context.Context, postID int64) ([]*Comment, error)
		Create(ctx context.Context, comment *Comment) error
	}

	Follow interface {
		Follow(ctx context.Context, followerID, followeeID int64) error
		Unfollow(ctx context.Context, followerID, followeeID int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostsStore{db: db},
		Users:    &UsersStore{db: db},
		Comments: &CommentStore{db: db},
		Follow:   &FollowStore{db: db},
	}
}
