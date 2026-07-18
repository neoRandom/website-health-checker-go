package usecase

import (
	"context"
	"database/sql"
	"http-server/internal/core/interface/driven"
	"http-server/internal/core/interface/driver"
	"http-server/internal/core/model"
)

type DashboardUseCases struct {
	siteRepository   driven.SiteRepository
	resultRepository driven.ResultRepository
}

func NewDashboardUseCases(
	siteRepository driven.SiteRepository,
	resultRepository driven.ResultRepository,
) *DashboardUseCases {
	return &DashboardUseCases{
		siteRepository:   siteRepository,
		resultRepository: resultRepository,
	}
}

func (uc *DashboardUseCases) GetSiteStatuses(
	ctx context.Context,
) ([]driver.SiteStatus, error) {
	list, err := uc.siteRepository.GetList(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]driver.SiteStatus, 0, len(list))
	for _, s := range list {
		// nil, non-nil-error only on real failure
		latest, err := uc.resultRepository.GetSiteLatest(ctx, s.Id)

		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		out = append(out, driver.SiteStatus{Site: s, Latest: latest})
	}

	return out, nil
}

func (uc *DashboardUseCases) GetSiteDetail(
	ctx context.Context, id model.SiteID, historyLimit int,
) (*driver.SiteStatus, []model.Result, error) {
	site, err := uc.siteRepository.GetByID(ctx, id)
	if err != nil || site == nil {
		return nil, nil, err
	}

	history, err := uc.resultRepository.GetHistory(ctx, id, historyLimit)
	if err != nil {
		return nil, nil, err
	}

	var latest *model.Result
	if len(history) > 0 {
		latest = &history[len(history)-1]
	}
	
	return &driver.SiteStatus{Site: *site, Latest: latest}, history, nil
}
