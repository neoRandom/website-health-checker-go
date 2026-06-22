package prometheusmetricsexporter

import (
	"log"
	"net/http"
)

type MetricsCollector interface {
	MetricsHandler() http.Handler
}

type PrometheusMetricsExporterAdapter struct {
	addr             string
	metricsCollector MetricsCollector
}

func NewPrometheusMetricsExporterAdapter(
	addr string, 
	metricsCollector MetricsCollector,
) *PrometheusMetricsExporterAdapter {
	return &PrometheusMetricsExporterAdapter{
		addr:             addr,
		metricsCollector: metricsCollector,
	}
}

func (ps *PrometheusMetricsExporterAdapter) Start() error {
	mux := http.NewServeMux()

	mux.Handle("/metrics", ps.metricsCollector.MetricsHandler())

	srv := http.Server{
		Addr:    ps.addr,
		Handler: mux,
	}

	log.Printf("Prometheus metric exporter starting at http://localhost%v...", ps.addr)
	return srv.ListenAndServe()
}
