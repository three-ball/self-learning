package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/three-ball/gopher-social/internal/store"
)

func Seed(store *store.Storage) error {
	ctx := context.Background()
	// Seed users
	users := generateUsers(100, 67000)
	for index, user := range users {
		if err := store.Users.Create(ctx, &user); err != nil {
			log.Printf("Error creating user %s: %v", user.Username, err)
			return err
		}
		users[index].ID = user.ID // Update the first user ID for post creation
	}
	// Seed posts
	posts := generatePosts(200, users)
	for index, post := range posts {
		if err := store.Posts.Create(ctx, &post); err != nil {
			log.Printf("Error creating post for user %d: %v", post.UserID, err)
			return err
		}
		// Update post with the created ID
		posts[index].ID = post.ID // This assumes the Create method sets the ID on the post
	}
	// Seed comments
	comments := generateComments(5000, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, &comment); err != nil {
			log.Printf("Error creating comment for post %d: %v", comment.PostID, err)
			return err
		}
	}
	return nil
}

func generateUsers(num int, from int) []store.User {
	users := make([]store.User, num)
	for i := 0; i < num; i++ {
		users[i] = store.User{
			Username:  "user" + fmt.Sprintf("%d", i+from),
			Email:     "user" + fmt.Sprintf("%d", i+from) + "@example.com",
			Password:  "password" + fmt.Sprintf("%d", i+from),
			CreatedAt: time.Now().Format(time.RFC3339),
		}
	}
	return users
}

func generatePosts(num int, users []store.User) []store.Post {
	posts := make([]store.Post, num)

	for i := 0; i < num; i++ {
		// pick random users from the provided list
		if len(users) == 0 {
			return nil // No users to create posts for
		}

		posts[i] = store.Post{
			Content:   "This is post content " + fmt.Sprintf("%d", i),
			Title:     "Post Title " + fmt.Sprintf("%d", i),
			UserID:    users[i%len(users)].ID,
			Tags:      []string{"tag1", "tag2"},
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
			Version:   1,
		}
	}
	return posts
}

func generateComments(num int, posts []store.Post) []store.Comment {
	comments := make([]store.Comment, num)
	for i := 0; i < num; i++ {
		comments[i] = store.Comment{
			Content:   "This is comment content " + fmt.Sprintf("%d", i),
			PostID:    posts[i%len(posts)].ID,
			UserID:    posts[i%len(posts)].UserID, // Assuming the user who created the post also comments
			CreatedAt: time.Now().Format(time.RFC3339),
		}
	}
	return comments
}
