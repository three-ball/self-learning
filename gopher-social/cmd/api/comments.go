package main

import (
	"net/http"

	"github.com/three-ball/gopher-social/internal/store"
)

type CreateCommentPayload struct {
	PostID  int64  `json:"post_id" validate:"required"`
	UserID  int64  `json:"user_id" validate:"required"`
	Content string `json:"content" validate:"required,max=1000"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var comment CreateCommentPayload

	if err := readJSON(w, r, &comment); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(comment); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	newComment := &store.Comment{
		PostID:  comment.PostID,
		UserID:  comment.UserID,
		Content: comment.Content,
	}

	if err := app.store.Comments.Create(r.Context(), newComment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, map[string]string{
		"message": "Comment created successfully",
	}); err != nil {
		app.internalServerError(w, r, err)
	}
}
