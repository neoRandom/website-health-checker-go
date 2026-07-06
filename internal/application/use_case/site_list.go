package usecase

import (
	"http-server/internal/core/interface/driven"
	"http-server/internal/core/model"
	"net/http"
)

type SiteListUseCases struct {
	siteRepository driven.SiteRepository
}

func NewSiteListUseCases(siteRepository driven.SiteRepository) *SiteListUseCases {
	return &SiteListUseCases{
		siteRepository: siteRepository,
	}
}

func (uc *SiteListUseCases) GetSiteList() ([]model.Site, error) {
	return uc.siteRepository.GetList()
}

func (uc *SiteListUseCases) AddSite(site *model.Site) (model.SiteID, error) {
	if site.ExpectedStatusCode == 0 {
		site.ExpectedStatusCode = http.StatusOK
	}

	id, err := uc.siteRepository.Save(site)
	if err != nil {
		return model.SiteID(0), err
	}

	return model.SiteID(id), nil
}

func (uc *SiteListUseCases) UpdateSite(site *model.Site) error {
	return uc.siteRepository.Update(site)
}

func (uc *SiteListUseCases) RemoveSite(id model.SiteID) error {
	return uc.siteRepository.Remove(id)
}
