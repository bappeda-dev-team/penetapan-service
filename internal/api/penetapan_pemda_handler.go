package api

import (
	"net/http"

	"github.com/bappeda-dev-team/penetapan-service/internal/pemda/web"
	"github.com/bappeda-dev-team/penetapan-service/internal/validator"
)

// SyncPenetapanTujuanPemdaHandler godoc
//
// @Summary Sync Tujuan pemda
// @Description Sinkron data tujuan pemda berdasarkan tahun
// @Tags Pemda
// @Accept json
// @Produce json
//
// @Param payload body web.SyncPenetapanRequest true "Payload sync penetapan pemda"
//
// @Success 200 {object} web.Response[web.SyncPenetapanResponse] "Success"
// @Failure 400 {object} web.ValidationErrorResponse "Bad Request"
// @Failure 500 {object} web.ErrorResponse "Internal Server Error"
//
// @Router /pemda/tujuan/sync [post]
func (app *Application) SyncPenetapanTujuanPemdaHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	input := web.SyncPenetapanRequest{}
	errors := map[string]string{}

	err := app.ReadJSON(w, r, &input)
	if err != nil {
		errors["invalid_request"] = err.Error()
		app.BadRequestResponse(
			w, r,
			web.ValidationErrorResponse{
				Error: errors,
			},
		)
		return
	}

	request := &web.SyncPenetapanRequest{
		Tahun: input.Tahun,
	}

	v := validator.New()
	if web.ValidateSyncPenetapanRequest(v, request); !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	result, err := app.PemdaService.SyncPenetapanTujuanPemda(
		r.Context(),
		request,
	)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	response := web.Response[web.SyncPenetapanResponse]{
		Data: result,
	}

	err = app.WriteJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

}
