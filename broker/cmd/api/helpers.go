package main

import (
	"encoding/json"
	"net/http"
)

//writeJSON marshals the provided interface into JSON and writes it to the response.
func writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {

	out, err := json.Marshal(data)

	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil

}

//errorJSON checks for errors on JSON
func (app *AppConfig) errorJSON(w http.ResponseWriter, err error, status ...int) error {

	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return writeJSON(w, statusCode, payload)

}

func ServerError(w http.ResponseWriter, err error) {

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}
