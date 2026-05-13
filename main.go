package main

import (
	"http-server/internal"
	"fmt"
	"time"
)

func main() {
	app := internal.App{
		Addr: "localhost:8080",
		ReadTimeout: 5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 15 * time.Second,
	}

	server := app.GetServer()

	fmt.Printf("Starting server at http://%v...\n", app.Addr)
	server.ListenAndServe()
}
