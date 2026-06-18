package driven

import "http-server/internal/domain/models"

type SiteRepository interface {
	GetList() ([]*models.Site, error)
	Save(s *models.Site) (models.SiteID, error)
}
