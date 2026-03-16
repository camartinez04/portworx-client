package metrics

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// PoolMetrics holds parsed Prometheus metrics for a single Portworx storage pool.
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

// NodeMetrics holds parsed Prometheus metrics for a Portworx node.
//
// Node-level I/O metrics use the target="1" label variant – this represents
// the node acting as storage target (actual physical disk activity).
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

// GetNodeMetrics fetches Prometheus metrics for the requested Portworx node ID.
//
// metricsURLs may be:
//   - empty string  → auto-discover all Portworx pod endpoints via the K8s API
//   - single URL    → fetch that one endpoint
//   - comma-separated list → fan-out across all listed endpoints (local dev)
//
// nodeID must match the `nodeID` label in the Prometheus output
// (Portworx UUID, e.g. "6784f98f-6f71-4e55-bbf0-c12bc1ae659e").
func GetNodeMetrics(metricsURLs, nodeID string) (*NodeMetrics, error) {
	urls, err := ResolveMetricsURLs(metricsURLs)
	if err != nil {
		return nil, err
	}

	if len(urls) == 1 {
		return fetchNodeMetrics(urls[0], nodeID)
	}

	// Fan-out across all endpoints concurrently.
	type result struct {
		nm  *NodeMetrics
		err error
	}
	ch := make(chan result, len(urls))
	for _, u := range urls {
		go func(url string) {
			nm, err := fetchNodeMetrics(url, nodeID)
			ch <- result{nm, err}
		}(u)
	}

	var best *NodeMetrics
	var lastErr error
	for range urls {
		r := <-ch
		if r.err != nil {
			lastErr = r.err
			continue
		}
		if best == nil || nodeIOScore(r.nm) > nodeIOScore(best) {
			best = r.nm
		}
	}

	if best == nil {
		return nil, lastErr
	}
	return best, nil
}

// nodeIOScore returns a scalar representing total I/O activity for a node result.
func nodeIOScore(nm *NodeMetrics) float64 {
	return nm.ReadThroughput + nm.WriteThroughput + nm.ReadIOPS + nm.WriteIOPS +
		nm.CPUPercent + nm.NumVolumes
}

// fetchNodeMetrics fetches a single metrics endpoint and parses node data.
func fetchNodeMetrics(metricsURL, nodeID string) (*NodeMetrics, error) {
	resp, err := http.Get(metricsURL) //nolint:gosec // URL comes from operator-configured env var
	if err != nil {
		return nil, fmt.Errorf("failed to reach metrics endpoint %s: %w", metricsURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("metrics endpoint returned HTTP %d", resp.StatusCode)
	}

	nm := &NodeMetrics{
		NodeID:       nodeID,
		StoragePools: make(map[string]*PoolMetrics),
	}

	// Label selectors
	nodeIDSelector := `nodeID="` + nodeID + `"`

	scanner := bufio.NewScanner(resp.Body)
	// 4 MiB buffer – node metrics lines can be very long.
	scanner.Buffer(make([]byte, 4*1024*1024), 4*1024*1024)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Only process lines for our node.
		if !strings.Contains(line, nodeIDSelector) {
			continue
		}

		// Parse Prometheus exposition format:
		//   metricname{label="val",...} value [timestamp]
		braceIdx := strings.Index(line, "{")
		if braceIdx < 0 {
			continue
		}
		metricName := line[:braceIdx]

		closeBraceIdx := strings.LastIndex(line, "}")
		if closeBraceIdx < 0 || closeBraceIdx+1 >= len(line) {
			continue
		}

		labelsSection := line[braceIdx+1 : closeBraceIdx]

		rest := strings.TrimSpace(line[closeBraceIdx+1:])
		parts := strings.Fields(rest)
		if len(parts) == 0 {
			continue
		}
		value, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			continue
		}

		// ── Determine which label variant this line belongs to ────────────────
		isTarget      := strings.Contains(labelsSection, `target="1"`)
		isCoordinator := strings.Contains(labelsSection, `coordinator="1"`)
		poolUUID      := extractLabelUUID(labelsSection, "pool")
		hasPoolInt    := hasIntegerPool(labelsSection)

		// Skip coordinator-role and integer-pool-index lines.
		if isCoordinator || hasPoolInt {
			continue
		}

		switch {

		// ── Node-level I/O metrics (target role) ─────────────────────────────
		case isTarget:
			switch metricName {
			case "px_node_stats_readthroughput":
				nm.ReadThroughput = value
			case "px_node_stats_writethroughput":
				nm.WriteThroughput = value
			case "px_node_stats_read_iops":
				nm.ReadIOPS = value
			case "px_node_stats_write_iops":
				nm.WriteIOPS = value
			case "px_node_stats_read_latency_seconds":
				nm.ReadLatencyMs = value * 1000
			case "px_node_stats_write_latency_seconds":
				nm.WriteLatencyMs = value * 1000
			}

		// ── Per-pool metrics (pool UUID label) ───────────────────────────────
		case poolUUID != "":
			pool := getOrCreatePool(nm, poolUUID)
			switch metricName {
			case "px_pool_stats_readthroughput":
				pool.ReadThroughput = value
			case "px_pool_stats_writethroughput":
				pool.WriteThroughput = value
			case "px_pool_stats_read_iops":
				pool.ReadIOPS = value
			case "px_pool_stats_write_iops":
				pool.WriteIOPS = value
			case "px_pool_stats_read_latency_seconds":
				pool.ReadLatencyMs = value * 1000
			case "px_pool_stats_write_latency_seconds":
				pool.WriteLatencyMs = value * 1000
			case "px_pool_stats_total_bytes":
				pool.TotalBytes = value
			case "px_pool_stats_used_bytes":
				pool.UsedBytes = value
			case "px_pool_stats_available_bytes":
				pool.AvailableBytes = value
			case "px_pool_stats_status", "px_pool_stats_pool_status":
				if pool.Status == 0 {
					pool.Status = value
				}
			}

		// ── Node-level gauges (no secondary label) ───────────────────────────
		default:
			switch metricName {
			case "px_node_stats_cpu_usage":
				nm.CPUPercent = value
			case "px_node_stats_total_mem":
				nm.TotalMemory = value
			case "px_node_stats_used_mem":
				nm.UsedMemory = value
			case "px_node_stats_free_mem":
				nm.FreeMemory = value
			case "px_node_stats_num_volumes":
				nm.NumVolumes = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading metrics response: %w", err)
	}

	return nm, nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

// extractLabelUUID returns the value of `key` from a Prometheus labels string
// only when it looks like a UUID (contains hyphens and is longer than 8 chars).
func extractLabelUUID(labels, key string) string {
	search := key + `="`
	idx := strings.Index(labels, search)
	if idx < 0 {
		return ""
	}
	start := idx + len(search)
	end := strings.Index(labels[start:], `"`)
	if end < 0 {
		return ""
	}
	val := labels[start : start+end]
	// Accept only UUID-like values (contains at least one hyphen and len > 8).
	if strings.Contains(val, "-") && len(val) > 8 {
		return val
	}
	return ""
}

// hasIntegerPool returns true if the labels contain pool="<integer>" (pool index).
func hasIntegerPool(labels string) bool {
	search := `pool="`
	idx := strings.Index(labels, search)
	if idx < 0 {
		return false
	}
	start := idx + len(search)
	end := strings.Index(labels[start:], `"`)
	if end < 0 {
		return false
	}
	val := labels[start : start+end]
	// Integer pool index has no hyphens.
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
