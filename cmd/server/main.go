package main

import (
	"context"
	"errors"
	"fmt"
	usecase "http-server/internal/application/use_case"
	"http-server/internal/infrastructure/config"
	"http-server/internal/infrastructure/database"
	"http-server/internal/infrastructure/driven"
	appserver "http-server/internal/infrastructure/driver/app_server"
	pprofserver "http-server/internal/infrastructure/driver/pprof_server"
	prometheusmetricsexporter "http-server/internal/infrastructure/driver/prometheus_metrics_exporter"
	"http-server/internal/infrastructure/driver/scheduler"
	"log"
	"time"

	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "modernc.org/sqlite"
)

func main() {
	if err := run(); err != nil {
		log.Printf("application error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(
		context.Background(), os.Interrupt, syscall.SIGTERM,
	)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// ====

	// TODO: Change hard-coded values for env variables
	dataSource := fmt.Sprintf(
		"file:%s?_journal_mode=%s&_busy_timeout=%d",
		cfg.DatabasePath, "WAL", 10000,
	)

	db, err := sql.Open(cfg.DatabaseType, dataSource)
	if err != nil {
		return fmt.Errorf(
			"open database %q at %q: %w",
			cfg.DatabaseType, cfg.DatabasePath, err,
		)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("database close error: %v", closeErr)
		}
	}()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database at %q: %w", cfg.DatabasePath, err)
	}

	err = database.Migrate(db)
	if err != nil {
		return fmt.Errorf("migrate database schema: %w", err)
	}

	// ====

	siteRepository := driven.NewSQLiteSiteRepositoryAdapter(db)
	resultRepository := driven.NewSQLiteResultRepositoryAdapter(db)

	metricsCollector := driven.NewPrometheusMetricsCollector(
		ctx, siteRepository,
	)
	httpRequester := driven.NewNetHttpRequesterAdapter(
		time.Duration(cfg.CheckTimeout) * time.Second,
	)

	siteListUseCases := usecase.NewSiteListUseCases(siteRepository)
	siteCheckUseCases := usecase.NewSiteCheckUseCases(
		httpRequester,
		resultRepository,
	)
	dashboardUseCases := usecase.NewDashboardUseCases(
		siteRepository,
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
		dashboardUseCases.GetSiteStatuses,
		dashboardUseCases.GetSiteDetail,
		metricsCollector,
		&cfg,
	)
	scheduler := scheduler.NewSchedulerAdapter(
		time.Duration(cfg.CheckInterval)*time.Second,
		siteRepository,
		siteCheckUseCases.CheckSites,
	)

	// ====

	var wg sync.WaitGroup
	errCh := make(chan error, 4)
	startService := func(name string, start func(ctx context.Context) error) {
		wg.Go(func() {
			if err := start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
				wrappedErr := fmt.Errorf(
					"%s failed to start or stopped unexpectedly: %w",
					name, err,
				)
				log.Printf("%v", wrappedErr)
				errCh <- wrappedErr
				stop()
				return
			}

			log.Printf("%s stopped", name)
		})
	}

	startService("pprof server", pprofServer.Start)
	startService("metrics exporter", metricsExporter.Start)
	startService("app server", appServer.Start)
	startService("scheduler", scheduler.Start)

	<-ctx.Done()
	log.Printf("Shutdown requested: %v", ctx.Err())

	wg.Wait()
	close(errCh)

	var serviceErr error

	for err := range errCh {
		if err != nil {
			log.Printf("service error: %v", err)
			if serviceErr == nil {
				serviceErr = err
			}
		}
	}

	if serviceErr != nil {
		return fmt.Errorf("shutdown after service failure: %w", serviceErr)
	}

	return nil
}
