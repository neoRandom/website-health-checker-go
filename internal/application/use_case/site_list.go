package usecase

import (
	"http-server/internal/core/model"
	"http-server/internal/core/interface/driven"
)

type SiteListUseCases struct {
	siteRepository driven.SiteRepository
}

func NewSiteListUseCases(siteRepository driven.SiteRepository) *SiteListUseCases {
	return &SiteListUseCases{
		siteRepository: siteRepository,
	}
}

func (uc *SiteListUseCases) GetSiteList() ([]*models.Site, error) {
	return uc.siteRepository.GetList()
}

func (uc *SiteListUseCases) AddSite(url string) (*models.Site, error) {
	s := &models.Site{
		Url: url,
	}

	id, err := uc.siteRepository.Save(s)
	if err != nil {
		return nil, err
	}
	s.Id = id

	return s, nil
}

func (uc *SiteListUseCases) UpdateSite(id models.SiteID, url string) error {
	s := &models.Site{
		Id:  id,
		Url: url,
	}

	return uc.siteRepository.Update(s)
}

func (uc *SiteListUseCases) RemoveSite(id models.SiteID) error {
	return uc.siteRepository.Remove(id)
}
