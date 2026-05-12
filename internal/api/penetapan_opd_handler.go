package api

import (
	"net/http"
	"strconv"

	"github.com/bappeda-dev-team/penetapan-service/internal/model/domain"
	"github.com/bappeda-dev-team/penetapan-service/internal/model/web"
)

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
	// query
	query := r.URL.Query()
	kodeOpd := query.Get("kodeOpd")
	tahun := query.Get("tahun")
	request := domain.PenetapanOpdRequest{
		KodeOpd: kodeOpd,
		Tahun:   tahun,
	}

	errors := map[string]string{}

	if request.KodeOpd == "" {
		errors["kodeOpd"] = "required"
	}

	if request.Tahun == "" {
		errors["tahun"] = "required"
	} else {
		_, err := strconv.Atoi(request.Tahun)
		if err != nil {
			errors["tahun"] = "tahun tidak valid"
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

	result, err := app.PenetapanOpdService.FindTujuan(r.Context(), request)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
	}

	response := web.Response[[]web.TujuanPenetapanOpdResponse]{
		Data: result,
	}
	err = app.WriteJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
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
	kodeOpd := query.Get("kodeOpd")
	tahun := query.Get("tahun")
	request := domain.PenetapanOpdRequest{
		KodeOpd: kodeOpd,
		Tahun:   tahun,
	}

	errors := map[string]string{}

	if request.KodeOpd == "" {
		errors["kode_opd"] = "required"
	}

	if request.Tahun == "" {
		errors["tahun"] = "required"
	} else {
		_, err := strconv.Atoi(request.Tahun)
		if err != nil {
			errors["tahun"] = "tahun tidak valid"
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
