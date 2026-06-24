package model

import "time"

type ResultID int64

type Result struct {
	Id         ResultID
	SiteId     SiteID
	StatusCode int
	CheckedAt  time.Time
}
