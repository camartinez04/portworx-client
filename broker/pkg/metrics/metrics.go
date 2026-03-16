package metrics

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// VolumeMetrics holds parsed Prometheus metrics for a specific Portworx volume.
//
// Units:
//   - Bytes fields     – raw bytes
//   - Throughput fields – bytes / second  (from px_volume_readthroughput / px_volume_writethroughput)
//   - IOPS fields      – operations / second
//   - Latency fields   – milliseconds (internally seconds are converted)
//   - IODepth          – average queue depth (gauge)
type VolumeMetrics struct {
	Error   bool   `json:"error,omitempty"`
	Message string `json:"message,omitempty"`

	VolumeName string `json:"volume_name"`

	// Cumulative byte counters (useful for rate calculation client-side)
	ReadBytes  float64 `json:"read_bytes"`
	WriteBytes float64 `json:"write_bytes"` // px_volume_vol_written_bytes

	// Direct throughput gauges (bytes/s, pre-computed by Portworx)
	ReadThroughput  float64 `json:"read_throughput_bytes_s"`
	WriteThroughput float64 `json:"write_throughput_bytes_s"`

	// IOPS gauges
	ReadIOPS  float64 `json:"read_iops"`
	WriteIOPS float64 `json:"write_iops"`
	IOPS      float64 `json:"iops"` // combined

	// Latency in milliseconds
	ReadLatencyMs  float64 `json:"read_latency_ms"`
	WriteLatencyMs float64 `json:"write_latency_ms"`

	// Queue depth
	IODepth float64 `json:"io_depth"`

	// Capacity / usage
	CapacityBytes float64 `json:"capacity_bytes"`
	UsageBytes    float64 `json:"usage_bytes"`
}

// GetVolumeMetrics fetches Prometheus metrics for the requested volume.
//
// metricsURLs may be:
//   - empty string  → auto-discover all Portworx pod endpoints via the K8s API
//   - single URL    → fetch that one endpoint
//   - comma-separated list → fan-out across all listed endpoints (local dev)
//
// When multiple endpoints are used the function fans out concurrently and
// returns the result with the highest I/O activity for the volume.
func GetVolumeMetrics(metricsURLs, volumeName string) (*VolumeMetrics, error) {
	urls, err := ResolveMetricsURLs(metricsURLs)
	if err != nil {
		return nil, err
	}

	if len(urls) == 1 {
		return fetchVolumeMetrics(urls[0], volumeName)
	}

	// Fan-out across all endpoints concurrently.
	type result struct {
		vm  *VolumeMetrics
		err error
	}
	ch := make(chan result, len(urls))
	for _, u := range urls {
		go func(url string) {
			vm, err := fetchVolumeMetrics(url, volumeName)
			ch <- result{vm, err}
		}(u)
	}

	var best *VolumeMetrics
	var lastErr error
	for range urls {
		r := <-ch
		if r.err != nil {
			lastErr = r.err
			continue
		}
		if best == nil || ioScore(r.vm) > ioScore(best) {
			best = r.vm
		}
	}

	if best == nil {
		return nil, lastErr
	}
	return best, nil
}

// ioScore returns a scalar representing total I/O activity for a volume result.
// Used to pick the "most active" node result when fanning out.
func ioScore(vm *VolumeMetrics) float64 {
	return vm.ReadThroughput + vm.WriteThroughput + vm.ReadIOPS + vm.WriteIOPS +
		vm.ReadBytes + vm.WriteBytes
}

// fetchVolumeMetrics fetches a single metrics endpoint and parses volume data.
func fetchVolumeMetrics(metricsURL, volumeName string) (*VolumeMetrics, error) {
	resp, err := http.Get(metricsURL) //nolint:gosec // URL comes from operator-configured env var
	if err != nil {
		return nil, fmt.Errorf("failed to reach metrics endpoint %s: %w", metricsURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("metrics endpoint returned HTTP %d", resp.StatusCode)
	}

	vm := &VolumeMetrics{VolumeName: volumeName}

	// Label selector we look for in each line.
	labelSelector := `volumename="` + volumeName + `"`

	scanner := bufio.NewScanner(resp.Body)
	// 2 MiB scan buffer – Portworx metrics lines can be very long (many labels).
	scanner.Buffer(make([]byte, 2*1024*1024), 2*1024*1024)

	for scanner.Scan() {
		line := scanner.Text()

		// Skip HELP/TYPE comments and blank lines.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Skip lines that don't mention our volume.
		if !strings.Contains(line, labelSelector) {
			continue
		}

		// Prometheus exposition format:
		//   metricname{label="val",...} value [unix_timestamp_ms]
		braceIdx := strings.Index(line, "{")
		if braceIdx < 0 {
			continue
		}
		metricName := line[:braceIdx]

		closeBraceIdx := strings.LastIndex(line, "}")
		if closeBraceIdx < 0 || closeBraceIdx+1 >= len(line) {
			continue
		}

		// Take only the first token after "} " – the numeric value.
		rest := strings.TrimSpace(line[closeBraceIdx+1:])
		parts := strings.Fields(rest)
		if len(parts) == 0 {
			continue
		}
		value, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			continue
		}

		// ── Map Portworx metric names to struct fields ──────────────────────
		switch metricName {

		// ── Cumulative byte counters ─────────────────────────────────────────
		case "px_volume_vol_read_bytes", "px_volume_read_bytes":
			vm.ReadBytes = value
		case "px_volume_vol_written_bytes", "px_volume_written_bytes":
			vm.WriteBytes = value

		// ── Direct throughput gauges (bytes/s, pre-computed) ─────────────────
		case "px_volume_readthroughput":
			vm.ReadThroughput = value
		case "px_volume_writethroughput":
			vm.WriteThroughput = value

		// ── IOPS gauges ──────────────────────────────────────────────────────
		case "px_volume_read_iops":
			vm.ReadIOPS = value
		case "px_volume_write_iops":
			vm.WriteIOPS = value
		case "px_volume_iops":
			vm.IOPS = value

		// ── Latency (Portworx reports in seconds → convert to ms) ────────────
		case "px_volume_vol_read_latency_seconds", "px_volume_read_latency_seconds":
			vm.ReadLatencyMs = value * 1000
		case "px_volume_vol_write_latency_seconds", "px_volume_write_latency_seconds":
			vm.WriteLatencyMs = value * 1000

		// ── IO queue depth ───────────────────────────────────────────────────
		case "px_volume_depth_io":
			vm.IODepth = value

		// ── Capacity / usage ─────────────────────────────────────────────────
		case "px_volume_capacity_bytes", "px_volume_fs_capacity_bytes":
			if vm.CapacityBytes == 0 {
				vm.CapacityBytes = value
			}
		case "px_volume_usage_bytes", "px_volume_fs_usage_bytes":
			if vm.UsageBytes == 0 {
				vm.UsageBytes = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading metrics response: %w", err)
	}

	return vm, nil
}

// splitURLs splits a comma-separated URL list, trimming whitespace and skipping blanks.
func splitURLs(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			out = append(out, s)
		}
	}
	return out
}
