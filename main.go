package main

import (
	"context"
	"errors"
	"fmt"
	"http-server/internal"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	app := internal.App{
		Addr:              "localhost:8080",
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       15 * time.Second,
	}

	server := app.GetServer()

	fmt.Printf("Starting server at http://%v...\n", app.Addr)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.ListenAndServe()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	watcher := internal.Watcher{
		Targets: []internal.Target{
			{
				ID:   "1",
				Name: "GitHub",
				URL:  "https://github.com",
			},
			{
				ID:   "2",
				Name: "Youtube",
				URL:  "https://www.youtube.com",
			},
			{
				ID:   "3",
				Name: "Google",
				URL:  "https://www.google.com",
			},
		},
	}

	go func() {
		watcher.Watch(ctx)
	}()

	select {
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(os.Stderr, "server error: %v\n", err)
			os.Exit(1)
		}
	case <-ctx.Done():
		fmt.Println("Shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Fprintf(os.Stderr, "graceful shutdown failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Server stopped gracefully")
}
