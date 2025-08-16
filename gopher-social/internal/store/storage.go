package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	// ErrNotFound is returned when a requested resource is not found.
	ErrNotFound = errors.New("resource not found")
	// ErrEntityExists is returned when trying to create an entity that already exists.
	ErrEntityExists = errors.New("entity already exists")
	// ErrDuplicateEmail is returned when the email already exists.
	ErrDuplicateEmail = errors.New("a user with that email already exists")
	// ErrDuplicateUsername is returned when the username already exists.
	ErrDuplicateUsername = errors.New("a user with that username already exists")
)

type Storage struct {
	Posts interface {
		Create(ctx context.Context, post *Post) error
		GetByID(ctx context.Context, id int64) (*Post, error)
		Update(ctx context.Context, post *Post) error
		Delete(ctx context.Context, id int64) error
		GetUserFeed(ctx context.Context, userID int64, pgq *PaginatedQuery) ([]PostWithMetadata, error)
	}

	Users interface {
		GetByID(ctx context.Context, id int64) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
		Create(context.Context, *sql.Tx, *User) error
		Update(ctx context.Context, user *User) error
		Delete(ctx context.Context, id int64) error
		Activate(context.Context, string) error
		CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error
	}

	Comments interface {
		GetByPostID(ctx context.Context, postID int64) ([]*Comment, error)
		Create(ctx context.Context, comment *Comment) error
	}

	Follow interface {
		Follow(ctx context.Context, followerID, followeeID int64) error
		Unfollow(ctx context.Context, followerID, followeeID int64) error
	}

	Roles interface {
		GetByName(context.Context, string) (*Role, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostsStore{db: db},
		Users:    &UsersStore{db: db},
		Comments: &CommentStore{db: db},
		Follow:   &FollowStore{db: db},
		Roles:    &RoleStore{db: db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
