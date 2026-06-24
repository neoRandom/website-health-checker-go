package main

import (
	"context"
	"errors"
	usecase "http-server/internal/application/use_case"
	"http-server/internal/infrastructure/config"
	"http-server/internal/infrastructure/database"
	"http-server/internal/infrastructure/driven"
	appserver "http-server/internal/infrastructure/driver/app_server"
	pprofserver "http-server/internal/infrastructure/driver/pprof_server"
	prometheusmetricsexporter "http-server/internal/infrastructure/driver/prometheus_metrics_exporter"
	"http-server/internal/infrastructure/driver/scheduler"
	"log"

	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "modernc.org/sqlite"
)

// TODO: Improve error handling.
// Change "log.Fatal", which forcibly terminates the program, for a softer warning
// Add error messages and runtime checks

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	//
	db, err := sql.Open(cfg.DatabaseType, cfg.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//
	err = database.Migrate(db)
	if err != nil {
		log.Fatalf("Failed to initialized database schema: %v", err)
	}

	//
	siteRepository := driven.NewSQLiteSiteRepositoryAdapter(db)
	resultRepository := driven.NewSQLiteResultRepositoryAdapter(db)

	metricsCollector := driven.NewPrometheusMetricsCollector(siteRepository)
	httpRequester := driven.NewNetHttpRequesterAdapter()

	siteListUseCases := usecase.NewSiteListUseCases(siteRepository)
	siteCheckUseCases := usecase.NewSiteCheckUseCases(
		httpRequester,
		resultRepository,
	)

	pprofServer := pprofserver.NewPprofServerAdapter(":6060")
	metricsExporter := prometheusmetricsexporter.NewPrometheusMetricsExporterAdapter(
		":2112",
		metricsCollector,
	)
	appServer := appserver.NewAppServerAdapter(
		":8080",
		siteListUseCases.GetSiteList,
		siteListUseCases.AddSite,
		siteListUseCases.UpdateSite,
		siteListUseCases.RemoveSite,
		metricsCollector,
	)
	scheduler := scheduler.NewSchedulerAdapter(
		siteRepository,
		siteCheckUseCases.CheckSites,
	)

	//

	var wg sync.WaitGroup
	errCh := make(chan error, 4)
	startService := func(name string, start func() error) {
		wg.Go(func() {
			if err := start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
				stop()
				return
			}

			log.Printf("%s stopped", name)
		})
	}

	startService("pprof server", func() error { return pprofServer.Start(ctx) })
	startService("metrics exporter", func() error { return metricsExporter.Start(ctx) })
	startService("app server", func() error { return appServer.Start(ctx) })
	startService("scheduler", func() error { return scheduler.Start(ctx) })

	<-ctx.Done()
	log.Printf("Shutdown requested: %v", ctx.Err())

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			log.Printf("service error: %v", err)
		}
	}
}
