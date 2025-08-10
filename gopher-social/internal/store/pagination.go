package store

import (
	"net/http"
	"strconv"
	"strings"
)

type PaginatedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=100"`
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"omitempty,oneof=asc desc"`
	Tags   []string `json:"tags,omitempty" validate:"dive,max=5"`
	Search string   `json:"search,omitempty" validate:"max=100"`
}

type PaginateOpts func(*PaginatedQuery)

// NewaginatedQuery creates a new PaginatedQuery with default values.
// It can be used to set default pagination parameters or to create a new query with custom values.
func NewPaginatedQuery(opts ...PaginateOpts) *PaginatedQuery {
	pq := PaginatedQuery{
		Limit:  10,         // Default limit
		Offset: 0,          // Default offset
		Sort:   "desc",     // Default sort order
		Tags:   []string{}, // Default empty tags
	}

	for _, opt := range opts {
		opt(&pq)
	}

	return &pq
}

// WithLimit sets the limit for the pagination query.
func WithLimit(limit int) PaginateOpts {
	return func(pq *PaginatedQuery) {
		if limit > 0 && limit <= 100 {
			pq.Limit = limit
		}
	}
}

// WithOffset sets the offset for the pagination query.
func WithOffset(offset int) PaginateOpts {
	return func(pq *PaginatedQuery) {
		if offset >= 0 {
			pq.Offset = offset
		}
	}
}

// WithSort sets the sort order for the pagination query.
func WithSort(sort string) PaginateOpts {
	return func(pq *PaginatedQuery) {
		if sort == "asc" || sort == "desc" {
			pq.Sort = sort
		}
	}
}

func (pq *PaginatedQuery) Parse(r *http.Request) (*PaginatedQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		parsedLimit, err := strconv.Atoi(limit)
		if err != nil || parsedLimit < 1 || parsedLimit > 100 {
			return pq, err
		}
		pq.Limit = parsedLimit
	}

	offset := qs.Get("offset")
	if offset != "" {
		parsedOffset, err := strconv.Atoi(offset)
		if err != nil || parsedOffset < 0 {
			return pq, err
		}
		pq.Offset = parsedOffset
	}

	sort := qs.Get("sort")
	if sort != "" {
		pq.Sort = sort
	}

	tags := qs.Get("tags")
	if tags != "" {
		pq.Tags = strings.Split(tags, ",")
	}

	search := qs.Get("search")
	if search != "" {
		pq.Search = search
	}

	return pq, nil
}
