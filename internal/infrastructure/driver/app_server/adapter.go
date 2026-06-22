package appserver

import (
	"encoding/json"
	"http-server/internal/core/interface/driver"
	models "http-server/internal/core/model"
	"http-server/internal/infrastructure/driver/app_server/dto"
	"http-server/internal/infrastructure/driver/app_server/middleware"
	"log"
	"net/http"
	"strconv"
)

type MetricsCollector interface {
	MetricsHandler() http.Handler
	MetricsMiddleware(next http.Handler) http.Handler
}

type AppServerAdapter struct {
	addr             string
	getSiteList      driver.GetSiteList
	addSite          driver.AddSite
	updateSite       driver.UpdateSite
	removeSite       driver.RemoveSite
	metricsCollector MetricsCollector
}

func NewAppServerAdapter(
	addr string,
	getSiteList driver.GetSiteList,
	addSite driver.AddSite,
	updateSite driver.UpdateSite,
	removeSite driver.RemoveSite,
	metricsCollector MetricsCollector,
) *AppServerAdapter {
	return &AppServerAdapter{
		addr:             addr,
		getSiteList:      getSiteList,
		addSite:          addSite,
		updateSite:       updateSite,
		removeSite:       removeSite,
		metricsCollector: metricsCollector,
	}
}

func (s *AppServerAdapter) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	mux.HandleFunc("GET /sites/list", s.handleGetSiteList)
	mux.HandleFunc("POST /sites/list", s.handleAddSite)
	mux.HandleFunc("PUT /sites/list", s.handleUpdateSite)
	mux.HandleFunc("DELETE /sites/list/{id}", s.handleRemoveSite)

	mux.Handle("/metrics", s.metricsCollector.MetricsHandler())

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
	return srv.ListenAndServe()
}

func (s *AppServerAdapter) handleGetSiteList(w http.ResponseWriter, r *http.Request) {
	sList, err := s.getSiteList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body := make([]*dto.SiteJSON, len(sList))
	for i, site := range sList {
		body[i] = &dto.SiteJSON{
			Id:  site.Id,
			Url: site.Url,
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.GetSiteListResponse{
		Body: body,
	})
}

func (s *AppServerAdapter) handleAddSite(w http.ResponseWriter, r *http.Request) {
	var req dto.AddSiteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := s.addSite(req.Url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.SiteJSON{
		Id:  created.Id,
		Url: created.Url,
	})
}

func (s *AppServerAdapter) handleUpdateSite(w http.ResponseWriter, r *http.Request) {
	var req dto.SiteJSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := s.updateSite(req.Id, req.Url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "success"}`))
}

func (s *AppServerAdapter) handleRemoveSite(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.removeSite(models.SiteID(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}
