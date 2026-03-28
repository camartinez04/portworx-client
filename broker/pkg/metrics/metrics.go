package metrics

import (
	"fmt"
	"sync"
	"time"
)

// VolumeMetrics holds the current Prometheus metrics for a specific Portworx volume.
//
// Units:
//   - Bytes fields      – raw bytes
//   - Throughput fields – bytes / second  (from px_volume_readthroughput / px_volume_writethroughput)
//   - IOPS fields       – operations / second
//   - Latency fields    – milliseconds (Portworx reports seconds; converted here)
//   - IODepth           – average queue depth (gauge)
type VolumeMetrics struct {
	Error   bool   `json:"error,omitempty"`
	Message string `json:"message,omitempty"`

	VolumeName string `json:"volume_name"`

	// Cumulative byte counters
	ReadBytes  float64 `json:"read_bytes"`
	WriteBytes float64 `json:"write_bytes"`

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

// VolumeMetricsHistory holds time-series data for a volume over the last
// ~10 minutes, suitable for pre-populating dashboard charts on page load.
// Timestamps are Unix milliseconds; all value arrays are index-aligned with Timestamps.
type VolumeMetricsHistory struct {
	VolumeName      string    `json:"volume_name"`
	Timestamps      []int64   `json:"timestamps"`
	ReadThroughput  []float64 `json:"read_throughput_bytes_s"`
	WriteThroughput []float64 `json:"write_throughput_bytes_s"`
	ReadIOPS        []float64 `json:"read_iops"`
	WriteIOPS       []float64 `json:"write_iops"`
	ReadLatencyMs   []float64 `json:"read_latency_ms"`
	WriteLatencyMs  []float64 `json:"write_latency_ms"`
}

// volumeMetricNames is the set of Portworx volume metric names fetched from Thanos.
// Using a regex match across all of them in a single instant query is the most
// efficient approach (one HTTP round-trip per call).
const volumeMetricRegex = `px_volume_readthroughput|px_volume_writethroughput|` +
	`px_volume_read_iops|px_volume_write_iops|px_volume_iops|` +
	`px_volume_vol_read_latency_seconds|px_volume_read_latency_seconds|` +
	`px_volume_vol_write_latency_seconds|px_volume_write_latency_seconds|` +
	`px_volume_depth_io|` +
	`px_volume_capacity_bytes|px_volume_fs_capacity_bytes|` +
	`px_volume_usage_bytes|px_volume_fs_usage_bytes|` +
	`px_volume_vol_read_bytes|px_volume_read_bytes|` +
	`px_volume_vol_written_bytes|px_volume_written_bytes`

// GetVolumeMetrics fetches current Prometheus metrics for the named volume from
// the Thanos Querier.
//
// metricsURL must be the Thanos HTTP API base URL, e.g.:
//
//	https://thanos-querier.openshift-monitoring.svc.cluster.local:9091
//
// token is an optional Bearer token for Thanos authentication (leave empty if
// the endpoint is unauthenticated within the cluster).
func GetVolumeMetrics(metricsURL, token, volumeName string) (*VolumeMetrics, error) {
	if metricsURL == "" {
		return nil, fmt.Errorf("metrics URL is not configured – set PORTWORX_METRICS_URL")
	}

	// One instant query fetches all px_volume metrics for this volume at once.
	sel := `volumename="` + volumeName + `"`
	promQL := `{__name__=~"` + volumeMetricRegex + `",` + sel + `}`

	results, err := instantQuery(metricsURL, token, promQL)
	if err != nil {
		return nil, err
	}

	// Aggregate: sum values per metric name across all series (handles data
	// replicated across multiple Portworx nodes in Thanos).
	sums := sumInstantByName(results)
	// For capacity/usage we want max (same physical value on all replicas).
	maxes := maxInstantByName(results)

	vm := &VolumeMetrics{VolumeName: volumeName}

	vm.ReadThroughput = sums["px_volume_readthroughput"]
	vm.WriteThroughput = sums["px_volume_writethroughput"]
	vm.ReadIOPS = sums["px_volume_read_iops"]
	vm.WriteIOPS = sums["px_volume_write_iops"]
	vm.IOPS = sums["px_volume_iops"]
	vm.IODepth = sums["px_volume_depth_io"]

	// Latency: Portworx reports seconds; convert to ms.
	latR := sums["px_volume_vol_read_latency_seconds"]
	if latR == 0 {
		latR = sums["px_volume_read_latency_seconds"]
	}
	vm.ReadLatencyMs = latR * 1000

	latW := sums["px_volume_vol_write_latency_seconds"]
	if latW == 0 {
		latW = sums["px_volume_write_latency_seconds"]
	}
	vm.WriteLatencyMs = latW * 1000

	// Capacity / usage: prefer the plain name; fall back to fs_ variant.
	cap := maxes["px_volume_capacity_bytes"]
	if cap == 0 {
		cap = maxes["px_volume_fs_capacity_bytes"]
	}
	vm.CapacityBytes = cap

	usage := maxes["px_volume_usage_bytes"]
	if usage == 0 {
		usage = maxes["px_volume_fs_usage_bytes"]
	}
	vm.UsageBytes = usage

	// Cumulative byte counters.
	rb := sums["px_volume_vol_read_bytes"]
	if rb == 0 {
		rb = sums["px_volume_read_bytes"]
	}
	vm.ReadBytes = rb

	wb := sums["px_volume_vol_written_bytes"]
	if wb == 0 {
		wb = sums["px_volume_written_bytes"]
	}
	vm.WriteBytes = wb

	return vm, nil
}

// historyWindowSec is the lookback window for chart pre-population (10 minutes).
const historyWindowSec = 600

// historyStepSec is the resolution step – matches the JS polling interval (20 s).
const historyStepSec = 20

// GetVolumeMetricsHistory fetches the last ~10 minutes of time-series data for
// the named volume from the Thanos Querier. All six chart metrics are fetched
// concurrently. The returned arrays are all index-aligned with Timestamps.
func GetVolumeMetricsHistory(metricsURL, token, volumeName string) (*VolumeMetricsHistory, error) {
	if metricsURL == "" {
		return nil, fmt.Errorf("metrics URL is not configured – set PORTWORX_METRICS_URL")
	}

	now := time.Now()
	start := now.Add(-historyWindowSec * time.Second)
	sel := `volumename="` + volumeName + `"`

	type rangeSpec struct {
		promQL  string
		scaleFn func(float64) float64 // optional transformation (e.g. s→ms)
	}

	// Six metrics that drive the six dashboard charts.
	queries := []rangeSpec{
		{promQL: `sum(px_volume_readthroughput{` + sel + `})`},
		{promQL: `sum(px_volume_writethroughput{` + sel + `})`},
		{promQL: `sum(px_volume_read_iops{` + sel + `})`},
		{promQL: `sum(px_volume_write_iops{` + sel + `})`},
		{
			promQL: `sum(px_volume_vol_read_latency_seconds{` + sel + `} or px_volume_read_latency_seconds{` + sel + `})`,
			scaleFn: func(v float64) float64 { return v * 1000 },
		},
		{
			promQL: `sum(px_volume_vol_write_latency_seconds{` + sel + `} or px_volume_write_latency_seconds{` + sel + `})`,
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

	// Collect results.
	seriesMap := make(map[int]seriesResult, len(queries))
	for sr := range ch {
		seriesMap[sr.idx] = sr
	}

	// Use the first series with data as the timestamp reference.
	var refTS []int64
	for i := 0; i < len(queries); i++ {
		if sr := seriesMap[i]; len(sr.ts) > 0 {
			refTS = sr.ts
			break
		}
	}
	if refTS == nil {
		// No data available – return an empty history rather than an error.
		return &VolumeMetricsHistory{VolumeName: volumeName}, nil
	}

	// Align all series to the reference timestamp grid.
	align := func(idx int) []float64 {
		sr := seriesMap[idx]
		if len(sr.ts) == 0 {
			return make([]float64, len(refTS))
		}
		return alignToTimestamps(refTS, buildTSMap(sr.ts, sr.vals))
	}

	return &VolumeMetricsHistory{
		VolumeName:      volumeName,
		Timestamps:      refTS,
		ReadThroughput:  align(0),
		WriteThroughput: align(1),
		ReadIOPS:        align(2),
		WriteIOPS:       align(3),
		ReadLatencyMs:   align(4),
		WriteLatencyMs:  align(5),
	}, nil
}
