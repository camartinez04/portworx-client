package main

import (
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// VolumeMetricsAPIHTTP is a thin proxy that forwards live-metrics requests from
// the browser to the broker and relays the JSON response back.
//
// Route: GET /portworx/client/api/metrics/{volume_name}
//
// The handler sits behind the AuthKeycloak middleware so the browser session is
// already validated when this runs.  The broker call is intentionally plain HTTP
// (same internal-cluster communication pattern used by all other repository
// functions in this service).
func (m *Repository) VolumeMetricsAPIHTTP(w http.ResponseWriter, r *http.Request) {

	volumeName := chi.URLParam(r, "volume_name")
	if volumeName == "" {
		http.Error(w, `{"error":true,"message":"volume name is required"}`, http.StatusBadRequest)
		return
	}

	targetURL := BrokerURL + "/broker/getvolumemetrics/" + volumeName

	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		log.Printf("VolumeMetricsAPIHTTP: build request: %v", err)
		http.Error(w, `{"error":true,"message":"internal error"}`, http.StatusInternalServerError)
		return
	}

	token := m.App.Session.GetString(r.Context(), "token")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("VolumeMetricsAPIHTTP: broker call failed: %v", err)
		http.Error(w, `{"error":true,"message":"metrics broker unavailable"}`, http.StatusServiceUnavailable)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("VolumeMetricsAPIHTTP: read body: %v", err)
		http.Error(w, `{"error":true,"message":"failed to read metrics response"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
