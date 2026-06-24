package scheduler

import (
	"http-server/internal/core/interface/driven"
	"http-server/internal/core/interface/driver"
	"log"
	"time"
)

type SchedulerAdapter struct {
	siteRepository driven.SiteRepository
	checkSites     driver.CheckSites
}

func NewSchedulerAdapter(
	siteRepository driven.SiteRepository,
	checkSites driver.CheckSites,
) *SchedulerAdapter {
	return &SchedulerAdapter{
		siteRepository: siteRepository,
		checkSites:     checkSites,
	}
}

func (a *SchedulerAdapter) Start() error {
	targets, err := a.siteRepository.GetList()
	if err != nil {
		return err
	}

	log.Printf("Scheduler starting with %d targets", len(targets))

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				log.Printf("Scheduler stopping")
				return

			case <-ticker.C:
				log.Printf("Scheduler tick: checking %d targets", len(targets))
				if err := a.checkSites(targets); err != nil {
					log.Printf("Scheduler check failed: %s", err)
				}
			}
		}
	}()

	// TODO: Implement graceful shutdown (see main.go)
	// When adding the context, change the done channel to stop the ticker as
	// the context stops too
	//
	// done<- true
	return nil
}
