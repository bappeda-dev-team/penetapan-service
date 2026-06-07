package api

import (
	"net/http"
	"strconv"

	"github.com/bappeda-dev-team/penetapan-service/internal/model/web"
	"github.com/bappeda-dev-team/penetapan-service/internal/model/web/individu"
)

// RekinIndividuHandler godoc
//
// @Summary     Get rencana kinerja individu penetapan
// @Description Mengambil data rencana kinerja individu berdasarkan id_pegawai, kode OPD dan tahun penetapan
// @Tags        Individu
// @Accept      json
// @Produce     json
//
// @Param       idPegawai query string true "Id Pegawai"
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

	idPegawai := query.Get("idPegawai")

	kodeOpd := query.Get("kodeOpd")

	tahunStr := query.Get("tahun")

	var tahun int

	if idPegawai == "" {
		errors["idPegawai"] = "required"
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

	request := individu.SyncPenetapanIndividuRequest{
		IdPegawai: idPegawai,
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

	response := web.Response[individu.RekinPenetapanIndividuResponse]{
		Data: result,
	}

	err = app.WriteJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}
