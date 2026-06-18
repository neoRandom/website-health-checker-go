package driven

import "http-server/internal/domain/models"

type SiteRepository interface {
	GetList() ([]*models.Site, error)
	Save(s *models.Site) (models.SiteID, error)
	Update(s *models.Site) error
	Remove(id models.SiteID) error
}
