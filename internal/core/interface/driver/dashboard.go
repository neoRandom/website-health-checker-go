package driver

import "http-server/internal/core/model"

type SiteStatus struct {
	Site   model.Site
	Latest *model.Result // nil if no check has run yet
}

type GetSiteStatuses func() ([]SiteStatus, error)
type GetSiteDetail func(id model.SiteID, historyLimit int) (*SiteStatus, []model.Result, error)
