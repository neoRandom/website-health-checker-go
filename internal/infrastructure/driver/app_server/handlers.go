package appserver

import (
	"encoding/json"
	"http-server/internal/core/model"
	"http-server/internal/infrastructure/driver/app_server/dto"
	"http-server/internal/infrastructure/web"
	"net/http"
	"strconv"
	"time"
)

// =============================================================================

func (s *AppServerAdapter) handleHome(
	w http.ResponseWriter, r *http.Request,
) {
	statuses, err := s.getSiteStatuses()
	if err != nil {
		http.Error(w, "unable to load sites", http.StatusInternalServerError)
		return
	}

	data := web.BuildDashboardData(s.cfg, statuses, time.Now())

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	web.HomePage(data).Render(r.Context(), w)
}

func (s *AppServerAdapter) handleSiteDetails(
	w http.ResponseWriter, r *http.Request,
) {
	idString := r.PathValue("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	status, history, err := s.getSiteDetail(model.SiteID(id), 100)
	if err != nil {
		http.Error(w, "unable to load site", http.StatusInternalServerError)
		return
	}
	if status == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`<p class="text-sm text-neutral-600">Site not found.</p>`))
		return
	}

	detail := web.BuildSiteDetail(status.Site, history, time.Now())

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	web.SiteDetailPanel(detail).Render(r.Context(), w)
}

// =============================================================================

func (s *AppServerAdapter) handleGetSiteList(
	w http.ResponseWriter, r *http.Request,
) {
	sList, err := s.getSiteList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body := make([]dto.SiteJSON, len(sList))
	for i, site := range sList {
		body[i] = dto.SiteJSON{
			Id:          site.Id,
			Url:         site.Url,
			Description: site.Description,
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.GetSiteListResponse{
		Body: body,
	})
}

func (s *AppServerAdapter) handleAddSite(
	w http.ResponseWriter, r *http.Request,
) {
	var req dto.AddSiteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	site := &model.Site{
		Url:                req.Url,
		ExpectedStatusCode: req.ExpectedStatusCode,
		Description:        req.Description,
	}

	id, err := s.addSite(site)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.SiteJSON{
		Id:                 id,
		Url:                site.Url,
		ExpectedStatusCode: site.ExpectedStatusCode,
		Description:        site.Description,
	})
}

func (s *AppServerAdapter) handleUpdateSite(
	w http.ResponseWriter, r *http.Request,
) {
	var req dto.SiteJSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	site := &model.Site{
		Id:                 req.Id,
		Url:                req.Url,
		ExpectedStatusCode: req.ExpectedStatusCode,
		Description:        req.Description,
	}

	err := s.updateSite(site)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "success"}`))
}

func (s *AppServerAdapter) handleRemoveSite(
	w http.ResponseWriter, r *http.Request,
) {
	idString := r.PathValue("id")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.removeSite(model.SiteID(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}