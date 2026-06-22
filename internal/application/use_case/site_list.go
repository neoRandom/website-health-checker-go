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

func (uc *SiteListUseCases) GetSiteList() ([]*model.Site, error) {
	return uc.siteRepository.GetList()
}

func (uc *SiteListUseCases) AddSite(url string) (*model.Site, error) {
	s := &model.Site{
		Url: url,
	}

	id, err := uc.siteRepository.Save(s)
	if err != nil {
		return nil, err
	}
	s.Id = id

	return s, nil
}

func (uc *SiteListUseCases) UpdateSite(id model.SiteID, url string) error {
	s := &model.Site{
		Id:  id,
		Url: url,
	}

	return uc.siteRepository.Update(s)
}

func (uc *SiteListUseCases) RemoveSite(id model.SiteID) error {
	return uc.siteRepository.Remove(id)
}
