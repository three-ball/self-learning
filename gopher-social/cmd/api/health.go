package main

import "net/http"

// healthcheckHandler godoc
//
//	@Summary		Healthcheck
//	@Description	Healthcheck endpoint
//	@Tags			ops
//	@Produce		json
//	@Success		200	{object}	string	"ok"
//	@Router			/health [get]
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
