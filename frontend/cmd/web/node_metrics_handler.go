package main

import (
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// NodeMetricsAPIHTTP is a thin proxy that forwards live node-metrics requests
// from the browser to the broker and relays the JSON response back.
//
// Route: GET /portworx/client/api/node-metrics/{node_id}
//
// The node_id is the Portworx UUID (e.g. "6784f98f-6f71-4e55-bbf0-c12bc1ae659e")
// and must match the nodeID label in the Prometheus metrics endpoint.
func (m *Repository) NodeMetricsAPIHTTP(w http.ResponseWriter, r *http.Request) {

	nodeID := chi.URLParam(r, "node_id")
	if nodeID == "" {
		http.Error(w, `{"error":true,"message":"node ID is required"}`, http.StatusBadRequest)
		return
	}

	targetURL := BrokerURL + "/broker/getnodemetrics/" + nodeID

	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		log.Printf("NodeMetricsAPIHTTP: build request: %v", err)
		http.Error(w, `{"error":true,"message":"internal error"}`, http.StatusInternalServerError)
		return
	}

	// Do not set Authorization header — broker middleware injects its global
	// KeycloakToken.  An empty bearer (from the "token" session key that is
	// never written) causes the broker to return 405.

	// Do not follow redirects.
	client := &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("NodeMetricsAPIHTTP: broker call failed: %v", err)
		http.Error(w, `{"error":true,"message":"metrics broker unavailable"}`, http.StatusServiceUnavailable)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("NodeMetricsAPIHTTP: broker returned %d", res.StatusCode)
		http.Error(w, `{"error":true,"message":"metrics unavailable"}`, http.StatusServiceUnavailable)
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("NodeMetricsAPIHTTP: read body: %v", err)
		http.Error(w, `{"error":true,"message":"failed to read metrics response"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
