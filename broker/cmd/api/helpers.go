package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// writeJSON marshals the provided interface into JSON and writes it to the response.
func writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Set the content type and other headers.
	w.Header().Set("Content-Type", "application/json")
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	// Commit headers and write the status code.
	// Ensure this is the only place where WriteHeader is called.
	w.WriteHeader(status)

	// Write the body after setting headers and status.
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// errorJSON checks for errors on JSON
func (App *AppConfig) errorJSON(w http.ResponseWriter, err error, status ...int) error {

	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return writeJSON(w, statusCode, payload)

}

// ServerError writes a 500 error to the response
func ServerError(w http.ResponseWriter, err error) {

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}

// DateFormat returns date from UNIX timestamp
func DateFormat(date int64) string {

	times := time.Unix(date, 0)

	layout := "2006-01-02 15:04:05"

	return times.Format(layout)
}
