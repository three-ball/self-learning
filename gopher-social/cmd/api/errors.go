package main

import (
	"net/http"
)

func (app *application) internalServerError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	app.logger.Errorf("internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

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
	app.logger.Errorf("bad request: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

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
	app.logger.Errorf("not found: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

	writeJSONError(
		w,
		http.StatusNotFound,
		"the requested resource could not be found",
	)
}

func (app *application) unauthorizedError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	app.logger.Errorf("unauthorized: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

	writeJSONError(
		w,
		http.StatusUnauthorized,
		"you are not authorized to access this resource",
	)
}

func (app *application) forbiddenError(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnw("forbidden", "method", r.Method, "path", r.URL.Path, "error")

	writeJSONError(w, http.StatusForbidden, "forbidden")
}
