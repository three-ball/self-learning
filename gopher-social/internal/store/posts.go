package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64      `json:"id"`
	Content   string     `json:"content"`
	Title     string     `json:"title"`
	UserID    int64      `json:"user_id"`
	Tags      []string   `json:"tags"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
	Version   int        `json:"version"`
	Comments  []*Comment `json:"comments"`
	User      *User      `json:"user,omitempty"` // Optional user information
}

type PostWithMetadata struct {
	Post
	CommentsCount int `json:"comments_count"`
}

type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) Create(ctx context.Context, post *Post) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
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
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	// Implementation for retrieving a post by ID from the database
	query := `
		SELECT id, content, title, user_id, tags, created_at, updated_at, version
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
		&post.Version, // Added version field
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
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
		UPDATE posts
		SET content = $2, title = $3, tags = $4, updated_at = NOW(), version = version + 1
		WHERE id = $1 AND version = $5
		RETURNING updated_at, version
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.ID,
		post.Content,
		post.Title,
		pq.Array(post.Tags),
		post.Version,
	).Scan(&post.UpdatedAt, &post.Version)

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
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
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

func (s *PostsStore) GetUserFeed(ctx context.Context, userID int64) ([]PostWithMetadata, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	// still can't get this part? Why the instructor use JOIN which OR condition?
	// query := `
	// SELECT
	// 	p.id,
	// 	p.user_id,
	// 	p.title,
	// 	p.content,
	// 	p.created_at,
	// 	p.version,
	// 	p.tags,
	// 	COUNT(c.id) AS comments_count
	// FROM posts p
	// LEFT JOIN comments c ON c.post_id = p.id
	// LEFT JOIN users u ON p.user_id = u.id
	// JOIN followers f ON (f.follower_id = p.user_id OR p.user_id = $1)
	// WHERE f.user_id = $1 OR p.user_id = $1
	// GROUP BY p.id, u.id
	// ORDER BY p.created_at DESC
	// `

	// just my solution as what I understand
	query := `
		select
			p.id,
			p.user_id,
			u.username,
			p.title,
			p."content",
			p.created_at,
			p."version" ,
			p.tags,
			count(c.id) as comments_count
		from
			posts p
		left join "comments" c on
			c.post_id = p.id
		left join users u on
			p.user_id = u.id
		left join followers f on
			p.user_id = f.user_id
		where
			f.follower_id = $1 -- Fetch posts the user is following
			or p.user_id = $1 -- or fetch posts by the user themselves
		group by
			(p.id, u.id)
		order by
			p.created_at desc;`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []PostWithMetadata
	for rows.Next() {
		var post PostWithMetadata
		var username string
		var tags pq.StringArray
		if err := rows.Scan(
			&post.ID,
			&post.UserID,
			&username,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.Version,
			pq.Array(&post.Tags),
			&post.CommentsCount,
		); err != nil {
			return nil, err
		}
		post.Tags = tags
		post.User = &User{Username: username} // Assuming User struct has a Username field
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
