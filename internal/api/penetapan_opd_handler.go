package api

import (
	"net/http"
	"strconv"

	"github.com/bappeda-dev-team/penetapan-service/internal/model/domain"
	"github.com/bappeda-dev-team/penetapan-service/internal/model/web"
	"github.com/bappeda-dev-team/penetapan-service/internal/validator"
)

// SyncPenetapanSasaranOpdHandler godoc
//
// @Summary     Sync Penetapan Sasaran OPD
// @Description Sinkron data penetapan sasaran OPD berdasarkan kode OPD dan tahun dari perencanaan
// @Tags        Sync
// @Accept      json
// @Produce     json
//
// @Param       payload body web.SyncPenetapanOpdRequest true "Payload sinkronisasi penetapan OPD"
//
// @Success     200 {object} web.Response[web.SyncPenetapanOpdResponse] "Success"
// @Failure     422 {object} web.ValidationErrorResponse                "Unprocessable Entity"
// @Failure     500 {object} web.ErrorResponse                          "Internal Server Error"
//
// @Router      /opd/sasaran/sync [post]
func (app *Application) SyncPenetapanSasaranOpdHandler(
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
		KodeOpd: input.KodeOpd,
		Tahun:   input.Tahun,
	}
	v := validator.New()
	if web.ValidateSyncPenetapanRequest(v, request); !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	result, err := app.PenetapanOpdService.SyncPenetapanOpd(
		r.Context(),
		request,
		domain.JenisPenetapanSasaran,
	)
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

// SyncPenetapanTujuanOpdHandler godoc
//
// @Summary     Sync Penetapan Tujuan OPD
// @Description Sinkron data penetapan tujuan OPD berdasarkan kode OPD dan tahun dari perencanaan
// @Tags        Sync
// @Accept      json
// @Produce     json
//
// @Param       payload body web.SyncPenetapanOpdRequest true "Payload sinkronisasi penetapan OPD"
//
// @Success     200 {object} web.Response[web.SyncPenetapanOpdResponse] "Success"
// @Failure     422 {object} web.ValidationErrorResponse                "Unprocessable Entity"
// @Failure     500 {object} web.ErrorResponse                          "Internal Server Error"
//
// @Router      /opd/tujuan/sync [post]
func (app *Application) SyncPenetapanTujuanOpdHandler(
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
		KodeOpd: input.KodeOpd,
		Tahun:   input.Tahun,
	}
	v := validator.New()
	if web.ValidateSyncPenetapanRequest(v, request); !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	result, err := app.PenetapanOpdService.SyncPenetapanOpd(
		r.Context(),
		request,
		domain.JenisPenetapanTujuan,
	)
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
// @Success     200 {object} web.Response[web.TujuanPenetapanOpdResponse] "Berhasil mengambil data tujuan OPD"
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
		errors["kodeOpd"] = "required"
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

	response := web.Response[web.TujuanPenetapanOpdResponse]{
		Data: result,
	}

	err = app.WriteJSON(w, http.StatusOK, response, nil)
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
// @Success     200 {object} web.Response[web.SasaranPenetapanOpdResponse] "Berhasil mengambil data sasaran OPD"
// @Failure     400 {object} web.ValidationErrorResponse                   "Bad Request"
// @Failure     500 {object} web.ErrorResponse                             "Internal Server Error"
//
// @Router      /opd/sasaran [get]
func (app *Application) SasaranOpdHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	query := r.URL.Query()

	errors := map[string]string{}

	kodeOpd := query.Get("kodeOpd")

	tahunStr := query.Get("tahun")

	var tahun int

	if kodeOpd == "" {
		errors["kodeOpd"] = "required"
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

	result, err := app.PenetapanOpdService.FindSasaran(r.Context(), request)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	response := web.Response[web.SasaranPenetapanOpdResponse]{
		Data: result,
	}

	err = app.WriteJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}

// SyncPenetapanRenjaOpdHandler godoc
//
// @Summary     Sync Penetapan Renja OPD
// @Description Sinkron data penetapan renja OPD berdasarkan kode OPD dan tahun dari perencanaan
// @Tags        Sync
// @Accept      json
// @Produce     json
//
// @Param       payload body web.SyncPenetapanOpdRequest true "Payload sinkronisasi penetapan OPD"
//
// @Success     200 {object} web.Response[web.SyncPenetapanOpdResponse] "Success"
// @Failure     422 {object} web.ValidationErrorResponse                "Unprocessable Entity"
// @Failure     500 {object} web.ErrorResponse                          "Internal Server Error"
//
// @Router      /opd/renja/sync [post]
func (app *Application) SyncPenetapanRenjaOpdHandler(
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
		KodeOpd: input.KodeOpd,
		Tahun:   input.Tahun,
	}
	v := validator.New()
	if web.ValidateSyncPenetapanRequest(v, request); !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	result, err := app.PenetapanOpdService.SyncPenetapanOpd(
		r.Context(),
		request,
		domain.JenisPenetapanRenja,
	)
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

// RenjaOpdHandler godoc
//
// @Summary     Get sasaran OPD penetapan
// @Description Mengambil data renja OPD berdasarkan kode OPD dan tahun penetapan
// @Tags        OPD
// @Accept      json
// @Produce     json
//
// @Param       kodeOpd query string true "Kode OPD"
// @Param       tahun   query int    true "Tahun Penetapan"
//
// @Success     200 {object} web.Response[web.RenjaPenetapanOpdResponse]   "Berhasil mengambil data renja OPD"
// @Failure     400 {object} web.ValidationErrorResponse                   "Bad Request"
// @Failure     500 {object} web.ErrorResponse                             "Internal Server Error"
//
// @Router      /opd/renja [get]
func (app *Application) RenjaOpdHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	query := r.URL.Query()

	errors := map[string]string{}

	kodeOpd := query.Get("kodeOpd")

	tahunStr := query.Get("tahun")

	var tahun int

	if kodeOpd == "" {
		errors["kodeOpd"] = "required"
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

	result, err := app.PenetapanOpdService.FindRenja(r.Context(), request)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	response := web.Response[web.RenjaPenetapanOpdResponse]{
		Data: result,
	}

	err = app.WriteJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}
