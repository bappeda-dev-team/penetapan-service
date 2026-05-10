package api

import (
	"encoding/json"
	"net/http"
)

func (app *Application) WriteJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {

	js, err := json.MarshalIndent(
		data,
		"",
		"\t",
	)

	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, values := range headers {
		for _, value := range values {
			w.Header().Set(key, value)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	w.Write(js)

	return nil
}
