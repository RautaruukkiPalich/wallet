package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
	"time"
)

var (
	TotalPostRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_post_requests_total",
		Help: "Number of post requests",
	})

	TotalGetRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_get_requests_total",
		Help: "Number of get requests",
	})
)

var histogramCodeVec = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "hist",
		Name:      "response_time",
		Help:      "Time handled by status code",
	},
	[]string{
		"code",
	},
)

func ObserveHistogramCodeResponseVec(code int, dur time.Duration) {
	histogramCodeVec.WithLabelValues(
		strconv.Itoa(code),
	).Observe(dur.Seconds())
}

var histogramMethodCounter = promauto.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace: "hist",
		Name:      "method",
		Help:      "Request method",
	},
	[]string{
		"method",
		"path",
	},
)

func ObserveHistogramMethodResponseVec(method string, path string, dur time.Duration) {
	histogramMethodCounter.WithLabelValues(
		method,
		path,
	).Observe(dur.Seconds())
}
