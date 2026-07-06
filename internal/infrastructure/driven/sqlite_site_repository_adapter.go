package driven

import (
	"database/sql"
	"http-server/internal/core/model"
)

type SQLiteSiteRepositoryAdapter struct {
	db *sql.DB
}

func NewSQLiteSiteRepositoryAdapter(db *sql.DB) *SQLiteSiteRepositoryAdapter {
	return &SQLiteSiteRepositoryAdapter{
		db: db,
	}
}

func (r *SQLiteSiteRepositoryAdapter) GetList() ([]model.Site, error) {
	rows, err := r.db.Query(
		`SELECT site_id, url, expected_status_code FROM sites`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Site

	for rows.Next() {
		var site model.Site
		if err := rows.Scan(
			&site.Id,
			&site.Url,
			&site.ExpectedStatusCode,
		); err != nil {
			return nil, err
		}
		list = append(list, site)
	}

	return list, nil
}

func (r *SQLiteSiteRepositoryAdapter) Save(s *model.Site) (model.SiteID, error) {
	results, err := r.db.Exec(
		`
			INSERT INTO sites (url, expected_status_code) 
			VALUES (?, ?)
		`,
		s.Url, s.ExpectedStatusCode,
	)
	if err != nil {
		return model.SiteID(0), err
	}

	id, err := results.LastInsertId()
	return model.SiteID(id), err
}

// TODO: refactor the Update implementation
// currently, the UPDATE updates every column except the primary key
// change it to dynamically build the query, or skip it entirely
func (r *SQLiteSiteRepositoryAdapter) Update(s *model.Site) error {
	_, err := r.db.Exec(
		`UPDATE sites SET url = ?, expected_status_code = ? WHERE site_id = ?`,
		s.Url, s.ExpectedStatusCode, s.Id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *SQLiteSiteRepositoryAdapter) Remove(id model.SiteID) error {
	_, err := r.db.Exec(`DELETE FROM sites WHERE site_id = ?`, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *SQLiteSiteRepositoryAdapter) Count() (uint64, error) {
	var c uint64
	err := r.db.QueryRow(`SELECT COUNT(*) FROM sites`).Scan(&c)
	return c, err
}
