package driver

import (
	"context"
	"http-server/internal/core/model"
)

type SiteStatus struct {
	Site   model.Site
	Latest *model.Result // nil if no check has run yet
}

type GetSiteStatuses func(ctx context.Context) ([]SiteStatus, error)
type GetSiteDetail (
	func(
		ctx context.Context, 
		id model.SiteID, historyLimit int,
	) (*SiteStatus, []model.Result, error))
