package web

import (
	"fmt"
	"net/url"
	"time"
)

// State is derived, never stored — it exists only as a function of
// (Result.StatusCode, Result.IsHealthy, Site.ExpectedStatusCode).
type State string

const (
	StateOK    State = "OK"
	StateWatch State = "WATCH"
	StateDown  State = "DOWN"
)

// DeriveState implements the exact rule from the spec:
//   DOWN  = 4xx/5xx or transport failure (IsHealthy false)
//   WATCH = healthy response, but StatusCode != ExpectedStatusCode
//   OK    = healthy response matching ExpectedStatusCode
func DeriveState(isHealthy bool, statusCode, expectedStatusCode int) State {
	if !isHealthy || statusCode >= 400 {
		return StateDown
	}
	if statusCode != expectedStatusCode {
		return StateWatch
	}
	return StateOK
}

func (s State) TextClass() string {
	switch s {
	case StateOK:
		return "text-emerald-400"
	case StateWatch:
		return "text-amber-400"
	case StateDown:
		return "text-red-400"
	default:
		return "text-neutral-500"
	}
}

func (s State) DotClass() string {
	switch s {
	case StateOK:
		return "bg-emerald-400"
	case StateWatch:
		return "bg-amber-400"
	case StateDown:
		return "bg-red-400"
	default:
		return "bg-neutral-600"
	}
}

func (s State) BarClass() string {
	switch s {
	case StateOK:
		return "bg-emerald-500/70 hover:bg-emerald-400"
	case StateWatch:
		return "bg-amber-500/70 hover:bg-amber-400"
	case StateDown:
		return "bg-red-500/70 hover:bg-red-400"
	default:
		return "bg-neutral-700"
	}
}

// SplitURL is the fix for the host/endpoint conflation: Site.Url is the
// only source of truth, so host and path are derived here, once, at the
// adapter boundary — templates never parse URLs themselves.
//
// This does NOT return a "secure" flag. Whether a check's TLS/SSL
// certificate is present and valid is a dynamic, per-Result finding
// (Result.IsSecure) established by the checker at request time — it is
// not something you can infer from the scheme in a stored URL string.
// A site can be configured https:// and still report IsSecure=false on
// a given check (expired cert, hostname mismatch, chain failure) while
// the HTTP request itself still succeeds.
func SplitURL(raw string) (host string, endpoint string) {
	u, err := url.Parse(raw)
	if err != nil {
		return raw, "/"
	}
	host = u.Host
	if host == "" {
		host = u.Path // fallback for malformed/relative input
	}
	endpoint = u.Path
	if endpoint == "" {
		endpoint = "/"
	}
	return host, endpoint
}

// RelativeTime renders a Duration-since as the compact "12s / 3m / 2h" form
// used throughout the dashboard, so this formatting rule exists in exactly
// one place.
func RelativeTime(t time.Time, now time.Time) string {
	d := now.Sub(t)
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
}

// toStr is the single int->string conversion path used from templ
// expressions (templ interpolation requires a string, not a Stringer).
func toStr(i int) string {
	return fmt.Sprintf("%d", i)
}

// Secure is tri-state, not bool: TLS/SSL validity is established by a
// check's Result, so a never-checked site has no data yet — that is
// "unknown", not "insecure". Collapsing the two would be a false
// positive in the UI (a red/absent badge implying a cert problem that
// was never actually observed).
type Secure int

const (
	SecureUnknown Secure = iota
	SecureYes
	SecureNo
)

func SecureFromResult(hasResult bool, isSecure bool) Secure {
	if !hasResult {
		return SecureUnknown
	}
	if isSecure {
		return SecureYes
	}
	return SecureNo
}

func (s Secure) Label() string {
	switch s {
	case SecureYes:
		return "SECURE"
	case SecureNo:
		return "INSECURE"
	default:
		return "—"
	}
}

func (s Secure) TextClass() string {
	switch s {
	case SecureYes:
		return "text-emerald-400"
	case SecureNo:
		return "text-red-400"
	default:
		return "text-neutral-600"
	}
}

// --- Row-level view models -------------------------------------------------

type SiteRow struct {
	ID                 string
	Host               string
	Endpoint           string
	Description        string
	State              State
	LastResponseTimeMS int64 // -1 when no successful response exists
	LastCheckedAgo     string
	Secure             Secure // from the latest Result, not the URL scheme
	ExpectedStatusCode int
	LastStatusCode     int
}

type IncidentRow struct {
	Host           string
	Endpoint       string
	ResponseTimeMS int64
	CheckedAgo     string
	StatusCode     int
	State          State
}

type ResultPoint struct {
	ResponseTimeMS  int64
	State           State
	StatusCode      int
	CheckedAgo      string
	BarHeightPct    float64 // 0-100, pre-computed against the window max
}

type SiteDetail struct {
	ID                 string
	Host               string
	Endpoint           string
	Description        string
	ExpectedStatusCode int
	CurrentSecure      Secure // from the most recent Result, not the URL scheme
	CurrentState       State
	Results            []ResultPoint // chronological, oldest first, capped at 100
}

// --- Page-level view model --------------------------------------------------

type DashboardData struct {
	TotalSites int
	OKCount    int
	DownCount  int
	WatchCount int

	Method                 string // global constant — no per-site field exists
	PollingIntervalSeconds int
	TimeoutSeconds         int
	LastCheckAgo           string

	Sites     []SiteRow
	Incidents []IncidentRow

	Selected *SiteDetail // nil until a site row is chosen via HTMX
}

// BuildResultPoints caps to the last 100 results, orders oldest-first, and
// pre-computes bar heights relative to the max response time in the window
// so the template performs zero arithmetic.
func BuildResultPoints(results []ResultPoint) []ResultPoint {
	var maxMS int64 = 1
	for _, r := range results {
		if r.ResponseTimeMS > maxMS {
			maxMS = r.ResponseTimeMS
		}
	}
	for i := range results {
		if results[i].State == StateDown {
			results[i].BarHeightPct = 12 // floor height for no-response bars
			continue
		}
		pct := float64(results[i].ResponseTimeMS) / float64(maxMS) * 100
		if pct < 8 {
			pct = 8
		}
		results[i].BarHeightPct = pct
	}
	if len(results) > 100 {
		results = results[len(results)-100:]
	}
	return results
}