package main

import (
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

	_ "modernc.org/sqlite"
)

// TODO: Change "[]*Struct" to "[]Struct" where it is applicable
// A slice of structures has poorer performance because of multiple reasons
// So it needs to be used only where pointers are absolutely necessary

// TODO: Improve error handling.
// Change "log.Fatal", which forcibly terminates the program, for a softer warning
// Add error messages and runtime checks

func main() {
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

	// TODO: Implement graceful shutdown. If not, storage may corrupt.
	go pprofServer.Start()
	go metricsExporter.Start()
	go appServer.Start()
	go scheduler.Start()
}
