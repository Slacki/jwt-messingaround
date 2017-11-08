package main

import (
	"net/http"
)

func authenticatedMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func auth(w http.ResponseWriter, r *http.Request) {
}
