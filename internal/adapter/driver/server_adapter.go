package driver

import (
	"encoding/json"
	"fmt"
	"http-server/internal/adapter/driver/dto"
	"http-server/internal/adapter/driver/middleware"
	"http-server/internal/domain/ports/driver"
	"log"
	"net/http"
)

type ServerAdapter struct {
	Addr        string
	GetSiteList driver.GetSiteList
	AddSite     driver.AddSite
}

func (s *ServerAdapter) Init() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})

	mux.HandleFunc("/sites/list", s.handleSiteList)

	wMux := middleware.ChainMiddleware(
		mux,
		middleware.LogMiddleware,
		middleware.CorsMiddleware,
	)

	srv := http.Server{
		Addr:    s.Addr,
		Handler: wMux,
	}

	log.Printf("Server starting at http://%v...", s.Addr)
	srv.ListenAndServe()
}

func (s *ServerAdapter) handleSiteList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		l, err := s.GetSiteList()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}

		body := make([]*dto.SiteResponse, len(l))
		for i, site := range l {
			body[i] = &dto.SiteResponse{
				Id:  site.Id,
				Url: site.Url,
			}
		}

		json.NewEncoder(w).Encode(dto.GetSiteListResponse{
			Body: body,
		})

	case http.MethodPost:
		var req dto.AddSiteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			break
		}

		s, err := s.AddSite(req.Url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}

		json.NewEncoder(w).Encode(dto.SiteResponse{
			Id:  s.Id,
			Url: s.Url,
		})

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
