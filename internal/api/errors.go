package api

import (
	"fmt"
	"github.com/bappeda-dev-team/penetapan-service/internal/model/web"
	"net/http"
)

// The logError() method is a helper for logging an error message
func (app *Application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	app.Logger.Error(err.Error(), "method", method, "uri", uri)
}

// The errorResponse() method is a helper for sending JSON-formatted error message
func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	res := web.ErrorResponse{
		Error: message,
	}

	err := app.WriteJSON(w, status, res, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// the serverErrorResponse() method will be used to send a 500 Server Error
// status code and JSON response to the client
func (app *Application) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process the request"

	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// the forbiddenResponse() method will be used to send a 403 Forbidden
// status code and JSON response to the client
func (app *Application) ForbiddenResponse(w http.ResponseWriter, r *http.Request) {
	message := "you do not have permission to access this resource"

	app.errorResponse(w, r, http.StatusForbidden, message)
}

// the notFoundResponse() method will be used to send a 404 Not Found
// status code and JSON response to the client
func (app *Application) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"

	app.errorResponse(w, r, http.StatusNotFound, message)
}

// the methodNotAllowedResponse() method will be used to send a 405 Method Not Allowed
// status code and JSON response to the client
func (app *Application) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("%s method is not supported for this resource", r.Method)

	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// the badRequestResponse() method will be used to send a 400 Bad Request
// status code and JSON response to the client
func (app *Application) BadRequestResponse(w http.ResponseWriter, r *http.Request, message any) {
	status := http.StatusBadRequest
	err := app.WriteJSON(w, status, message, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// the FailedValidationResponse() method will be used to send a 400 Bad Request
// for invalid form submit
// from post request
func (app *Application) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
