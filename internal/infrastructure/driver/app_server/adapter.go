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
	addr             string
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
	addr string,
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
		addr:             addr,
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

	mux.HandleFunc("GET /sites/list", s.handleGetSiteList)
	mux.HandleFunc("POST /sites/list", s.handleAddSite)
	mux.HandleFunc("PATCH /sites/list", s.handleUpdateSite)
	mux.HandleFunc("DELETE /sites/list/{id}", s.handleRemoveSite)

	wMux := middleware.ChainMiddleware(
		mux,
		s.metricsCollector.MetricsMiddleware,
		middleware.LogMiddleware,
		middleware.CorsMiddleware,
	)

	srv := http.Server{
		Addr:    s.addr,
		Handler: wMux,
	}

	log.Printf("Server starting at http://localhost%v...", s.addr)
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Printf("Server stopping at http://localhost%v...", s.addr)
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return err
		}
		return <-errCh
	case err := <-errCh:
		return err
	}
}
