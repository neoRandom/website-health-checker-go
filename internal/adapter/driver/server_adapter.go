package driver

import (
	"fmt"
	"http-server/internal/adapter/driver/middleware"
	"log"
	"net/http"
)

type ServerAdapter struct {
	Addr string
}

func (s ServerAdapter) Init() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})

	wMux := middleware.ChainMiddleware(
		mux,
		middleware.LogMiddleware,
		middleware.CorsMiddleware,
	)

	srv := http.Server{
		Addr: s.Addr,
		Handler: wMux,
	}

	log.Printf("Server starting at http://%v...", s.Addr)
	srv.ListenAndServe()
}
