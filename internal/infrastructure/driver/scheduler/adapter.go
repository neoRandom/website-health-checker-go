package scheduler

import (
	"context"
	"fmt"
	"http-server/internal/core/interface/driven"
	"http-server/internal/core/interface/driver"
	"log"
	"time"
)

type SchedulerAdapter struct {
	interval       time.Duration
	siteRepository driven.SiteRepository
	checkSites     driver.CheckSites
}

func NewSchedulerAdapter(
	interval time.Duration,
	siteRepository driven.SiteRepository,
	checkSites driver.CheckSites,
) *SchedulerAdapter {
	return &SchedulerAdapter{
		interval: interval,
		siteRepository: siteRepository,
		checkSites:     checkSites,
	}
}

func (a *SchedulerAdapter) Start(ctx context.Context) error {
	targets, err := a.siteRepository.GetList()
	if err != nil {
		return fmt.Errorf("load scheduled targets: %w", err)
	}

	log.Printf("Scheduler starting with %d targets", len(targets))

	ticker := time.NewTicker(a.interval)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Printf("Scheduler stopping")
				return

			case <-ticker.C:
				updatedTargets, err := a.siteRepository.GetList()
				if err != nil {
					log.Printf("Scheduler could not refresh targets: %v", fmt.Errorf("refresh scheduled targets: %w", err))
				} else {
					targets = updatedTargets
				}

				log.Printf("Scheduler tick: checking %d targets", len(targets))
				if err := a.checkSites(targets); err != nil {
					log.Printf("Scheduler check failed: %v", fmt.Errorf("check scheduled targets: %w", err))
				}
			}
		}
	}()

	<-ctx.Done()
	return nil
}
