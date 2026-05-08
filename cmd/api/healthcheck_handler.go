package main

import (
	"net/http"
)

// healthcheck godoc
// @Summary healthcheck
// @Description checking is this service is accessible
// @Tags healthcheck
// @Accept json
// @Produce json
// @Success 200	{object} map[string]interface{}	"Service is healthy"
// @Router /v1/healthcheck [get]
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := envelope{
		"status": "up",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
