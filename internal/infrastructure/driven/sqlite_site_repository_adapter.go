package driven

import (
	"context"
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

func (r *SQLiteSiteRepositoryAdapter) GetList(
	ctx context.Context,
) ([]model.Site, error) {
	rows, err := r.db.QueryContext(
		ctx,
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

func (r *SQLiteSiteRepositoryAdapter) GetByID(
	ctx context.Context, id model.SiteID,
) (*model.Site, error) {
	var s model.Site
	
	err := r.db.QueryRowContext(
		ctx,
		`SELECT site_id, url, expected_status_code, description FROM sites WHERE site_id = ?`,
		id,
	).Scan(&s.Id, &s.Url, &s.ExpectedStatusCode, &s.Description)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return &s, nil
}

func (r *SQLiteSiteRepositoryAdapter) Save(
	ctx context.Context, s *model.Site,
) (model.SiteID, error) {
	results, err := r.db.ExecContext(
		ctx,
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

func (r *SQLiteSiteRepositoryAdapter) Update(
	ctx context.Context, s *model.Site,
) error {
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

	_, err := r.db.ExecContext(ctx, finalQuery, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *SQLiteSiteRepositoryAdapter) Remove(
	ctx context.Context, id model.SiteID,
) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM sites WHERE site_id = ?`, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *SQLiteSiteRepositoryAdapter) Count(
	ctx context.Context,
) (uint64, error) {
	var c uint64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sites`).Scan(&c)
	return c, err
}
