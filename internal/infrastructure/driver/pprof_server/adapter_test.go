package pprofserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	pprof "net/http/pprof"
)

func TestPprofServerRoutes(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(testMux())
	t.Cleanup(server.Close)

	t.Run("root redirects to pprof index", func(t *testing.T) {
		t.Parallel()

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		resp, err := client.Get(server.URL + "/")
		if err != nil {
			t.Fatalf("GET / failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusFound {
			t.Fatalf("expected status %d, got %d", http.StatusFound, resp.StatusCode)
		}

		location := resp.Header.Get("Location")
		if location != "/debug/pprof/" {
			t.Fatalf("expected redirect to /debug/pprof/, got %q", location)
		}
	})

	t.Run("pprof index is served", func(t *testing.T) {
		t.Parallel()

		resp, err := http.Get(server.URL + "/debug/pprof/")
		if err != nil {
			t.Fatalf("GET /debug/pprof/ failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("read response body failed: %v", err)
		}

		if len(body) == 0 {
			t.Fatal("expected non-empty pprof index body")
		}
	})

	t.Run("pprof cmdline is served", func(t *testing.T) {
		t.Parallel()

		resp, err := http.Get(server.URL + "/debug/pprof/cmdline")
		if err != nil {
			t.Fatalf("GET /debug/pprof/cmdline failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("pprof heap is served", func(t *testing.T) {
		t.Parallel()

		resp, err := http.Get(server.URL + "/debug/pprof/heap")
		if err != nil {
			t.Fatalf("GET /debug/pprof/heap failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})
}

func testMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/debug/pprof/", http.StatusFound)
	})

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))

	return mux
}
