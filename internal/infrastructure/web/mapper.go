package web

import (
	"http-server/internal/core/interface/driver"
	"http-server/internal/core/model"
	"http-server/internal/infrastructure/config"
	"strconv"
	"time"
)

func BuildDashboardData(
	cfg *config.Config, statuses []driver.SiteStatus, now time.Time,
) DashboardData {
	data := DashboardData{
		Method:                 "GET/HEAD",
		PollingIntervalSeconds: cfg.CheckInterval,
		TimeoutSeconds:         cfg.CheckTimeout,
		TotalSites:             len(statuses),
	}

	var mostRecentCheck time.Time
	var incidents []IncidentRow

	for _, st := range statuses {
		host, endpoint := SplitURL(st.Site.Url)

		row := SiteRow{
			ID:          strconv.FormatInt(int64(st.Site.Id), 10),
			Host:        host,
			Endpoint:    endpoint,
			Description: st.Site.Description,
			Secure: SecureFromResult(
				st.Latest != nil,
				st.Latest != nil && st.Latest.IsSecure,
			),
			ExpectedStatusCode: st.Site.ExpectedStatusCode,
			LastResponseTimeMS: -1,
		}

		if st.Latest != nil {
			latest := st.Latest
			state := DeriveState(
				latest.IsHealthy,
				latest.StatusCode,
				st.Site.ExpectedStatusCode,
			)

			row.State = state
			row.LastStatusCode = latest.StatusCode
			row.LastCheckedAgo = RelativeTime(latest.CheckedAt, now)

			if latest.IsHealthy {
				row.LastResponseTimeMS = latest.ResponseTime.Milliseconds()
			}
			if latest.CheckedAt.After(mostRecentCheck) {
				mostRecentCheck = latest.CheckedAt
			}

			switch state {
			case StateOK:
				data.OKCount++
			case StateWatch:
				data.WatchCount++
			case StateDown:
				data.DownCount++
			}

			if state != StateOK {
				incidents = append(incidents, IncidentRow{
					Host:           host,
					Endpoint:       endpoint,
					ResponseTimeMS: latest.ResponseTime.Milliseconds(),
					CheckedAgo:     RelativeTime(latest.CheckedAt, now),
					StatusCode:     latest.StatusCode,
					State:          state,
				})
			}
		}

		data.Sites = append(data.Sites, row)
	}

	data.Incidents = incidents
	data.LastCheckAgo = RelativeTime(mostRecentCheck, now)
	return data
}

// BuildSiteDetail is the presenter for the GET /sites/{id} HTMX partial.
func BuildSiteDetail(site model.Site, history []model.Result, now time.Time) SiteDetail {
	host, endpoint := SplitURL(site.Url)

	points := make([]ResultPoint, 0, len(history))
	for _, r := range history {
		state := DeriveState(
			r.IsHealthy,
			r.StatusCode,
			site.ExpectedStatusCode,
		)

		var ms int64
		if r.IsHealthy {
			ms = r.ResponseTime.Milliseconds()
		}

		points = append(points, ResultPoint{
			ResponseTimeMS: ms,
			State:          state,
			StatusCode:     r.StatusCode,
			CheckedAgo:     RelativeTime(r.CheckedAt, now),
		})
	}

	points = BuildResultPoints(points)

	var current State = StateOK
	var currentSecure Secure = SecureUnknown

	if len(points) > 0 {
		current = points[len(points)-1].State
	}
	if len(history) > 0 {
		latest := history[len(history)-1]
		currentSecure = SecureFromResult(true, latest.IsSecure)
	}

	return SiteDetail{
		ID:                 strconv.FormatInt(int64(site.Id), 10),
		Host:               host,
		Endpoint:           endpoint,
		Description:        site.Description,
		ExpectedStatusCode: site.ExpectedStatusCode,
		CurrentSecure:      currentSecure,
		CurrentState:       current,
		Results:            points,
	}
}
