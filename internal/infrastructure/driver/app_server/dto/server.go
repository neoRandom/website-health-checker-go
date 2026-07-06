package dto

import "http-server/internal/core/model"

// TODO: change to.AddSiteRequest AND dto.SiteJSON to support ExpectedStatusCode
// As of v-0.0.2, ExpectedStatusCode is optional

type AddSiteRequest struct {
	Url string `json:"url"`
}

type GetSiteListResponse struct {
	Body []SiteJSON `json:"body"`
}

type SiteJSON struct {
	Id  model.SiteID `json:"id"`
	Url string       `json:"url"`
}
