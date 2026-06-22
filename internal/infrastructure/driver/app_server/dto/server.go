package dto

import "http-server/internal/core/model"

type AddSiteRequest struct {
	Url string `json:"url"`
}

type GetSiteListResponse struct {
	Body []*SiteJSON `json:"body"`
}

type SiteJSON struct {
	Id  model.SiteID `json:"id"`
	Url string        `json:"url"`
}
