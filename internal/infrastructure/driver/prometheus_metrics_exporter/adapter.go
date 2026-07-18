package prometheusmetricsexporter

import (
	"context"
	"log"
	"net/http"
	"time"
)

type MetricsCollector interface {
	MetricsHandler() http.Handler
}

type PrometheusMetricsExporterAdapter struct {
	port             string
	metricsCollector MetricsCollector
}

func NewPrometheusMetricsExporterAdapter(
	port string,
	metricsCollector MetricsCollector,
) *PrometheusMetricsExporterAdapter {
	return &PrometheusMetricsExporterAdapter{
		port:             port,
		metricsCollector: metricsCollector,
	}
}

func (ps *PrometheusMetricsExporterAdapter) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	mux.Handle("/metrics", ps.metricsCollector.MetricsHandler())

	srv := http.Server{
		Addr:    ps.port,
		Handler: mux,
	}

	log.Printf("Prometheus metric exporter starting at http://localhost%v...", ps.port)
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Printf("Prometheus metric exporter stopping at http://localhost%v...", ps.port)
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return err
		}
		return <-errCh
	case err := <-errCh:
		return err
	}
}
