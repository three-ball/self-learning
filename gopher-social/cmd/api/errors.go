package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	log.Printf("internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

	writeJSONError(
		w,
		http.StatusInternalServerError,
		"the server encountered an internal error and was unable to complete your request",
	)
}

func (app *application) badRequestError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	log.Printf("bad request: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

	writeJSONError(
		w,
		http.StatusBadRequest,
		err.Error(),
	)
}

func (app *application) notFoundError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	log.Printf("not found: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

	writeJSONError(
		w,
		http.StatusNotFound,
		"the requested resource could not be found",
	)
}
