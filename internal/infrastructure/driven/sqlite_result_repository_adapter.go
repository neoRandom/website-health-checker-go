package driven

import (
	"database/sql"
	"http-server/internal/core/model"
	"time"
)

type SQLiteResultRepositoryAdapter struct {
	db *sql.DB
}

func NewSQLiteResultRepositoryAdapter(db *sql.DB) *SQLiteResultRepositoryAdapter {
	return &SQLiteResultRepositoryAdapter{
		db: db,
	}
}

func (a *SQLiteResultRepositoryAdapter) GetSiteLatest(siteId model.SiteID) (*model.Result, error) {
	var r model.Result
	var responseTimeMs, unixCheckedAt int64

	err := a.db.QueryRow(
		`
			SELECT
				result_id, 
				site_id, 
				status_code, 
				is_healthy, 
				is_secure, 
				response_time_ms, 
				checked_at 
			FROM results 
			WHERE site_id = ?
			ORDER BY checked_at DESC
			LIMIT 1
		`,
		siteId,
	).Scan(
		&r.Id,
		&r.SiteId,
		&r.StatusCode,
		&r.IsHealthy,
		&r.IsSecure,
		&responseTimeMs,
		&unixCheckedAt,
	)

	if err != nil {
		return nil, err
	}
	r.ResponseTime = time.Duration(responseTimeMs) * time.Millisecond
	r.CheckedAt = time.Unix(unixCheckedAt, 0)

	return &r, nil
}

func (a *SQLiteResultRepositoryAdapter) GetHistory(
	siteId model.SiteID, limit int,
) ([]model.Result, error) {
	rows, err := a.db.Query(
		`
			SELECT
				result_id, site_id, status_code, is_healthy,
				is_secure, response_time_ms, checked_at
			FROM results
			WHERE site_id = ?
			ORDER BY checked_at DESC
			LIMIT ?
		`,
		siteId, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Result
	for rows.Next() {
		var r model.Result
		var responseTimeMs, unixCheckedAt int64
		
		if err := rows.Scan(
			&r.Id, &r.SiteId, &r.StatusCode, &r.IsHealthy,
			&r.IsSecure, &responseTimeMs, &unixCheckedAt,
		); err != nil {
			return nil, err
		}
		
		r.ResponseTime = time.Duration(responseTimeMs) * time.Millisecond
		r.CheckedAt = time.Unix(unixCheckedAt, 0)
		out = append(out, r)
	}

	// DESC query, but the sparkline wants oldest-first — reverse in place.
	// Can't be solved by a simple "ASC"
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out, rows.Err()
}

func (a *SQLiteResultRepositoryAdapter) Save(r *model.Result) (model.ResultID, error) {
	results, err := a.db.Exec(
		`
			INSERT INTO results
			(
				site_id, 
				status_code, 
				is_healthy, 
				is_secure, 
				response_time_ms, 
				checked_at
			) 
			VALUES (?, ?, ?, ?, ?, ?)
		`,
		r.SiteId,
		r.StatusCode,
		r.IsHealthy,
		r.IsSecure,
		r.ResponseTime.Milliseconds(),
		r.CheckedAt.Unix(),
	)
	if err != nil {
		return model.ResultID(0), err
	}

	id, err := results.LastInsertId()
	return model.ResultID(id), err
}
