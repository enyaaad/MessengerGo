package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "endpoint", "status"},
	)

	HTTPRequestsDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_requests_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

func init() {
	prometheus.MustRegister(HTTPRequestsTotal, HTTPRequestsDuration)
}
