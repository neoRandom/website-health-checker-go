package driver

import (
	"encoding/json"
	"http-server/internal/adapter/driver/dto"
	"http-server/internal/adapter/driver/middleware"
	"http-server/internal/domain/ports/driver"
	"log"
	"net/http"
)

type ServerAdapter struct {
	addr        string
	getSiteList driver.GetSiteList
	addSite     driver.AddSite
	updateSite  driver.UpdateSite
	removeSite  driver.RemoveSite
}

func NewServerAdapter(
	addr string,
	getSiteList driver.GetSiteList,
	addSite driver.AddSite,
	updateSite driver.UpdateSite,
	removeSite driver.RemoveSite,
) *ServerAdapter {
	return &ServerAdapter{
		addr:        addr,
		getSiteList: getSiteList,
		addSite:     addSite,
		updateSite:  updateSite,
		removeSite:  removeSite,
	}
}

func (s *ServerAdapter) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	mux.HandleFunc("/sites/list", s.handleSiteList)

	wMux := middleware.ChainMiddleware(
		mux,
		middleware.LogMiddleware,
		middleware.CorsMiddleware,
	)

	srv := http.Server{
		Addr:    s.addr,
		Handler: wMux,
	}

	log.Printf("Server starting at http://%v...", s.addr)
	srv.ListenAndServe()
}

func (s *ServerAdapter) handleSiteList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		l, err := s.getSiteList()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}

		body := make([]*dto.SiteJSON, len(l))
		for i, site := range l {
			body[i] = &dto.SiteJSON{
				Id:  site.Id,
				Url: site.Url,
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(dto.GetSiteListResponse{
			Body: body,
		})

	case http.MethodPost:
		var req dto.AddSiteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			break
		}

		s, err := s.addSite(req.Url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(dto.SiteJSON{
			Id:  s.Id,
			Url: s.Url,
		})

	case http.MethodPut:
		var req dto.SiteJSON
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			break
		}

		err := s.updateSite(req.Id, req.Url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))

	case http.MethodDelete:
		break

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
