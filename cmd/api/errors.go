package main

import (
	"fmt"
	"net/http"
)

// The logError() method is a helper for logging an error message
func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

// The errorResponse() method is a helper for sending JSON-formatted error message
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	res := envelope{"error": message}

	err := app.writeJSON(w, status, res, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// the serverErrorResponse() method will be used to send a 500 Server Error
// status code and JSON response to the client
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process the request"

	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// the forbiddenResponse() method will be used to send a 403 Forbidden
// status code and JSON response to the client
func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	message := "you do not have permission to access this resource"

	app.errorResponse(w, r, http.StatusForbidden, message)
}

// the notFoundResponse() method will be used to send a 404 Not Found
// status code and JSON response to the client
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"

	app.errorResponse(w, r, http.StatusNotFound, message)
}

// the methodNotAllowedResponse() method will be used to send a 405 Method Not Allowed
// status code and JSON response to the client
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("%s method is not supported for this resource", r.Method)

	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}
