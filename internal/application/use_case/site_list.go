package usecase

import (
	"context"
	"http-server/internal/core/interface/driven"
	"http-server/internal/core/model"
	"net/http"
)

type SiteListUseCases struct {
	siteRepository driven.SiteRepository
}

func NewSiteListUseCases(
	siteRepository driven.SiteRepository,
) *SiteListUseCases {
	return &SiteListUseCases{
		siteRepository: siteRepository,
	}
}

func (uc *SiteListUseCases) GetSiteList(
	ctx context.Context,
) ([]model.Site, error) {
	return uc.siteRepository.GetList(ctx)
}

func (uc *SiteListUseCases) AddSite(
	ctx context.Context, site *model.Site,
) (model.SiteID, error) {
	if site.ExpectedStatusCode == 0 {
		site.ExpectedStatusCode = http.StatusOK
	}

	id, err := uc.siteRepository.Save(ctx, site)
	if err != nil {
		return model.SiteID(0), err
	}

	return model.SiteID(id), nil
}

func (uc *SiteListUseCases) UpdateSite(
	ctx context.Context, site *model.Site,
) error {
	return uc.siteRepository.Update(ctx, site)
}

func (uc *SiteListUseCases) RemoveSite(
	ctx context.Context, id model.SiteID,
) error {
	return uc.siteRepository.Remove(ctx, id)
}
