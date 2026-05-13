package internal

import (
	"log"
	"net/http"
	"time"
)

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Printf(
				"%s %s %v", 
				r.Method, r.URL.Path, time.Since(start),
			)
	})
}

func WrapAllMiddleware(next http.Handler) http.Handler {
	return logMiddleware(next)
}
