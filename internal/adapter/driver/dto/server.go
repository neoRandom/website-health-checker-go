package dto

import "http-server/internal/domain/models"

type AddSiteRequest struct {
	Url string `json:"url"`
}

type GetSiteListResponse struct {
	Body []*SiteJSON `json:"body"`
}

type SiteJSON struct {
	Id  models.SiteID `json:"id"`
	Url string        `json:"url"`
}
