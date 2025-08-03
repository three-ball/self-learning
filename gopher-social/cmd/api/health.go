package main

import "net/http"

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if err := writeJSON(w, http.StatusOK, map[string]string{
		"status":  "available",
		"system":  "gopher-social",
		"version": version,
		"env":     app.config.env,
	}); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Internal Server Error")
	}
}
