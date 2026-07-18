package appserver

import (
	"context"
	"http-server/internal/core/interface/driver"
	"http-server/internal/infrastructure/config"
	"http-server/internal/infrastructure/driver/app_server/middleware"
	"log"
	"net/http"
	"time"
)

type MetricsCollector interface {
	MetricsMiddleware(next http.Handler) http.Handler
}

type AppServerAdapter struct {
	port             string
	getSiteList      driver.GetSiteList
	addSite          driver.AddSite
	updateSite       driver.UpdateSite
	removeSite       driver.RemoveSite
	getSiteStatuses  driver.GetSiteStatuses
	getSiteDetail    driver.GetSiteDetail
	metricsCollector MetricsCollector
	cfg              *config.Config
}

func NewAppServerAdapter(
	port string,
	getSiteList driver.GetSiteList,
	addSite driver.AddSite,
	updateSite driver.UpdateSite,
	removeSite driver.RemoveSite,
	getSiteStatuses driver.GetSiteStatuses,
	getSiteDetail driver.GetSiteDetail,
	metricsCollector MetricsCollector,
	cfg *config.Config,
) *AppServerAdapter {
	return &AppServerAdapter{
		port:             port,
		getSiteList:      getSiteList,
		addSite:          addSite,
		updateSite:       updateSite,
		removeSite:       removeSite,
		getSiteStatuses:  getSiteStatuses,
		getSiteDetail:    getSiteDetail,
		metricsCollector: metricsCollector,
		cfg: cfg,
	}
}

func (s *AppServerAdapter) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", s.handleHome)
	mux.HandleFunc("GET /sites/{id}", s.handleSiteDetails)

	mux.HandleFunc("GET /api/v1/site", s.handleGetSiteList)
	mux.HandleFunc("POST /api/v1/site", s.handleAddSite)
	mux.HandleFunc("PATCH /api/v1/site", s.handleUpdateSite)
	mux.HandleFunc("DELETE /api/v1/site/{id}", s.handleRemoveSite)

	wMux := middleware.ChainMiddleware(
		mux,
		s.metricsCollector.MetricsMiddleware,
		middleware.LogMiddleware,
		middleware.CorsMiddleware,
	)

	srv := http.Server{
		Addr:    s.port,
		Handler: wMux,
	}

	log.Printf("Server starting at http://localhost%v...", s.port)
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Printf("Server stopping at http://localhost%v...", s.port)
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return err
		}
		return <-errCh
	case err := <-errCh:
		return err
	}
}
