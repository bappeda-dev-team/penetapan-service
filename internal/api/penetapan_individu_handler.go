package api

import (
	"net/http"
	"strconv"

	"github.com/bappeda-dev-team/penetapan-service/internal/individu/web"
	"github.com/bappeda-dev-team/penetapan-service/internal/validator"
)

// RekinIndividuHandler godoc
//
// @Summary     Get rencana kinerja individu penetapan
// @Description Mengambil data rencana kinerja individu berdasarkan id_pegawai, kode OPD dan tahun penetapan
// @Tags        Individu
// @Accept      json
// @Produce     json
//
// @Param       pegawaiId query string true "Id Pegawai"
// @Param       kodeOpd query string true "Kode OPD"
// @Param       tahun   query int    true "Tahun Penetapan"
//
// @Success     200 {object} web.Response[web.RekinIndividuResponse] 	"Berhasil mengambil data rekin individu"
// @Failure     400 {object} web.ValidationErrorResponse                "Bad Request"
// @Failure     500 {object} web.ErrorResponse                          "Internal Server Error"
//
// @Router      /individu/rekin [get]
func (app *Application) RekinIndividuHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	query := r.URL.Query()

	errors := map[string]string{}

	pegawaiId := query.Get("pegawaiId")

	kodeOpd := query.Get("kodeOpd")

	tahunStr := query.Get("tahun")

	var tahun int

	if pegawaiId == "" {
		errors["pegawaiId"] = "required"
	}

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

	request := web.SyncPenetapanRequest{
		PegawaiId: pegawaiId,
		KodeOpd:   kodeOpd,
		Tahun:     tahun,
	}

	result, err := app.IndividuService.FindRekinsIndividu(
		r.Context(),
		request,
	)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	response := web.Response[web.RekinPenetapanIndividuResponse]{
		Data: result,
	}

	err = app.WriteJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}

// SyncRekinIndividuHandler godoc
//
// @Summary     Sync Rekin individu
// @Description Sinkron data pk penetapan berdasarkan id_pegawai, kode OPD dan tahun penetapan
// @Tags        Individu
// @Accept      json
// @Produce     json
//
// @Param       payload body web.SyncPenetapanRequest true "Payload sync penetapan individu"
//
// @Success     200 {object} web.Response[web.SyncPenetapanResponse] 	"Success"
// @Failure     400 {object} web.ValidationErrorResponse                "Bad Request"
// @Failure     500 {object} web.ErrorResponse                          "Internal Server Error"
//
// @Router      /individu/rekin/sync [post]
func (app *Application) SyncRekinIndividuHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	input := web.SyncPenetapanRequest{}
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

	request := &web.SyncPenetapanRequest{
		PegawaiId: input.PegawaiId,
		KodeOpd:   input.KodeOpd,
		Tahun:     input.Tahun,
	}
	v := validator.New()
	if web.ValidateSyncPenetapanRequest(v, request); !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	result, err := app.IndividuService.SyncPenetapanPkIndividu(
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

// SyncRenjaIndividuHandler godoc
//
// @Summary     Sync Renja individu
// @Description Sinkron data pk renja individu berdasarkan id_pegawai, kode OPD dan tahun penetapan
// @Tags        Individu
// @Accept      json
// @Produce     json
//
// @Param       payload body web.SyncPenetapanRequest true "Payload sync penetapan individu"
//
// @Success     200 {object} web.Response[web.SyncPenetapanResponse] 	"Success"
// @Failure     400 {object} web.ValidationErrorResponse                "Bad Request"
// @Failure     500 {object} web.ErrorResponse                          "Internal Server Error"
//
// @Router      /individu/renja/sync [post]
func (app *Application) SyncRenjaIndividuHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	input := web.SyncPenetapanRequest{}
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

	request := &web.SyncPenetapanRequest{
		PegawaiId: input.PegawaiId,
		KodeOpd:   input.KodeOpd,
		Tahun:     input.Tahun,
	}
	v := validator.New()
	if web.ValidateSyncPenetapanRequest(v, request); !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	result, err := app.IndividuService.SyncRenjaIndividu(
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

// RenjaIndividuHandler godoc
//
// @Summary     Get rencana kinerja individu penetapan
// @Description Mengambil data renja individu berdasarkan id_pegawai, kode OPD dan tahun penetapan
// @Tags        Individu
// @Accept      json
// @Produce     json
//
// @Param       pegawaiId query string true "Id Pegawai"
// @Param       kodeOpd query string true "Kode OPD"
// @Param       tahun   query int    true "Tahun Penetapan"
//
// @Success     200 {object} web.Response[web.RenjaIndividuResponse] 	"Berhasil mengambil data rekin individu"
// @Failure     400 {object} web.ValidationErrorResponse                "Bad Request"
// @Failure     500 {object} web.ErrorResponse                          "Internal Server Error"
//
// @Router      /individu/renja [get]
func (app *Application) RenjaIndividuHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	query := r.URL.Query()

	errors := map[string]string{}

	pegawaiId := query.Get("pegawaiId")

	kodeOpd := query.Get("kodeOpd")

	tahunStr := query.Get("tahun")

	var tahun int

	if pegawaiId == "" {
		errors["pegawaiId"] = "required"
	}

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

	request := web.SyncPenetapanRequest{
		PegawaiId: pegawaiId,
		KodeOpd:   kodeOpd,
		Tahun:     tahun,
	}

	result, err := app.IndividuService.FindRenjaIndividu(
		r.Context(),
		request,
	)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	response := web.Response[web.RenjaIndividuResponse]{
		Data: result,
	}

	err = app.WriteJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}
