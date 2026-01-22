package httpapi

import (
	"net/http"
	"strconv"
	"time"

	"hirifyGOTest/pkg/metrics"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(sw, r)

		status := strconv.Itoa(sw.status)
		method := r.Method
		endpoint := r.URL.Path
		metrics.HTTPRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
		metrics.HTTPRequestsDuration.WithLabelValues(method, endpoint).Observe(time.Since(start).Seconds())
	})
}
