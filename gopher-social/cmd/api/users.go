package main

import (
	"errors"
	"net/http"

	"github.com/three-ball/gopher-social/internal/store"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := readJSON(w, r, &user); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if user.Email == "" || user.Username == "" || user.Password == "" {
		app.badRequestError(w, r, errors.New("email, username, and password are required"))
		return
	}
	newUser := &store.User{
		Email:    user.Email,
		Username: user.Username,
		Password: user.Password, // In a real application, you should hash the password before storing it
	}
	if err := app.store.Users.Create(r.Context(), newUser); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := writeJSON(w, http.StatusCreated, map[string]string{
		"message": "User created successfully",
	}); err != nil {
		app.internalServerError(w, r, err)
	}
}
