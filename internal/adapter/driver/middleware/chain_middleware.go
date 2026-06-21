package middleware

import "net/http"

func ChainMiddleware(
	next http.Handler,
	middlewares ...func(http.Handler) http.Handler,
) http.Handler {
	for i := range middlewares {
		next = middlewares[i](next)
	}

	return next
}
