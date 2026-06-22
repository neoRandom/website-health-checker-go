package server

import (
	"log"
	"http-server/internal/infrastructure/config"
	"http-server/internal/infrastructure/database"
	"http-server/internal/infrastructure/driven"
	usecase "http-server/internal/application/use_case"
	pprofserver "http-server/internal/infrastructure/driver/pprof_server"
	appserver "http-server/internal/infrastructure/driver/app_server"
	
	"database/sql"
	_ "modernc.org/sqlite"
)

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
	metricsCollector := driven.NewPrometheusMetricsCollector(siteRepository)
	siteListUseCase := usecase.NewSiteListUseCases(siteRepository)

	pprofServer := pprofserver.NewPprofServerAdapter()
	appServer := appserver.NewAppServerAdapter(
		":8080",
		siteListUseCase.GetSiteList,
		siteListUseCase.AddSite,
		siteListUseCase.UpdateSite,
		siteListUseCase.RemoveSite,
		metricsCollector,
	)

	//
	go pprofServer.Start()
	go appServer.Start()
}
