package model

import "time"

type ResultID int64

type Result struct {
	Id           ResultID
	SiteId       SiteID
	StatusCode   int
	IsHealthy    bool
	IsSecure     bool
	ResponseTime time.Duration
	CheckedAt    time.Time
}
