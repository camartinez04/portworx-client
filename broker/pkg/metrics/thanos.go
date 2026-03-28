package metrics

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// instantResult is one time series from a Thanos/Prometheus instant-query response.
type instantResult struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"` // [float64 timestamp, string value]
}

// rangeResult is one time series from a Thanos/Prometheus range-query response.
type rangeResult struct {
	Metric map[string]string `json:"metric"`
	Values [][]interface{}   `json:"values"` // [[float64 ts, string val], ...]
}

// thanosHTTPClient is shared across all Thanos calls.
// TLS verification is skipped for cluster-internal HTTPS endpoints where the
// cluster CA is not necessarily present in the pod's trust store.
var thanosHTTPClient = &http.Client{ //nolint:gochecknoglobals
	Timeout: 15 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // cluster-internal endpoint
	},
}

// instantQuery runs a PromQL instant query against the Thanos HTTP API and
// returns all result series.
func instantQuery(baseURL, token, promQL string) ([]instantResult, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		baseURL+"/api/v1/query?"+url.Values{"query": {promQL}}.Encode(),
		nil,
	)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := thanosHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("thanos query: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("thanos HTTP %d: %.300s", resp.StatusCode, body)
	}

	var envelope struct {
		Status string `json:"status"`
		Data   struct {
			Result []instantResult `json:"result"`
		} `json:"data"`
		Error string `json:"error,omitempty"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("thanos parse: %w", err)
	}
	if envelope.Status != "success" {
		return nil, fmt.Errorf("thanos: %s", envelope.Error)
	}
	return envelope.Data.Result, nil
}

// rangeQuery runs a PromQL range query against the Thanos HTTP API and returns
// all result series. stepSec is the resolution step in seconds.
func rangeQuery(baseURL, token, promQL string, start, end time.Time, stepSec int) ([]rangeResult, error) {
	params := url.Values{
		"query": {promQL},
		"start": {fmt.Sprintf("%d", start.Unix())},
		"end":   {fmt.Sprintf("%d", end.Unix())},
		"step":  {fmt.Sprintf("%ds", stepSec)},
	}
	req, err := http.NewRequest(http.MethodGet, baseURL+"/api/v1/query_range?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := thanosHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("thanos range query: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("thanos HTTP %d: %.300s", resp.StatusCode, body)
	}

	var envelope struct {
		Status string `json:"status"`
		Data   struct {
			Result []rangeResult `json:"result"`
		} `json:"data"`
		Error string `json:"error,omitempty"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("thanos range parse: %w", err)
	}
	if envelope.Status != "success" {
		return nil, fmt.Errorf("thanos range: %s", envelope.Error)
	}
	return envelope.Data.Result, nil
}

// parseScalar converts a Prometheus scalar value (returned as a JSON string or
// float64) to float64. Returns (0, false) for NaN, ±Inf, or parse errors.
func parseScalar(v interface{}) (float64, bool) {
	switch s := v.(type) {
	case string:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil || f != f { // f != f catches NaN
			return 0, false
		}
		return f, true
	case float64:
		return s, true
	}
	return 0, false
}

// sumInstantByName aggregates instant-query results: for each unique __name__
// it sums the values across all series (handles data replicated across nodes).
func sumInstantByName(results []instantResult) map[string]float64 {
	out := make(map[string]float64, len(results))
	for _, r := range results {
		name := r.Metric["__name__"]
		if name == "" || len(r.Value) < 2 {
			continue
		}
		if f, ok := parseScalar(r.Value[1]); ok {
			out[name] += f
		}
	}
	return out
}

// maxInstantByName is like sumInstantByName but takes the maximum value per
// metric name (used for capacity/usage fields that should not be summed).
func maxInstantByName(results []instantResult) map[string]float64 {
	out := make(map[string]float64)
	for _, r := range results {
		name := r.Metric["__name__"]
		if name == "" || len(r.Value) < 2 {
			continue
		}
		if f, ok := parseScalar(r.Value[1]); ok && f > out[name] {
			out[name] = f
		}
	}
	return out
}

// extractTimeSeries pulls timestamps (Unix ms) and values out of the first
// non-empty range-query result series.
func extractTimeSeries(results []rangeResult) (timestamps []int64, values []float64) {
	for _, r := range results {
		if len(r.Values) == 0 {
			continue
		}
		ts := make([]int64, 0, len(r.Values))
		vs := make([]float64, 0, len(r.Values))
		for _, pair := range r.Values {
			if len(pair) < 2 {
				continue
			}
			tsF, ok1 := parseScalar(pair[0])
			val, ok2 := parseScalar(pair[1])
			if ok1 && ok2 {
				ts = append(ts, int64(tsF*1000)) // seconds → ms
				vs = append(vs, val)
			}
		}
		if len(ts) > 0 {
			return ts, vs
		}
	}
	return nil, nil
}

// alignToTimestamps takes a reference timestamp slice and a map of
// timestamp→value, and returns a value slice aligned to refTS, filling
// missing timestamps with 0.
func alignToTimestamps(refTS []int64, byTS map[int64]float64) []float64 {
	out := make([]float64, len(refTS))
	for i, t := range refTS {
		out[i] = byTS[t]
	}
	return out
}

// buildTSMap converts parallel timestamp/value slices into a map for alignment.
func buildTSMap(timestamps []int64, values []float64) map[int64]float64 {
	m := make(map[int64]float64, len(timestamps))
	for i, t := range timestamps {
		m[t] = values[i]
	}
	return m
}
