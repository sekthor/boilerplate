package boilerplate

import (
	"net/http"
	"strings"
)

func corsHandler(origins []string, methods []string) func(http.Handler) http.Handler {

	allowedOrigins := "*"
	if len(origins) > 0 {
		allowedOrigins = strings.Join(origins, ", ")
	}

	allowedMethods := "GET, PUT, POST, DELETE, HEAD, OPTIONS"
	if len(methods) > 0 {
		allowedMethods = strings.Join(methods, ", ")
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
			w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
