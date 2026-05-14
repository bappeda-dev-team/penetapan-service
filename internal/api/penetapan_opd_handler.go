package api

import (
	"net/http"
	"strconv"

	"github.com/bappeda-dev-team/penetapan-service/internal/model/domain"
	"github.com/bappeda-dev-team/penetapan-service/internal/model/web"
	"github.com/bappeda-dev-team/penetapan-service/internal/validator"
)

// SyncPenetapanOpdHandler godoc
//
// @Summary     Sync Penetapan OPD
// @Description Sinkron data penetapan OPD berdasarkan kode OPD, tahun dan jenis penetapan dari perencanaan
// @Tags        Sync
// @Accept      json
// @Produce     json
//
// @Param       payload body web.SyncPenetapanOpdRequest true "Payload sinkronisasi penetapan OPD"
//
// @Success     200 {object} web.Response[web.SyncPenetapanOpdResponse] "Berhasil mengambil data tujuan OPD"
// @Failure     422 {object} web.ValidationErrorResponse                "Unprocessable Entity"
// @Failure     500 {object} web.ErrorResponse                          "Internal Server Error"
//
// @Router      /opd/sync_penetapan [post]
func (app *Application) SyncPenetapanOpdHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	input := web.SyncPenetapanOpdRequest{}
	errors := map[string]string{}

	err := app.ReadJSON(w, r, &input)
	if err != nil {
		errors["invalid_request"] = err.Error()
		app.BadRequestResponse(
			w,
			r,
			web.ValidationErrorResponse{
				Error: errors,
			},
		)
		return
	}

	request := &web.SyncPenetapanOpdRequest{
		KodeOpd:        input.KodeOpd,
		Tahun:          input.Tahun,
		JenisPenetapan: input.JenisPenetapan,
	}
	v := validator.New()
	if web.ValidateSyncPenetapanRequest(v, request); !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	result, err := app.PenetapanOpdService.SyncPenetapanOpd(r.Context(), request)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	response := web.Response[web.SyncPenetapanOpdResponse]{
		Data: result,
	}
	err = app.WriteJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
	}
}

// TujuanOpdHandler godoc
//
// @Summary     Get tujuan OPD penetapan
// @Description Mengambil data tujuan OPD berdasarkan kode OPD dan tahun penetapan
// @Tags        OPD
// @Accept      json
// @Produce     json
//
// @Param       kodeOpd query string true "Kode OPD"
// @Param       tahun   query int    true "Tahun Penetapan"
//
// @Success     200 {array}  web.Response[web.TujuanPenetapanOpdResponse] "Berhasil mengambil data tujuan OPD"
// @Failure     400 {object} web.ValidationErrorResponse                  "Bad Request"
// @Failure     500 {object} web.ErrorResponse                            "Internal Server Error"
//
// @Router      /opd/tujuan [get]
func (app *Application) TujuanOpdHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	query := r.URL.Query()

	errors := map[string]string{}

	kodeOpd := query.Get("kodeOpd")

	tahunStr := query.Get("tahun")

	var tahun int

	if kodeOpd == "" {
		errors["kode_opd"] = "required"
	}

	if tahunStr == "" {

		errors["tahun"] = "required"

	} else {

		parsedTahun, err := strconv.Atoi(tahunStr)

		if err != nil {

			errors["tahun"] = "tahun tidak valid"

		} else {

			tahun = parsedTahun
		}
	}

	if len(errors) > 0 {

		app.BadRequestResponse(
			w,
			r,
			web.ValidationErrorResponse{
				Error: errors,
			},
		)

		return
	}

	request := domain.PenetapanOpdRequest{
		KodeOpd: kodeOpd,
		Tahun:   tahun,
	}

	result, err := app.PenetapanOpdService.FindTujuan(
		r.Context(),
		request,
	)

	if err != nil {

		app.ServerErrorResponse(w, r, err)

		return
	}

	response := web.Response[[]web.TujuanPenetapanOpdResponse]{
		Data: result,
	}

	err = app.WriteJSON(
		w,
		http.StatusOK,
		response,
		nil,
	)

	if err != nil {

		app.ServerErrorResponse(w, r, err)

		return
	}
}

// SasaranOpdHandler godoc
//
// @Summary     Get sasaran OPD penetapan
// @Description Mengambil data sasaran OPD berdasarkan kode OPD dan tahun penetapan
// @Tags        OPD
// @Accept      json
// @Produce     json
//
// @Param       kodeOpd query string true "Kode OPD"
// @Param       tahun   query int    true "Tahun Penetapan"
//
// @Success     200 {array}  web.Response[web.SasaranPenetapanOpdResponse] "Berhasil mengambil data sasaran OPD"
// @Failure     400 {object} web.ValidationErrorResponse                   "Bad Request"
// @Failure     500 {object} web.ErrorResponse                             "Internal Server Error"
//
// @Router      /opd/sasaran [get]
func (app *Application) SasaranOpdHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	// query
	query := r.URL.Query()
	errors := map[string]string{}

	kodeOpd := query.Get("kodeOpd")
	tahun, err := strconv.Atoi(query.Get("tahun"))
	if err != nil {
		errors["tahun"] = "tahun tidak valid"
	}
	request := domain.PenetapanOpdRequest{
		KodeOpd: kodeOpd,
		Tahun:   tahun,
	}

	if request.KodeOpd == "" {
		errors["kode_opd"] = "required"
	}

	if request.Tahun == 0 {
		errors["tahun"] = "required"
	}

	if len(errors) > 0 {
		app.BadRequestResponse(
			w,
			r,
			web.ValidationErrorResponse{
				Error: errors,
			},
		)
		return
	}

	result, err := app.PenetapanOpdService.FindSasaran(r.Context(), request)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
	}

	response := web.Response[[]web.SasaranPenetapanOpdResponse]{
		Data: result,
	}
	err = app.WriteJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
	}
}
