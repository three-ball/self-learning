package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64    `json:"id"`
	Content   string   `json:"content"`
	Title     string   `json:"title"`
	UserID    int64    `json:"user_id"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`

	Comments []*Comment `json:"comments"`
}

type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) Create(ctx context.Context, post *Post) error {
	// Implementation for creating a post in the database
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostsStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	// Implementation for retrieving a post by ID from the database
	query := `
		SELECT id, content, title, user_id, tags, created_at, updated_at
		FROM posts WHERE id = $1
	`

	post := &Post{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return post, nil
}

func (s *PostsStore) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts
		SET content = $2, title = $3, tags = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.ID,
		post.Content,
		post.Title,
		pq.Array(post.Tags),
	).Scan(&post.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *PostsStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
