package metrics

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// PoolMetrics holds Prometheus metrics for a single Portworx storage pool.
//
// Units:
//   - Bytes fields      – raw bytes
//   - Throughput fields – bytes / second
//   - IOPS fields       – operations / second
//   - Latency fields    – milliseconds (seconds converted ×1000)
type PoolMetrics struct {
	PoolID          string  `json:"pool_id"`
	ReadThroughput  float64 `json:"read_throughput_bytes_s"`
	WriteThroughput float64 `json:"write_throughput_bytes_s"`
	ReadIOPS        float64 `json:"read_iops"`
	WriteIOPS       float64 `json:"write_iops"`
	ReadLatencyMs   float64 `json:"read_latency_ms"`
	WriteLatencyMs  float64 `json:"write_latency_ms"`
	TotalBytes      float64 `json:"total_bytes"`
	UsedBytes       float64 `json:"used_bytes"`
	AvailableBytes  float64 `json:"available_bytes"`
	Status          float64 `json:"status"` // 1 = online, 0 = degraded/offline
}

// NodeMetrics holds Prometheus metrics for a Portworx node.
//
// Node-level I/O metrics use the target="1" label variant (actual physical
// disk activity at the storage-target role).
//
// Units:
//   - Bytes fields      – raw bytes
//   - Throughput fields – bytes / second
//   - IOPS fields       – operations / second
//   - Latency fields    – milliseconds
//   - CPU               – percentage (0–100)
//   - Memory fields     – raw bytes
type NodeMetrics struct {
	Error   bool   `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	NodeID  string `json:"node_id"`

	// Node-level I/O (target role = actual disk activity)
	ReadThroughput  float64 `json:"read_throughput_bytes_s"`
	WriteThroughput float64 `json:"write_throughput_bytes_s"`
	ReadIOPS        float64 `json:"read_iops"`
	WriteIOPS       float64 `json:"write_iops"`
	ReadLatencyMs   float64 `json:"read_latency_ms"`
	WriteLatencyMs  float64 `json:"write_latency_ms"`

	// System resources
	CPUPercent  float64 `json:"cpu_percent"`
	TotalMemory float64 `json:"total_memory_bytes"`
	UsedMemory  float64 `json:"used_memory_bytes"`
	FreeMemory  float64 `json:"free_memory_bytes"`
	NumVolumes  float64 `json:"num_volumes"`

	// Per-pool breakdown (keyed by pool UUID)
	StoragePools map[string]*PoolMetrics `json:"storage_pools"`
}

// NodeMetricsHistory holds time-series I/O data for a node over the last
// ~10 minutes for dashboard chart pre-population.
// Timestamps are Unix milliseconds; all value arrays are index-aligned.
type NodeMetricsHistory struct {
	NodeID          string    `json:"node_id"`
	Timestamps      []int64   `json:"timestamps"`
	ReadThroughput  []float64 `json:"read_throughput_bytes_s"`
	WriteThroughput []float64 `json:"write_throughput_bytes_s"`
	ReadIOPS        []float64 `json:"read_iops"`
	WriteIOPS       []float64 `json:"write_iops"`
	ReadLatencyMs   []float64 `json:"read_latency_ms"`
	WriteLatencyMs  []float64 `json:"write_latency_ms"`
}

// nodeMetricRegex covers all px_node_stats_* and px_pool_stats_* metrics
// fetched in a single instant query for a given node.
const nodeMetricRegex = `px_node_stats_readthroughput|px_node_stats_writethroughput|` +
	`px_node_stats_read_iops|px_node_stats_write_iops|` +
	`px_node_stats_read_latency_seconds|px_node_stats_write_latency_seconds|` +
	`px_node_stats_cpu_usage|px_node_stats_total_mem|px_node_stats_used_mem|` +
	`px_node_stats_free_mem|px_node_stats_num_volumes|` +
	`px_pool_stats_readthroughput|px_pool_stats_writethroughput|` +
	`px_pool_stats_read_iops|px_pool_stats_write_iops|` +
	`px_pool_stats_read_latency_seconds|px_pool_stats_write_latency_seconds|` +
	`px_pool_stats_total_bytes|px_pool_stats_used_bytes|px_pool_stats_available_bytes|` +
	`px_pool_stats_status|px_pool_stats_pool_status`

// GetNodeMetrics fetches current Prometheus metrics for the named Portworx node
// from the Thanos Querier.
//
// metricsURL must be the Thanos HTTP API base URL, e.g.:
//
//	https://thanos-querier.openshift-monitoring.svc.cluster.local:9091
//
// nodeID must match the nodeID label in the Prometheus output
// (Portworx UUID, e.g. "6784f98f-6f71-4e55-bbf0-c12bc1ae659e").
func GetNodeMetrics(metricsURL, token, nodeID string) (*NodeMetrics, error) {
	if metricsURL == "" {
		return nil, fmt.Errorf("metrics URL is not configured – set PORTWORX_METRICS_URL")
	}

	// One instant query fetches all node + pool metrics at once.
	sel := `nodeID="` + nodeID + `"`
	promQL := `{__name__=~"` + nodeMetricRegex + `",` + sel + `}`

	results, err := instantQuery(metricsURL, token, promQL)
	if err != nil {
		return nil, err
	}

	nm := &NodeMetrics{
		NodeID:       nodeID,
		StoragePools: make(map[string]*PoolMetrics),
	}

	for _, r := range results {
		name := r.Metric["__name__"]
		if name == "" || len(r.Value) < 2 {
			continue
		}
		val, ok := parseScalar(r.Value[1])
		if !ok {
			continue
		}

		labels := r.Metric
		isTarget := labels["target"] == "1"
		isCoord := labels["coordinator"] == "1"
		poolID := uuidLabel(labels, "pool")
		hasIntPool := isIntegerPool(labels)

		// Skip coordinator-role and integer-pool-index lines (same as before).
		if isCoord || hasIntPool {
			continue
		}

		switch {
		case isTarget:
			// Node-level I/O metrics (target role).
			switch name {
			case "px_node_stats_readthroughput":
				nm.ReadThroughput = val
			case "px_node_stats_writethroughput":
				nm.WriteThroughput = val
			case "px_node_stats_read_iops":
				nm.ReadIOPS = val
			case "px_node_stats_write_iops":
				nm.WriteIOPS = val
			case "px_node_stats_read_latency_seconds":
				nm.ReadLatencyMs = val * 1000
			case "px_node_stats_write_latency_seconds":
				nm.WriteLatencyMs = val * 1000
			}

		case poolID != "":
			// Per-pool metrics (UUID pool label).
			pool := getOrCreatePool(nm, poolID)
			switch name {
			case "px_pool_stats_readthroughput":
				pool.ReadThroughput = val
			case "px_pool_stats_writethroughput":
				pool.WriteThroughput = val
			case "px_pool_stats_read_iops":
				pool.ReadIOPS = val
			case "px_pool_stats_write_iops":
				pool.WriteIOPS = val
			case "px_pool_stats_read_latency_seconds":
				pool.ReadLatencyMs = val * 1000
			case "px_pool_stats_write_latency_seconds":
				pool.WriteLatencyMs = val * 1000
			case "px_pool_stats_total_bytes":
				pool.TotalBytes = val
			case "px_pool_stats_used_bytes":
				pool.UsedBytes = val
			case "px_pool_stats_available_bytes":
				pool.AvailableBytes = val
			case "px_pool_stats_status", "px_pool_stats_pool_status":
				if pool.Status == 0 {
					pool.Status = val
				}
			}

		default:
			// Node-level gauges (no secondary label).
			switch name {
			case "px_node_stats_cpu_usage":
				nm.CPUPercent = val
			case "px_node_stats_total_mem":
				nm.TotalMemory = val
			case "px_node_stats_used_mem":
				nm.UsedMemory = val
			case "px_node_stats_free_mem":
				nm.FreeMemory = val
			case "px_node_stats_num_volumes":
				nm.NumVolumes = val
			}
		}
	}

	return nm, nil
}

// GetNodeMetricsHistory fetches the last ~10 minutes of I/O time-series data
// for the named node from the Thanos Querier. Only node-level chart metrics are
// included (per-pool data is snapshot-only via GetNodeMetrics).
func GetNodeMetricsHistory(metricsURL, token, nodeID string) (*NodeMetricsHistory, error) {
	if metricsURL == "" {
		return nil, fmt.Errorf("metrics URL is not configured – set PORTWORX_METRICS_URL")
	}

	now := time.Now()
	start := now.Add(-historyWindowSec * time.Second)
	sel := `nodeID="` + nodeID + `",target="1"`

	type rangeSpec struct {
		promQL  string
		scaleFn func(float64) float64
	}

	queries := []rangeSpec{
		{promQL: `sum(px_node_stats_readthroughput{` + sel + `})`},
		{promQL: `sum(px_node_stats_writethroughput{` + sel + `})`},
		{promQL: `sum(px_node_stats_read_iops{` + sel + `})`},
		{promQL: `sum(px_node_stats_write_iops{` + sel + `})`},
		{
			promQL:  `sum(px_node_stats_read_latency_seconds{` + sel + `})`,
			scaleFn: func(v float64) float64 { return v * 1000 },
		},
		{
			promQL:  `sum(px_node_stats_write_latency_seconds{` + sel + `})`,
			scaleFn: func(v float64) float64 { return v * 1000 },
		},
	}

	type seriesResult struct {
		idx  int
		ts   []int64
		vals []float64
	}

	ch := make(chan seriesResult, len(queries))
	var wg sync.WaitGroup

	for i, q := range queries {
		wg.Add(1)
		go func(idx int, spec rangeSpec) {
			defer wg.Done()
			results, err := rangeQuery(metricsURL, token, spec.promQL, start, now, historyStepSec)
			if err != nil {
				ch <- seriesResult{idx: idx}
				return
			}
			ts, vals := extractTimeSeries(results)
			if spec.scaleFn != nil {
				for j, v := range vals {
					vals[j] = spec.scaleFn(v)
				}
			}
			ch <- seriesResult{idx: idx, ts: ts, vals: vals}
		}(i, q)
	}

	wg.Wait()
	close(ch)

	seriesMap := make(map[int]seriesResult, len(queries))
	for sr := range ch {
		seriesMap[sr.idx] = sr
	}

	var refTS []int64
	for i := 0; i < len(queries); i++ {
		if sr := seriesMap[i]; len(sr.ts) > 0 {
			refTS = sr.ts
			break
		}
	}
	if refTS == nil {
		return &NodeMetricsHistory{NodeID: nodeID}, nil
	}

	align := func(idx int) []float64 {
		sr := seriesMap[idx]
		if len(sr.ts) == 0 {
			return make([]float64, len(refTS))
		}
		return alignToTimestamps(refTS, buildTSMap(sr.ts, sr.vals))
	}

	return &NodeMetricsHistory{
		NodeID:          nodeID,
		Timestamps:      refTS,
		ReadThroughput:  align(0),
		WriteThroughput: align(1),
		ReadIOPS:        align(2),
		WriteIOPS:       align(3),
		ReadLatencyMs:   align(4),
		WriteLatencyMs:  align(5),
	}, nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

// uuidLabel returns the value of the given label key only when it looks like a
// UUID (contains hyphens and is longer than 8 chars).
func uuidLabel(labels map[string]string, key string) string {
	val := labels[key]
	if strings.Contains(val, "-") && len(val) > 8 {
		return val
	}
	return ""
}

// isIntegerPool returns true when the "pool" label contains an integer value
// (pool index) rather than a UUID.
func isIntegerPool(labels map[string]string) bool {
	val := labels["pool"]
	if val == "" {
		return false
	}
	return !strings.Contains(val, "-")
}

// getOrCreatePool returns the PoolMetrics for poolID, creating it if absent.
func getOrCreatePool(nm *NodeMetrics, poolID string) *PoolMetrics {
	if p, ok := nm.StoragePools[poolID]; ok {
		return p
	}
	p := &PoolMetrics{PoolID: poolID}
	nm.StoragePools[poolID] = p
	return p
}
