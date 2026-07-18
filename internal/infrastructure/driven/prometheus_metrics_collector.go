package driven

import (
	"context"
	"http-server/internal/core/interface/driven"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

type PrometheusMetricsCollector struct {
	siteRepository  driven.SiteRepository
	totalRequests   prometheus.Counter
	requestDuration *prometheus.HistogramVec
}

func NewPrometheusMetricsCollector(
	ctx context.Context, siteRepository driven.SiteRepository,
) *PrometheusMetricsCollector {
	totalRequests := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request latency",
		Buckets: prometheus.DefBuckets,
	}, []string{"route", "method", "status"})

	monitoredTargets := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "monitored_targets",
		Help: "Number of actively monitored targets",
	}, func() float64 {
		c, _ := siteRepository.Count(ctx)
		return float64(c)
	})

	prometheus.MustRegister(
		totalRequests, 
		requestDuration, 
		monitoredTargets,
	)

	return &PrometheusMetricsCollector{
		siteRepository:  siteRepository,
		totalRequests:   totalRequests,
		requestDuration: requestDuration,
	}
}

func (mc *PrometheusMetricsCollector) MetricsHandler() http.Handler {
	return promhttp.Handler()
}

func (mc *PrometheusMetricsCollector) MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			wrappedW := &responseWriterWrapper{
				ResponseWriter: w,
				statusCode: http.StatusOK,
			}

			start := time.Now()
			next.ServeHTTP(wrappedW, r)
			
			mc.requestDuration.WithLabelValues(
				r.URL.Path,
				r.Method,
				strconv.Itoa(wrappedW.statusCode),
			).Observe(float64(time.Since(start)))
			mc.totalRequests.Inc()
		})
}
