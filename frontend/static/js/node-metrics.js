/**
 * node-metrics.js
 * Live Portworx node metrics dashboard – Chart.js v2 compatible.
 *
 * Expects window.nodeID to be set by the host template before this script
 * is loaded.  On page load, fetches the last ~10 minutes of historical data
 * from Thanos via /portworx/client/api/node-metrics/{nodeID}/history to
 * pre-populate charts (no dead lines on reopen), then polls
 * /portworx/client/api/node-metrics/{nodeID} every REFRESH_INTERVAL ms for
 * new data points and a live pool snapshot.
 */
(function () {
  'use strict';

  /* ─── Configuration ───────────────────────────────────────────────────── */
  var REFRESH_INTERVAL = 20000;   // 20 s between polls
  var MAX_POINTS       = 30;      // rolling window depth

  /* ─── Theme colours (light theme matching app palette) ───────────────── */
  var COLOR_READ      = '#71c016';  // success green
  var COLOR_WRITE     = '#248afd';  // primary blue
  var COLOR_LATENCY_R = '#f0a500';  // amber
  var COLOR_LATENCY_W = '#e63946';  // red
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
  var pollTimer = null;

  /* ─── Helpers ─────────────────────────────────────────────────────────── */
  function push(arr, val) {
    arr.push(val);
    if (arr.length > MAX_POINTS) arr.shift();
  }

  function fmtBytes(b) {
    if (b < 0) b = 0;
    if (b < 1024)        return b.toFixed(0)       + ' B/s';
    if (b < 1048576)     return (b / 1024).toFixed(1)    + ' KB/s';
    if (b < 1073741824)  return (b / 1048576).toFixed(2) + ' MB/s';
    return (b / 1073741824).toFixed(2) + ' GB/s';
  }

  function fmtBytesSize(b) {
    if (b < 0) b = 0;
    if (b < 1073741824)  return (b / 1048576).toFixed(1)    + ' MB';
    if (b < 1099511627776) return (b / 1073741824).toFixed(2) + ' GB';
    return (b / 1099511627776).toFixed(2) + ' TB';
  }

  function fmtIops(v) {
    if (v >= 1000) return (v / 1000).toFixed(1) + 'K';
    return Math.round(v).toString();
  }

  function fmtLatency(ms) {
    if (ms < 1) return (ms * 1000).toFixed(0) + ' µs';
    return ms.toFixed(2) + ' ms';
  }

  function fmtPercent(v) {
    return (v || 0).toFixed(1) + '%';
  }

  function setText(id, val) {
    var el = document.getElementById(id);
    if (el) el.textContent = val;
  }

  function setRefreshState(state) {
    var dot = document.getElementById('node-metrics-refresh-dot');
    if (!dot) return;
    dot.className = 'metrics-refresh-dot metrics-dot-' + state;
  }

  /** Format a Unix-ms timestamp as a locale time string for chart labels. */
  function tsToLabel(unixMs) {
    return new Date(unixMs).toLocaleTimeString();
  }

  /* ─── Chart factory ───────────────────────────────────────────────────── */
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
          backgroundColor:      color + '20',
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
          mode:          'index',
          intersect:     false,
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
              fontColor:   TICK_COLOR,
              beginAtZero: true,
              callback:    tickFmt
            }
          }]
        }
      }
    });
  }

  /* ─── Chart initialisation ────────────────────────────────────────────── */
  function initCharts() {
    charts.readRate    = makeChart('nd-read-throughput',  'Read Throughput',  COLOR_READ,      fmtBytes,   fmtBytes,                                  rollingData.readRate);
    charts.writeRate   = makeChart('nd-write-throughput', 'Write Throughput', COLOR_WRITE,     fmtBytes,   fmtBytes,                                  rollingData.writeRate);
    charts.readIops    = makeChart('nd-read-iops',        'Read IOPS',        COLOR_READ,      fmtIops,    function(v){ return fmtIops(v) + ' IOPS'; }, rollingData.readIops);
    charts.writeIops   = makeChart('nd-write-iops',       'Write IOPS',       COLOR_WRITE,     fmtIops,    function(v){ return fmtIops(v) + ' IOPS'; }, rollingData.writeIops);
    charts.readLatency = makeChart('nd-read-latency',     'Read Latency',     COLOR_LATENCY_R, fmtLatency, fmtLatency,                                rollingData.readLatency);
    charts.writeLatency= makeChart('nd-write-latency',    'Write Latency',    COLOR_LATENCY_W, fmtLatency, fmtLatency,                                rollingData.writeLatency);
  }

  /* ─── Chart update ────────────────────────────────────────────────────── */
  function refreshChart(chart, labels, data) {
    if (!chart) return;
    chart.data.labels            = labels.slice();
    chart.data.datasets[0].data = data.slice();
    chart.update(0);
  }

  function refreshAllCharts() {
    refreshChart(charts.readRate,    rollingData.labels, rollingData.readRate);
    refreshChart(charts.writeRate,   rollingData.labels, rollingData.writeRate);
    refreshChart(charts.readIops,    rollingData.labels, rollingData.readIops);
    refreshChart(charts.writeIops,   rollingData.labels, rollingData.writeIops);
    refreshChart(charts.readLatency, rollingData.labels, rollingData.readLatency);
    refreshChart(charts.writeLatency,rollingData.labels, rollingData.writeLatency);
  }

  /* ─── Pool cards renderer ─────────────────────────────────────────────── */
  function renderPoolCards(pools) {
    var container = document.getElementById('nd-pool-cards');
    if (!container) return;

    if (!pools || Object.keys(pools).length === 0) {
      container.innerHTML = '<p class="text-muted" style="font-size:12px;">No pool data available.</p>';
      return;
    }

    var html = '<div class="row">';
    Object.keys(pools).forEach(function (uuid) {
      var p        = pools[uuid];
      var usedPct  = p.total_bytes > 0 ? ((p.used_bytes / p.total_bytes) * 100).toFixed(1) : 0;
      var shortId  = uuid.substring(0, 8) + '…';
      var barColor = usedPct > 85 ? '#e63946' : usedPct > 65 ? '#f0a500' : '#71c016';

      html += '<div class="col-md-6 col-xl-4 mb-3">' +
        '<div class="nd-pool-card">' +
          '<div class="nd-pool-header">' +
            '<span class="nd-pool-id" title="' + uuid + '">' + shortId + '</span>' +
            '<span class="nd-pool-pct" style="color:' + barColor + ';">' + usedPct + '%</span>' +
          '</div>' +
          '<div class="nd-pool-bar-track">' +
            '<div class="nd-pool-bar-fill" style="width:' + usedPct + '%;background:' + barColor + ';"></div>' +
          '</div>' +
          '<div class="nd-pool-stats">' +
            '<span>' + fmtBytesSize(p.used_bytes) + ' used</span>' +
            '<span>' + fmtBytesSize(p.total_bytes) + ' total</span>' +
          '</div>' +
          '<div class="nd-pool-io">' +
            '<span style="color:#71c016;">↓ ' + fmtBytes(p.read_throughput_bytes_s || 0) + '</span>' +
            '<span style="color:#248afd;">↑ ' + fmtBytes(p.write_throughput_bytes_s || 0) + '</span>' +
            '<span>' + fmtIops(p.read_iops || 0) + ' R IOPS</span>' +
            '<span>' + fmtIops(p.write_iops || 0) + ' W IOPS</span>' +
          '</div>' +
        '</div>' +
      '</div>';
    });
    html += '</div>';
    container.innerHTML = html;
  }

  /* ─── History pre-population ──────────────────────────────────────────── */
  /**
   * Fetches the last ~10 minutes of Thanos range-query data and pre-populates
   * the rolling window so charts show real history on first open.
   */
  function loadHistory(nid) {
    return fetch('/portworx/client/api/node-metrics/' + encodeURIComponent(nid) + '/history')
      .then(function (r) {
        if (!r.ok) throw new Error('HTTP ' + r.status);
        return r.json();
      })
      .then(function (h) {
        if (!h.timestamps || h.timestamps.length === 0) return;

        var pts   = h.timestamps.length;
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
        console.warn('[node-metrics] history load failed:', err);
      });
  }

  /* ─── Metrics fetch & update ──────────────────────────────────────────── */
  function fetchMetrics() {
    var nid = window.nodeID;
    if (!nid) return;

    setRefreshState('loading');

    fetch('/portworx/client/api/node-metrics/' + encodeURIComponent(nid))
      .then(function (r) {
        if (!r.ok) throw new Error('HTTP ' + r.status);
        return r.json();
      })
      .then(function (data) {
        if (data.error) throw new Error(data.message || 'broker error');

        var now = new Date().toLocaleTimeString();

        var readRate  = data.read_throughput_bytes_s  || 0;
        var writeRate = data.write_throughput_bytes_s || 0;
        var readIops  = data.read_iops                || 0;
        var writeIops = data.write_iops               || 0;
        var readLat   = data.read_latency_ms          || 0;
        var writeLat  = data.write_latency_ms         || 0;

        /* ── Push to rolling window ── */
        push(rollingData.labels,       now);
        push(rollingData.readRate,     readRate);
        push(rollingData.writeRate,    writeRate);
        push(rollingData.readIops,     readIops);
        push(rollingData.writeIops,    writeIops);
        push(rollingData.readLatency,  readLat);
        push(rollingData.writeLatency, writeLat);

        refreshAllCharts();

        /* ── Update stat cards ── */
        setText('nd-stat-read-tp',    fmtBytes(readRate));
        setText('nd-stat-write-tp',   fmtBytes(writeRate));
        setText('nd-stat-read-iops',  fmtIops(readIops) + ' IOPS');
        setText('nd-stat-write-iops', fmtIops(writeIops) + ' IOPS');
        setText('nd-stat-read-lat',   fmtLatency(readLat));
        setText('nd-stat-write-lat',  fmtLatency(writeLat));
        setText('nd-stat-cpu',        fmtPercent(data.cpu_percent));
        setText('nd-stat-mem',        fmtPercent(
          data.total_memory_bytes > 0
            ? (data.used_memory_bytes / data.total_memory_bytes) * 100
            : 0
        ));
        setText('nd-stat-vols',       Math.round(data.num_volumes || 0).toString());

        /* ── Render pool cards ── */
        renderPoolCards(data.storage_pools);

        setText('node-metrics-last-updated', 'Updated ' + now);
        setRefreshState('ok');
      })
      .catch(function (err) {
        console.error('[node-metrics]', err);
        setRefreshState('error');
        setText('node-metrics-last-updated', 'Error – retrying…');
      });
  }

  /* ─── Bootstrap ───────────────────────────────────────────────────────── */
  document.addEventListener('DOMContentLoaded', function () {
    if (!window.nodeID) return;

    initCharts();

    // Pre-populate charts from Thanos history, then start live polling.
    loadHistory(window.nodeID).then(function () {
      fetchMetrics();
      pollTimer = setInterval(fetchMetrics, REFRESH_INTERVAL);
    });
  });

})();
