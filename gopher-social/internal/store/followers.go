package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowStore struct {
	db *sql.DB
}

func (s *FollowStore) Follow(ctx context.Context, followerID, followeeID int64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO followers (user_id, follower_id)
		VALUES ($1, $2) ON CONFLICT DO NOTHING
	`

	_, err := s.db.ExecContext(ctx, query, followeeID, followerID)
	if err != nil {
		if poErr, ok := err.(*pq.Error); ok && poErr.Code == "23505" {
			return ErrEntityExists
		}
		return err
	}

	return nil
}

func (s *FollowStore) Unfollow(ctx context.Context, followerID, followeeID int64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		DELETE FROM followers
		WHERE user_id = $1 AND follower_id = $2
	`

	_, err := s.db.ExecContext(ctx, query, followeeID, followerID)
	if err != nil {
		return err
	}

	return nil
}
