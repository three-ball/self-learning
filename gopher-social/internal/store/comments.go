package store

import (
	"context"
	"database/sql"
	"time"
)

/*
CREATE TABLE IF NOT EXISTS comments (

	id bigserial PRIMARY KEY,
	post_id bigserial NOT NULL,
	user_id bigserial NOT NULL,
	content TEXT NOT NULL,
	created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()

);
*/
type Comment struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) GetByPostID(ctx context.Context, postID int64) ([]*Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
		SELECT id, post_id, user_id, content, created_at
		FROM comments WHERE post_id = $1
	`

	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		comment := &Comment{}
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *CommentStore) Create(ctx context.Context, comment *Comment) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
		INSERT INTO comments (post_id, user_id, content)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
