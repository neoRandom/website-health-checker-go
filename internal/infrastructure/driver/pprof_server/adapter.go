package pprofserver

import (
	"context"
	"log"
	"net/http"
	pprof "net/http/pprof"
	"time"
)

type PprofServerAdapter struct {
	addr string
}

func NewPprofServerAdapter(addr string) *PprofServerAdapter {
	return &PprofServerAdapter{
		addr: addr,
	}
}

func (ps *PprofServerAdapter) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("pprof server received request for %s", r.URL.Path)
		http.Redirect(w, r, "/debug/pprof/", http.StatusFound)
	})

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	srv := http.Server{
		Addr:    ps.addr,
		Handler: mux,
	}

	log.Printf("pprof server starting at http://localhost%v...", ps.addr)
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Printf("pprof server stopping at http://localhost%v...", ps.addr)
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return err
		}
		return <-errCh
	case err := <-errCh:
		return err
	}
}
