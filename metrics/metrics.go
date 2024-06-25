package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_total",
			Help: "Total number of requests received",
		},
		[]string{"method", "endpoint"},
	)
	RequestSuccessTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_success_total",
			Help: "Total number of successful requests (2xx, 3xx)",
		},
		[]string{"method", "endpoint"},
	)
	RequestFailureTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_failure_total",
			Help: "Total number of failed requests (4xx, 5xx)",
		},
		[]string{"method", "endpoint"},
	)
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "Duration of requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
	RequestsPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "requests_per_second",
			Help: "Number of requests per second",
		},
		[]string{"method", "endpoint"},
	)
)

func Init() {
	prometheus.MustRegister(RequestTotal)
	prometheus.MustRegister(RequestSuccessTotal)
	prometheus.MustRegister(RequestFailureTotal)
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(RequestsPerSecond)
}

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
