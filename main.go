package main

import (
	"database/sql"
	"http-server/database"
	"http-server/internal/adapter/driven"
	"http-server/internal/adapter/driver"
	"http-server/internal/config"
	usecases "http-server/internal/use_cases"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open(cfg.DatabaseType, cfg.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = database.Migrate(db)
	if err != nil {
		log.Fatalf("Failed to initialized database schema: %v", err)
	}

	siteRepository := driven.NewSQLiteSiteRepositoryAdapter(db)

	siteListUseCase := usecases.NewSiteListUseCases(siteRepository)

	server := driver.NewServerAdapter(
		"localhost:8080",
		siteListUseCase.GetSiteList,
		siteListUseCase.AddSite,
		siteListUseCase.UpdateSite,
		siteListUseCase.RemoveSite,
	)

	server.Run()
}
