package usecase

import (
	"context"
	"http-server/internal/core/interface/driven"
	"http-server/internal/core/model"
	"log"
	"sync"
)

type SiteCheckUseCases struct {
	httpRequester driven.HTTPRequester
	resultRepository driven.ResultRepository
}

func NewSiteCheckUseCases(
	httpRequester driven.HTTPRequester,
	resultRepository driven.ResultRepository,
) *SiteCheckUseCases {
	return &SiteCheckUseCases{
		httpRequester: httpRequester,
		resultRepository: resultRepository,
	}
}

func (uc *SiteCheckUseCases) CheckSites(
	ctx context.Context, targets []model.Site,
) error {
	var wg sync.WaitGroup

	for _, t := range targets {
		wg.Add(1)
		go uc.worker(ctx, &t, &wg)
	}

	wg.Wait()

	return nil
}

func (uc *SiteCheckUseCases) worker(
	ctx context.Context, target *model.Site, wg *sync.WaitGroup,
) {
	defer wg.Done()

	results, err := uc.httpRequester.CheckSite(target)
	if err != nil {
		log.Printf("Error worker checking site '%s': %s", target.Url, err)
		return
	}

	_, err = uc.resultRepository.Save(ctx, results)
	if err != nil {
		log.Printf("Error worker saving site '%s': %s", target.Url, err)
		return
	}
}
