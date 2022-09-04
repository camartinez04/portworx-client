(function ($) {
  'use strict';
  $(function () {

    if ($("#node-chart").length) {
      var areaData = {
        labels: ["Used", "Available"],
        datasets: [{
            data: [window.usagePercent, window.availablePercent],
            backgroundColor: [
              "#71c016", "#8caaff", ,
            ],
            borderColor: "rgba(0,0,0,0)"
          }
        ]
      };
      var areaOptions = {
        responsive: true,
        maintainAspectRatio: true,
        segmentShowStroke: false,
        cutoutPercentage: 78,
        elements: {
          arc: {
              borderWidth: 4
          }
        },      
        legend: {
          display: false
        },
        tooltips: {
          enabled: true
        },
        legendCallback: function (chart) {
          var text = [];
          text.push('<div class="report-chart">');
          text.push('<div class="d-flex justify-content-between mx-4 mx-xl-5 mt-3"><div class="d-flex align-items-center"><div class="me-3" style="width:20px; height:20px; border-radius: 50%; background-color: ' + chart.data.datasets[0].backgroundColor[0] + '"></div><p class="mb-0">Used space</p></div>');
          text.push('<p class="mb-0">' + Math.round((window.usage / 1024) * 100) / 100 + ' GB</p>');
          text.push('</div>');
          text.push('<div class="d-flex justify-content-between mx-4 mx-xl-5 mt-3"><div class="d-flex align-items-center"><div class="me-3" style="width:20px; height:20px; border-radius: 50%; background-color: ' + chart.data.datasets[0].backgroundColor[1] + '"></div><p class="mb-0">Available space</p></div>');
          text.push('<p class="mb-0">' + Math.round((window.available / 1024) * 100) / 100 + ' GB </p>');
          text.push('</div>');
          text.push('</div>');
          return text.join("");
        },
      }
      var northAmericaChartPlugins = {
        beforeDraw: function(chart) {
          var width = chart.chart.width,
              height = chart.chart.height,
              ctx = chart.chart.ctx;
      
          ctx.restore();
          var fontSize = 3.125;
          ctx.font = "600 " + fontSize + "em sans-serif";
          ctx.textBaseline = "middle";
          ctx.fillStyle = "#000";
      
          var text = window.total / 1024,
              textX = Math.round((width - ctx.measureText(text).width) / 2),
              textY = height / 2;
      
          ctx.fillText(text, textX, textY);
          ctx.save();
        }
      }
      var northAmericaChartCanvas = $("#node-chart").get(0).getContext("2d");
      var northAmericaChart = new Chart(northAmericaChartCanvas, {
        type: 'doughnut',
        data: areaData,
        options: areaOptions,
        plugins: northAmericaChartPlugins
      });
      document.getElementById('node-legend').innerHTML = northAmericaChart.generateLegend();
      console.log(window.usagePercent,window.availablePercent, window.usage, window.available, window.total);

    }
    
    if ($("#volume-utilization-chart").length) {


      var areaData = {
        labels: ["Used", "Available"],

        datasets: [{
          data: [window.usagePercent, window.availablePercent],
          backgroundColor: [
            "#71c016", "#8caaff",
          ],
          borderColor: "rgba(0,0,0,0)"
        }
        ]
      };
      var areaOptions = {
        responsive: true,
        maintainAspectRatio: true,
        segmentShowStroke: false,
        cutoutPercentage: 78,
        elements: {
          arc: {
            borderWidth: 4
          }
        },
        legend: {
          display: false
        },
        tooltips: {
          enabled: true
        },
        legendCallback: function (chart) {
          var text = [];
          text.push('<div class="report-chart">');
          text.push('<div class="d-flex justify-content-between mx-4 mx-xl-5 mt-3"><div class="d-flex align-items-center"><div class="me-3" style="width:20px; height:20px; border-radius: 50%; background-color: ' + chart.data.datasets[0].backgroundColor[0] + '"></div><p class="mb-0">Used space</p></div>');
          text.push('<p class="mb-0">' + Math.round((window.usage / 1024) * 100) / 100 + ' GB</p>');
          text.push('</div>');
          text.push('<div class="d-flex justify-content-between mx-4 mx-xl-5 mt-3"><div class="d-flex align-items-center"><div class="me-3" style="width:20px; height:20px; border-radius: 50%; background-color: ' + chart.data.datasets[0].backgroundColor[1] + '"></div><p class="mb-0">Available space</p></div>');
          text.push('<p class="mb-0">' + Math.round((window.available / 1024) * 100) / 100 + ' GB </p>');
          text.push('</div>');
          text.push('</div>');
          return text.join("");
        },
      }
      var volumeUtilizationChartPlugins = {
        beforeDraw: function (chart) {
          var width = chart.chart.width,
            height = chart.chart.height,
            ctx = chart.chart.ctx;

          ctx.restore();
          var fontSize = 3.125;
          ctx.font = "600 " + fontSize + "em sans-serif";
          ctx.textBaseline = "middle";
          ctx.fillStyle = "#000";

          var text = window.total / 1024,
            textX = Math.round((width - ctx.measureText(text).width) / 2),
            textY = height / 2;

          ctx.fillText(text, textX, textY);
          ctx.save();
        }
      }
      var volumeUtilizationChartCanvas = $("#volume-utilization-chart").get(0).getContext("2d");
      var volumeUtilizationChart = new Chart(volumeUtilizationChartCanvas, {
        type: 'doughnut',
        data: areaData,
        options: areaOptions,
        plugins: volumeUtilizationChartPlugins
      });
      document.getElementById('volume-utilization-legend').innerHTML = volumeUtilizationChart.generateLegend();
      console.log(window.usagePercent,window.availablePercent, window.usage, window.available, window.total);

    }

  }
  );
})(jQuery);