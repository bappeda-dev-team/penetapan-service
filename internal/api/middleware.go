package api

import (
	"fmt"
	"net/http"
)

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// built-in function recover() is used to check
			// if a panic occured
			pv := recover()

			if pv != nil {
				w.Header().Set("Connection", "close")

				app.ServerErrorResponse(w, r, fmt.Errorf("%v", pv))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
