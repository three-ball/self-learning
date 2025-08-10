package main

import (
	"errors"
	"net/http"

	"github.com/three-ball/gopher-social/internal/store"
)

func (a *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// pagination, filtering, and sorting logic can be added here
	pg := store.NewPaginatedQuery()

	pg, err := pg.Parse(r)
	if err != nil {
		a.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(pg); err != nil {
		a.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	feed, err := a.store.Posts.GetUserFeed(ctx, getUserFromCtx(r).ID, pg)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			a.notFoundError(w, r, err)
			return
		default:
			a.internalServerError(w, r, err)
			return
		}
	}

	if err := a.jsonResponse(w, http.StatusOK, feed); err != nil {
		a.internalServerError(w, r, err)
		return
	}
}
