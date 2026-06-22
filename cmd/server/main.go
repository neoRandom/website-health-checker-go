package main

import (
	"database/sql"
	"http-server/internal/infrastructure/driven"
	appserver "http-server/internal/infrastructure/driver/app_server"
	"http-server/internal/infrastructure/config"
	"http-server/internal/infrastructure/database"
	usecase "http-server/internal/application/use_case"
	"log"

	"net/http"
	_ "net/http/pprof"

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

	server := appserver.NewAppServerAdapter(
		":8080",
		siteListUseCase.GetSiteList,
		siteListUseCase.AddSite,
		siteListUseCase.UpdateSite,
		siteListUseCase.RemoveSite,
		metricsCollector,
	)

	//
	go func() {
		log.Println("Initializing pprof at http://localhost:6060...")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("Error initializing pprof: %v", err)
		}
	}()

	server.Run()
}
