package usecase

import (
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

func (uc *SiteCheckUseCases) CheckSites(targets []model.Site) error {
	var wg sync.WaitGroup

	for _, t := range targets {
		wg.Add(1)
		go uc.worker(&t, &wg)
	}

	wg.Wait()

	return nil
}

func (uc *SiteCheckUseCases) worker(target *model.Site, wg *sync.WaitGroup) {
	defer wg.Done()

	results, err := uc.httpRequester.CheckSite(target)
	if err != nil {
		log.Printf("Error worker checking site '%s': %s", target.Url, err)
		return
	}

	_, err = uc.resultRepository.Save(results)
	if err != nil {
		log.Printf("Error worker saving site '%s': %s", target.Url, err)
		return
	}
}
