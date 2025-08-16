package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/three-ball/gopher-social/internal/store"
)

// func decodeBase64(encoded string) (string, error) {
// 	decoded, err := base64.StdEncoding.DecodeString(encoded)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(decoded), nil
// }

// func (app *application) BasicAuthencateMW() func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			authHeader := r.Header.Get("Authorization")
// 			if authHeader == "" {
// 				app.unauthorizedError(w, r, errors.New("missing Authorization header"))
// 				return
// 			}

// 			// parse base64 encoded credentials
// 			parts := strings.SplitN(authHeader, " ", 2)
// 			if len(parts) != 2 || parts[0] != "Basic" {
// 				app.unauthorizedError(w, r, errors.New("invalid Authorization header format"))
// 				return
// 			}

// 			credentials, err := decodeBase64(parts[1])
// 			if err != nil {
// 				app.unauthorizedError(w, r, errors.New("invalid base64 encoding"))
// 				return
// 			}

// 			// split username and password
// 			creds := strings.SplitN(credentials, ":", 2)
// 			if len(creds) != 2 {
// 				app.unauthorizedError(w, r, errors.New("invalid credentials format"))
// 				return
// 			}

// 			username, password := creds[0], creds[1]
// 			if username == "" || password == "" {
// 				app.unauthorizedError(w, r, errors.New("username or password cannot be empty"))
// 				return
// 			}

// 			// authenticate user
// 			user, err := app.store.Users.GetByEmail(ctx, username)
// 			if err != nil {
// 				app.unauthorizedError(w, r, errors.New("user not found"))
// 				return
// 			}

// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedError(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedError(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		token := parts[1]
		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)

		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetByID(ctx, userID)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, uck, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) checkPostOwnership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromCtx(r)
		post := getPostFromCtx(r)

		// RULE 1: Owner can always access their own content
		if post.UserID == user.ID {
			next.ServeHTTP(w, r)
			return
		}

		// RULE 2: Check if user has sufficient role level
		allowed, err := app.checkRolePrecedence(r.Context(), user, requiredRole)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		if !allowed {
			app.forbiddenError(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil
}
