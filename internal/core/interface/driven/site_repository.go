package driven

import "http-server/internal/core/model"

type SiteRepository interface {
	GetList() ([]*models.Site, error)
	Save(s *models.Site) (models.SiteID, error)
	Update(s *models.Site) error
	Remove(id models.SiteID) error
	Count() (uint64, error)
}
