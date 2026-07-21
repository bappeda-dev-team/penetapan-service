package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/bappeda-dev-team/penetapan-service/docs"
)

func (app *Application) Routes() http.Handler {
	router := httprouter.New()

	// 404
	router.NotFound = http.HandlerFunc(app.NotFoundResponse)
	// 405
	router.MethodNotAllowed = http.HandlerFunc(app.MethodNotAllowedResponse)

	router.Handler(http.MethodGet, "/swagger/*filepath", httpSwagger.WrapHandler)

	router.HandlerFunc(http.MethodGet, "/v3/api-docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	router.HandlerFunc(http.MethodGet, "/healthcheck", app.HealthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/individu/rekin", app.RekinIndividuHandler)
	router.HandlerFunc(http.MethodPost, "/individu/rekin/sync", app.SyncRekinIndividuHandler)

	// opd
	// renaksi
	router.HandlerFunc(http.MethodGet, "/opd/renaksi", app.RenaksiOpdHandler)

	// renja
	router.HandlerFunc(http.MethodGet, "/opd/renja", app.RenjaOpdHandler)
	router.HandlerFunc(http.MethodPost, "/opd/renja/sync", app.SyncPenetapanRenjaOpdHandler)

	// sasaran opd
	router.HandlerFunc(http.MethodGet, "/opd/sasaran", app.SasaranOpdHandler)
	router.HandlerFunc(http.MethodPost, "/opd/sasaran/sync", app.SyncPenetapanSasaranOpdHandler)

	// tujuan opd
	router.HandlerFunc(http.MethodGet, "/opd/tujuan", app.TujuanOpdHandler)
	router.HandlerFunc(http.MethodPost, "/opd/tujuan/sync", app.SyncPenetapanTujuanOpdHandler)

	// pemda
	// tujuan pemda
	router.HandlerFunc(http.MethodPost, "/pemda/tujuan/sync", app.SyncPenetapanTujuanPemdaHandler)

	return app.recoverPanic(router)
}
