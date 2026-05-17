package internal

import (
	"fmt"
	"net/http"
	"time"
)

type App struct {
	Addr              string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello, internet!</h1>")
}

func (a *App) GetServer() *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlerFunc)

	wMux := WrapAllMiddleware(mux)

	srv := &http.Server{
		Addr:              a.Addr,
		Handler:           wMux,
		ReadTimeout:       a.ReadTimeout,
		ReadHeaderTimeout: a.ReadHeaderTimeout,
		WriteTimeout:      a.WriteTimeout,
		IdleTimeout:       a.IdleTimeout,
	}

	return srv
}
