package database

import "database/sql"

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
		PRAGMA foreign_keys = ON;

		CREATE TABLE IF NOT EXISTS sites (
			site_id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT NOT NULL UNIQUE
		);

		CREATE TABLE IF NOT EXISTS results (
			result_id INTEGER PRIMARY KEY AUTOINCREMENT,
			site_id INTEGER NOT NULL,
			status_code INTEGER NOT NULL,
			checked_at INTEGER NOT NULL,

			FOREIGN KEY (site_id) REFERENCES sites(site_id) ON DELETE CASCADE
		);
	`)
	
	return err
}
