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

func (ps *PrometheusMetricsExporterAdapter) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	mux.Handle("/metrics", ps.metricsCollector.MetricsHandler())

	srv := http.Server{
		Addr:    ps.addr,
		Handler: mux,
	}

	log.Printf("Prometheus metric exporter starting at http://localhost%v...", ps.addr)
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Printf("Prometheus metric exporter stopping at http://localhost%v...", ps.addr)
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return err
		}
		return <-errCh
	case err := <-errCh:
		return err
	}
}
