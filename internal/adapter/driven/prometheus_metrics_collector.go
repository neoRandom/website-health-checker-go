package driven

import (
	"http-server/internal/domain/ports/driven"
	"log"
	"net/http"
	"time"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusMetricsCollector struct {
	siteRepository driven.SiteRepository
}

func NewPrometheusMetricsCollector(siteRepository driven.SiteRepository) *PrometheusMetricsCollector {

	return &PrometheusMetricsCollector{
		siteRepository: siteRepository,
	}
}

func (mc *PrometheusMetricsCollector) MetricsHandler() http.Handler {
	return promhttp.Handler()
}

func (mc *PrometheusMetricsCollector) MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			mc.requestDuration(time.Since(start))
			mc.increaseRequestCounter()
		})
}

func (mc *PrometheusMetricsCollector) increaseRequestCounter() {
	log.Println("NEW REQUEST")
}

func (mc *PrometheusMetricsCollector) requestDuration(duration time.Duration) {
	log.Println("REQUEST DURATION: ", duration)
}
