package internal

import "time"

type CheckStatus string

const (
    CheckStatusUnknown   CheckStatus = "unknown"
    CheckStatusHealthy   CheckStatus = "healthy"
    CheckStatusUnhealthy CheckStatus = "unhealthy"
)

type Target struct {
    ID            string
    Name          string
    URL           string
    CreatedAt     time.Time
    UpdatedAt     time.Time
    LastCheckedAt *time.Time
    LastStatus    CheckStatus
    LastError     string
}

type CheckResult struct {
    TargetID       string
    URL            string
    Status         CheckStatus
    HTTPStatusCode int
    Duration       time.Duration
    CheckedAt      time.Time
    Error          string
}
