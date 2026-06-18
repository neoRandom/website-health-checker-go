package dto

import "http-server/internal/domain/models"

type GetSiteListResponse struct {
	Body []*SiteResponse `json:"body"`
}

type AddSiteRequest struct {
	Url string `json:"url"`
}

type SiteResponse struct {
	Id  models.SiteID `json:"id"`
	Url string        `json:"url"`
}
