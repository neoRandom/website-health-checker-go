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
	var responseTimeMs int64
	var unixCheckedAt int64

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
