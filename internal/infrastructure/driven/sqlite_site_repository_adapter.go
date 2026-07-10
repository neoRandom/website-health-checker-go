package driven

import (
	"database/sql"
	"fmt"
	"http-server/internal/core/model"
	"strings"
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
		`SELECT site_id, url, expected_status_code, description FROM sites`,
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
			&site.Description,
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
			INSERT INTO sites (url, expected_status_code, description) 
			VALUES (?, ?, ?)
		`,
		s.Url, s.ExpectedStatusCode, s.Description,
	)
	if err != nil {
		return model.SiteID(0), err
	}

	id, err := results.LastInsertId()
	return model.SiteID(id), err
}

func (r *SQLiteSiteRepositoryAdapter) Update(s *model.Site) error {
	var queryParts []string
	var args []any
	argCounter := 1

	if s.Id == model.SiteID(0) {
		return fmt.Errorf("Invalid site ID '%v'", s.Id)
	}

	if s.Url != "" {
		queryParts = append(
			queryParts, fmt.Sprintf("url = $%d", argCounter),
		)
		args = append(args, s.Url)
		argCounter++
	}
	if s.ExpectedStatusCode != 0 {
		queryParts = append(
			queryParts, fmt.Sprintf("expected_status_code = $%d", argCounter),
		)
		args = append(args, s.ExpectedStatusCode)
		argCounter++
	}
	if s.Description != "" {
		queryParts = append(
			queryParts, fmt.Sprintf("description = $%d", argCounter),
		)
		args = append(args, s.Description)
		argCounter++
	}

	if len(queryParts) == 0 {
		return nil
	}

	args = append(args, s.Id)
	finalQuery := fmt.Sprintf(
		`UPDATE sites SET %s WHERE site_id = $%d`, 
		strings.Join(queryParts, ", "),
		argCounter,
	)

	_, err := r.db.Exec(finalQuery, args...)
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
