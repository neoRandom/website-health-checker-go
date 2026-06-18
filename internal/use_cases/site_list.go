package usecases

import (
	"http-server/internal/domain/models"
	"http-server/internal/domain/ports/driven"
)

type SiteListUseCases struct {
	SiteRepository driven.SiteRepository
}

func (uc *SiteListUseCases) GetSiteList() ([]*models.Site, error) {
	return uc.SiteRepository.GetList()
}

func (uc *SiteListUseCases) AddSite(url string) (*models.Site, error) {
	s := &models.Site{
		Url: url,
	}

	id, err := uc.SiteRepository.Save(s)
	if err != nil {
		return nil, err
	}
	s.Id = id
	
	return s, nil
}
