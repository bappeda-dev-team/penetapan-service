package api

import (
	"net/http"

	"github.com/bappeda-dev-team/penetapan-service/internal/model/web"
)

// healthcheck godoc
// @Summary healthcheck
// @Description checking if this service is accessible
// @Tags healthcheck
// @Produce json
// @Success 200	{object} web.Response[web.HealthcheckResponse]
// @Failure 500 {object} web.ErrorResponse
// @Router /healthcheck [get]
func (app *Application) HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := web.Response[web.HealthcheckResponse]{
		Data: web.HealthcheckResponse{
			Status: "up",
			SystemInfo: web.HealthcheckSystemInfo{
				Environment: app.Config.Env,
				Version:     Version,
			},
		},
	}

	err := app.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
	}
}
