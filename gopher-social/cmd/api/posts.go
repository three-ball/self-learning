package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/three-ball/gopher-social/internal/store"
)

type postContextKey string

const (
	pck postContextKey = "post"
)

func getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(pck).(*store.Post)
	return post
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postID := chi.URLParam(r, "postID")
		if postID == "" {
			app.badRequestError(w, r, errors.New("post ID is required"))
			return
		}

		id, err := strconv.ParseInt(postID, 10, 64)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		post, err := app.store.Posts.GetByID(r.Context(), id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundError(w, r, err)
				return
			default:
				app.internalServerError(w, r, err)
				return
			}
		}

		ctx := context.WithValue(r.Context(), pck, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=10000"`
	Tags    []string `json:"tags"`
}

// CreatePost godoc
//
//	@Summary		Creates a post
//	@Description	Creates a post
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreatePostPayload	true	"Post payload"
//	@Success		201		{object}	store.Post
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	// Read the post data from the request body
	var post CreatePostPayload

	if err := readJSON(w, r, &post); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(post); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Create a new Post instance
	newPost := &store.Post{
		Title:   post.Title,
		Content: post.Content,
		UserID:  1,
		Tags:    post.Tags,
	}

	// Use the PostsStore to create the post in the database
	if err := app.store.Posts.Create(r.Context(), newPost); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Placeholder for post creation logic
	// This will handle the creation of a new post
	// For now, we can just return a success message
	if err := app.jsonResponse(w, http.StatusCreated, map[string]string{
		"message": "Post created successfully",
		"post_id": strconv.FormatInt(newPost.ID, 10),
	}); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetPost godoc
//
//	@Summary		Fetches a post
//	@Description	Fetches a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	store.Post
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [get]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

type UpdatePostPayload struct {
	Title   *string  `json:"title" validate:"omitempty,max=100"`
	Content *string  `json:"content" validate:"omitempty,max=10000"`
	Tags    []string `json:"tags"`
}

// UpdatePost godoc
//
//	@Summary		Updates a post
//	@Description	Updates a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Post ID"
//	@Param			payload	body		UpdatePostPayload	true	"Post payload"
//	@Success		200		{object}	store.Post
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [patch]
func (app *application) patchPostHandler(w http.ResponseWriter, r *http.Request) {
	existingPost := getPostFromCtx(r)

	// Read the update payload
	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate only the fields that are provided
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Check if at least one field is provided for update
	if payload.Title == nil && payload.Content == nil && payload.Tags == nil {
		app.badRequestError(w, r, errors.New("at least one field must be provided for update"))
		return
	}

	// Update only the provided fields (PATCH semantics)
	if payload.Title != nil {
		existingPost.Title = *payload.Title
	}
	if payload.Content != nil {
		existingPost.Content = *payload.Content
	}
	if payload.Tags != nil {
		existingPost.Tags = payload.Tags
	}

	// Save the updated post
	if err := app.store.Posts.Update(r.Context(), existingPost); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, existingPost); err != nil {
		app.internalServerError(w, r, err)
	}
}

// DeletePost godoc
//
//	@Summary		Deletes a post
//	@Description	Delete a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		204	{object}	string
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	if postID == "" {
		app.badRequestError(w, r, errors.New("post ID is required"))
		return
	}

	// Convert postID to int64
	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Delete the post
	if err := app.store.Posts.Delete(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, map[string]string{
		"message": "Post deleted successfully",
	}); err != nil {
		app.internalServerError(w, r, err)
	}
}
