package driven

import "http-server/internal/core/model"

type SiteRepository interface {
	GetList() ([]*model.Site, error)
	Save(s *model.Site) (model.SiteID, error)
	Update(s *model.Site) error
	Remove(id model.SiteID) error
	Count() (uint64, error)
}
