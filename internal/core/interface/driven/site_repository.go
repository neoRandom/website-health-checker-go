package driven

import (
	"context"
	"http-server/internal/core/model"
)

type SiteRepository interface {
	GetList(ctx context.Context) ([]model.Site, error)
	GetByID(ctx context.Context, id model.SiteID) (*model.Site, error)
	Save(ctx context.Context, s *model.Site) (model.SiteID, error)
	Update(ctx context.Context, s *model.Site) error
	Remove(ctx context.Context, id model.SiteID) error
	Count(ctx context.Context) (uint64, error)
}
