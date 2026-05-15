package internal

import (
	"context"
	"log"
	"net/http"
	"time"
)

const (
	StatusCodeMethodNotAllowed = 405
	MinimumStatusCodeForError = 400
)


type Watcher struct {
	Targets []Target
}

func CheckHealth(t *Target) {
	log.Printf("watcher: Checking %s (%s) at %s", t.Name, t.ID, t.URL)

	start := time.Now()

	hRes, err := http.Head(t.URL)
	if hRes.StatusCode == StatusCodeMethodNotAllowed {
		hRes, err = http.Get(t.URL)
	}

	dur := time.Since(start)

	cRes := CheckResult{
		// TargetID:       t.ID,
		// URL:            t.URL,
		// HTTPStatusCode: hRes.StatusCode,
		Duration:       dur,
		// CheckedAt:      start,
	}

	if err != nil {
		cRes.Status = CheckStatusUnhealthy
		cRes.Error = err.Error()
	} else {
		if hRes.StatusCode >= MinimumStatusCodeForError {
			cRes.Status = CheckStatusUnhealthy
		} else {
			cRes.Status = CheckStatusHealthy
		}
		cRes.Error = ""
	}

	log.Printf("watcher: %s (%s) is %s - %v", t.Name, t.ID, cRes.Status, cRes.Duration)
}

func (w *Watcher) Watch(ctx context.Context) {
	if len(w.Targets) == 0 {
		return
	}

	c := 0

outer:
	for {
		select {
		case <-ctx.Done():
			break outer
		default:
			if c >= len(w.Targets) {
				c = 0
			}

			time.Sleep(1 * time.Second)

			CheckHealth(&w.Targets[c])
			c++
		}
	}
}
