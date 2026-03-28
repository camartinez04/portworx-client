/**
 * volume-metrics.js
 * Live Portworx volume metrics dashboard – Chart.js v2 compatible.
 *
 * Expects window.volumeName to be set by the host template before this script
 * is loaded.  On page load, fetches the last ~10 minutes of historical data
 * from Thanos via /portworx/client/api/metrics/{volumeName}/history to
 * pre-populate charts (no dead lines on reopen), then polls
 * /portworx/client/api/metrics/{volumeName} every REFRESH_INTERVAL ms for
 * new data points.
 */
(function () {
  'use strict';

  /* ─── Configuration ───────────────────────────────────────────────────── */
  var REFRESH_INTERVAL = 20000;   // 20 s between polls
  var MAX_POINTS       = 30;      // rolling window depth

  /* ─── Theme colours (light theme matching app palette) ───────────────── */
  var COLOR_READ      = '#71c016';       // success green  ($success)
  var COLOR_WRITE     = '#248afd';       // primary blue   ($primary)
  var COLOR_LATENCY_R = '#f0a500';       // amber/warning
  var COLOR_LATENCY_W = '#e63946';       // danger red
  var GRID_COLOR      = 'rgba(0,0,0,0.06)';
  var TICK_COLOR      = '#787878';

  /* ─── State ───────────────────────────────────────────────────────────── */
  var charts      = {};

  function zeros() { var a = []; for (var i = 0; i < MAX_POINTS; i++) a.push(0); return a; }
  function emptyLabels() { var a = []; for (var i = 0; i < MAX_POINTS; i++) a.push(''); return a; }

  var rollingData = {
    labels:       emptyLabels(),
    readRate:     zeros(),
    writeRate:    zeros(),
    readIops:     zeros(),
    writeIops:    zeros(),
    readLatency:  zeros(),
    writeLatency: zeros()
  };
  var lastSample  = null;   // { ts, readBytes, writeBytes }
  var pollTimer   = null;

  /* ─── Helpers ─────────────────────────────────────────────────────────── */
  function push(arr, val) {
    arr.push(val);
    if (arr.length > MAX_POINTS) arr.shift();
  }

  /** Format bytes/s → human-readable throughput string. */
  function fmtBytes(b) {
    if (b < 0) b = 0;
    if (b < 1024)        return b.toFixed(0)      + ' B/s';
    if (b < 1048576)     return (b / 1024).toFixed(1)    + ' KB/s';
    if (b < 1073741824)  return (b / 1048576).toFixed(2) + ' MB/s';
    return (b / 1073741824).toFixed(2) + ' GB/s';
  }

  /** Format IOPS value. */
  function fmtIops(v) {
    if (v >= 1000) return (v / 1000).toFixed(1) + 'K';
    return Math.round(v).toString();
  }

  /** Format latency value. */
  function fmtLatency(ms) {
    if (ms < 1) return (ms * 1000).toFixed(0) + ' µs';
    return ms.toFixed(2) + ' ms';
  }

  function setText(id, val) {
    var el = document.getElementById(id);
    if (el) el.textContent = val;
  }

  function setRefreshState(state) {
    var dot = document.getElementById('metrics-refresh-dot');
    if (!dot) return;
    dot.className = 'metrics-refresh-dot metrics-dot-' + state;
  }

  /** Format a Unix-ms timestamp as a locale time string for chart labels. */
  function tsToLabel(unixMs) {
    return new Date(unixMs).toLocaleTimeString();
  }

  /* ─── Chart factory ───────────────────────────────────────────────────── */
  /**
   * @param {string} canvasId
   * @param {string} label
   * @param {string} color
   * @param {function} tickFmt  – value formatter for y-axis ticks
   * @param {function} tooltipFmt – value formatter for tooltips
   */
  function makeChart(canvasId, label, color, tickFmt, tooltipFmt, dataArr) {
    var el = document.getElementById(canvasId);
    if (!el) return null;

    return new Chart(el.getContext('2d'), {
      type: 'line',
      data: {
        labels: rollingData.labels.slice(),
        datasets: [{
          label:                label,
          data:                 dataArr.slice(),
          borderColor:          color,
          backgroundColor:      color + '22',
          borderWidth:          2,
          pointRadius:          2,
          pointHoverRadius:     5,
          pointBackgroundColor: color,
          fill:                 true,
          lineTension:          0.35,
          spanGaps:             true
        }]
      },
      options: {
        responsive:          true,
        maintainAspectRatio: false,
        animation:           { duration: 250 },
        legend:              { display: false },
        tooltips: {
          mode:        'index',
          intersect:   false,
          displayColors: false,
          callbacks: {
            label: function (item) {
              return label + ': ' + tooltipFmt(item.yLabel);
            }
          }
        },
        scales: {
          xAxes: [{
            gridLines: { color: GRID_COLOR, zeroLineColor: GRID_COLOR },
            ticks:     { fontColor: TICK_COLOR, maxTicksLimit: 6, maxRotation: 0 }
          }],
          yAxes: [{
            gridLines: { color: GRID_COLOR, zeroLineColor: GRID_COLOR },
            ticks: {
              fontColor:  TICK_COLOR,
              beginAtZero: true,
              callback:   tickFmt
            }
          }]
        }
      }
    });
  }

  /* ─── Chart initialisation ────────────────────────────────────────────── */
  function initCharts() {
    charts.readRate    = makeChart('mx-read-bytes',    'Read Throughput',  COLOR_READ,      fmtBytes,   fmtBytes,                                  rollingData.readRate);
    charts.writeRate   = makeChart('mx-write-bytes',   'Write Throughput', COLOR_WRITE,     fmtBytes,   fmtBytes,                                  rollingData.writeRate);
    charts.readIops    = makeChart('mx-read-iops',     'Read IOPS',        COLOR_READ,      fmtIops,    function(v){ return fmtIops(v) + ' IOPS'; }, rollingData.readIops);
    charts.writeIops   = makeChart('mx-write-iops',    'Write IOPS',       COLOR_WRITE,     fmtIops,    function(v){ return fmtIops(v) + ' IOPS'; }, rollingData.writeIops);
    charts.readLatency = makeChart('mx-read-latency',  'Read Latency',     COLOR_LATENCY_R, fmtLatency, fmtLatency,                                rollingData.readLatency);
    charts.writeLatency= makeChart('mx-write-latency', 'Write Latency',    COLOR_LATENCY_W, fmtLatency, fmtLatency,                                rollingData.writeLatency);
  }

  /* ─── Chart update ────────────────────────────────────────────────────── */
  function refreshChart(chart, labels, data) {
    if (!chart) return;
    chart.data.labels              = labels.slice();
    chart.data.datasets[0].data   = data.slice();
    chart.update(0);   // 0 = skip animation for live updates
  }

  function refreshAllCharts() {
    refreshChart(charts.readRate,    rollingData.labels, rollingData.readRate);
    refreshChart(charts.writeRate,   rollingData.labels, rollingData.writeRate);
    refreshChart(charts.readIops,    rollingData.labels, rollingData.readIops);
    refreshChart(charts.writeIops,   rollingData.labels, rollingData.writeIops);
    refreshChart(charts.readLatency, rollingData.labels, rollingData.readLatency);
    refreshChart(charts.writeLatency,rollingData.labels, rollingData.writeLatency);
  }

  /* ─── History pre-population ──────────────────────────────────────────── */
  /**
   * Fetches the last ~10 minutes of Thanos range-query data and pre-populates
   * the rolling window so charts show real history instead of flat zeros.
   * Resolves (with or without data) so that polling always starts afterward.
   */
  function loadHistory(volName) {
    return fetch('/portworx/client/api/metrics/' + encodeURIComponent(volName) + '/history')
      .then(function (r) {
        if (!r.ok) throw new Error('HTTP ' + r.status);
        return r.json();
      })
      .then(function (h) {
        if (!h.timestamps || h.timestamps.length === 0) return;

        // Take up to MAX_POINTS most-recent points from the history window.
        var pts = h.timestamps.length;
        var start = Math.max(0, pts - MAX_POINTS);

        rollingData.labels       = [];
        rollingData.readRate     = [];
        rollingData.writeRate    = [];
        rollingData.readIops     = [];
        rollingData.writeIops    = [];
        rollingData.readLatency  = [];
        rollingData.writeLatency = [];

        for (var i = start; i < pts; i++) {
          rollingData.labels.push(tsToLabel(h.timestamps[i]));
          rollingData.readRate.push(h.read_throughput_bytes_s[i]  || 0);
          rollingData.writeRate.push(h.write_throughput_bytes_s[i] || 0);
          rollingData.readIops.push(h.read_iops[i]                || 0);
          rollingData.writeIops.push(h.write_iops[i]              || 0);
          rollingData.readLatency.push(h.read_latency_ms[i]       || 0);
          rollingData.writeLatency.push(h.write_latency_ms[i]     || 0);
        }

        // Pad to MAX_POINTS with leading empty slots if fewer points returned.
        while (rollingData.labels.length < MAX_POINTS) {
          rollingData.labels.unshift('');
          rollingData.readRate.unshift(0);
          rollingData.writeRate.unshift(0);
          rollingData.readIops.unshift(0);
          rollingData.writeIops.unshift(0);
          rollingData.readLatency.unshift(0);
          rollingData.writeLatency.unshift(0);
        }

        refreshAllCharts();
      })
      .catch(function (err) {
        // Non-fatal – history unavailable; charts start from scratch.
        console.warn('[metrics] history load failed:', err);
      });
  }

  /* ─── Metrics fetch & update ──────────────────────────────────────────── */
  function fetchMetrics() {
    var volName = window.volumeName;
    if (!volName) return;

    setRefreshState('loading');

    fetch('/portworx/client/api/metrics/' + encodeURIComponent(volName))
      .then(function (r) {
        if (!r.ok) throw new Error('HTTP ' + r.status);
        return r.json();
      })
      .then(function (data) {
        if (data.error) throw new Error(data.message || 'broker error');

        var now    = new Date().toLocaleTimeString();
        var nowMs  = Date.now();

        /* ── Throughput (prefer pre-computed gauges; fall back to delta/s) ── */
        var readRate = 0, writeRate = 0;

        if (data.read_throughput_bytes_s > 0 || data.write_throughput_bytes_s > 0) {
          readRate  = data.read_throughput_bytes_s  || 0;
          writeRate = data.write_throughput_bytes_s || 0;
        } else if (lastSample) {
          var dtSec = (nowMs - lastSample.ts) / 1000;
          if (dtSec > 0) {
            readRate  = Math.max(0, (data.read_bytes  - lastSample.readBytes)  / dtSec);
            writeRate = Math.max(0, (data.write_bytes - lastSample.writeBytes) / dtSec);
          }
        }
        lastSample = { ts: nowMs, readBytes: data.read_bytes, writeBytes: data.write_bytes };

        var readIops     = data.read_iops       || 0;
        var writeIops    = data.write_iops      || 0;
        var readLatMs    = data.read_latency_ms || 0;
        var writeLatMs   = data.write_latency_ms|| 0;

        /* ── Push to rolling window ── */
        push(rollingData.labels,       now);
        push(rollingData.readRate,     readRate);
        push(rollingData.writeRate,    writeRate);
        push(rollingData.readIops,     readIops);
        push(rollingData.writeIops,    writeIops);
        push(rollingData.readLatency,  readLatMs);
        push(rollingData.writeLatency, writeLatMs);

        refreshAllCharts();

        /* ── Update stat cards ── */
        setText('stat-read-throughput',  fmtBytes(readRate));
        setText('stat-write-throughput', fmtBytes(writeRate));
        setText('stat-read-iops',        fmtIops(readIops) + ' IOPS');
        setText('stat-write-iops',       fmtIops(writeIops) + ' IOPS');
        setText('stat-read-latency',     fmtLatency(readLatMs));
        setText('stat-write-latency',    fmtLatency(writeLatMs));
        setText('stat-io-depth',         Math.round(data.io_depth || 0).toString());

        setText('metrics-last-updated', 'Updated ' + now);
        setRefreshState('ok');
      })
      .catch(function (err) {
        console.error('[metrics]', err);
        setRefreshState('error');
        setText('metrics-last-updated', 'Error – retrying…');
      });
  }

  /* ─── Bootstrap ───────────────────────────────────────────────────────── */
  document.addEventListener('DOMContentLoaded', function () {
    if (!window.volumeName) return;

    initCharts();

    // Pre-populate charts from Thanos history, then start live polling.
    loadHistory(window.volumeName).then(function () {
      fetchMetrics();
      pollTimer = setInterval(fetchMetrics, REFRESH_INTERVAL);
    });
  });

})();
