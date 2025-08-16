package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/three-ball/gopher-social/internal/store"
)

var (
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
)

type userContextKey string

const (
	uck userContextKey = "user"
)

func getUserFromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(uck).(*store.User)
	return user
}

func (app *application) usersContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		if userID == "" {
			app.badRequestError(w, r, errors.New("user ID is required"))
			return
		}
		id, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}
		user, err := app.store.Users.GetByID(r.Context(), id)
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
		ctx := context.WithValue(r.Context(), uck, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	// Revert back to authentication middleware to get user from context
	user := getUserFromCtx(r)

	if err := writeJSON(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

type UpdateUserPayload struct {
	Username *string `json:"username" validate:"omitempty,min=3,max=50"`
	Email    *string `json:"email" validate:"omitempty,email"`
}

func (app *application) patchUserHandler(w http.ResponseWriter, r *http.Request) {
	existingUser := getUserFromCtx(r)

	// Read the update payload
	var payload UpdateUserPayload
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
	if payload.Username == nil && payload.Email == nil {
		app.badRequestError(w, r, errors.New("at least one field must be provided for update"))
		return
	}

	// Update only the provided fields (PATCH semantics)
	if payload.Username != nil {
		existingUser.Username = *payload.Username
	}
	if payload.Email != nil {
		existingUser.Email = *payload.Email
	}

	// Save the updated user
	if err := app.store.Users.Update(r.Context(), existingUser); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, existingUser); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		app.badRequestError(w, r, errors.New("user ID is required"))
		return
	}

	// Convert userID to int64
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Delete the user
	if err := app.store.Users.Delete(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := writeJSON(w, http.StatusOK, map[string]string{
		"message": "User deleted successfully",
	}); err != nil {
		app.internalServerError(w, r, err)
	}
}

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

// FollowUser godoc
//
//	@Summary		Follows a user
//	@Description	Follows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User followed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromCtx(r)
	var follow FollowUser
	if err := readJSON(w, r, &follow); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if follow.UserID == 0 {
		app.badRequestError(w, r, errors.New("user_id is required"))
		return
	}
	if followerUser.ID == follow.UserID {
		app.badRequestError(w, r, errors.New("you cannot follow yourself"))
		return
	}
	if err := app.store.Follow.Follow(r.Context(), followerUser.ID, follow.UserID); err != nil {
		switch {
		case errors.Is(err, store.ErrEntityExists):
			app.badRequestError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}
	if err := writeJSON(w, http.StatusCreated, map[string]string{
		"message": "User followed successfully",
		"user_id": strconv.FormatInt(follow.UserID, 10),
	}); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// UnfollowUser gdoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User unfollowed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromCtx(r)
	unfollowUserID := chi.URLParam(r, "userID")
	if unfollowUserID == "" {
		app.badRequestError(w, r, errors.New("user ID is required"))
		return
	}
	id, err := strconv.ParseInt(unfollowUserID, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if followerUser.ID == id {
		app.badRequestError(w, r, errors.New("you cannot unfollow yourself"))
		return
	}
	if err := app.store.Follow.Unfollow(r.Context(), followerUser.ID, id); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}
	if err := writeJSON(w, http.StatusOK, map[string]string{
		"message": "User unfollowed successfully",
	}); err != nil {
		app.internalServerError(w, r, err)
	}
}

// ActivateUser godoc
//
//	@Summary		Activates/Register a user
//	@Description	Activates/Register a user by invitation token
//	@Tags			users
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := app.store.Users.Activate(r.Context(), token)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, ""); err != nil {
		app.internalServerError(w, r, err)
	}
}
