package driver

import (
	"context"
	"http-server/internal/core/model"
)

type GetSiteList func(ctx context.Context) ([]model.Site, error)
type AddSite func(ctx context.Context, site *model.Site) (model.SiteID, error)
type UpdateSite func(ctx context.Context, site *model.Site) error
type RemoveSite func(ctx context.Context, id model.SiteID) error
