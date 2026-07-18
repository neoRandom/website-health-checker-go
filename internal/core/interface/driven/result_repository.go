package driven

import (
	"context"
	"http-server/internal/core/model"
)

type ResultRepository interface {
	GetSiteLatest(ctx context.Context, siteId model.SiteID) (*model.Result, error)
	GetHistory(ctx context.Context, siteId model.SiteID, limit int) ([]model.Result, error)
	Save(ctx context.Context, r *model.Result) (model.ResultID, error)
}
