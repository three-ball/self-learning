package main

import (
	"errors"
	"net/http"

	"github.com/three-ball/gopher-social/internal/store"
)

// getUserFeedHandler godoc
//
//	@Summary		Fetches the user feed
//	@Description	Fetches the user feed
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			since	query		string	false	"Since"
//	@Param			until	query		string	false	"Until"
//	@Param			limit	query		int		false	"Limit"
//	@Param			offset	query		int		false	"Offset"
//	@Param			sort	query		string	false	"Sort"
//	@Param			tags	query		string	false	"Tags"
//	@Param			search	query		string	false	"Search"
//	@Success		200		{object}	[]store.PostWithMetadata
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
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
