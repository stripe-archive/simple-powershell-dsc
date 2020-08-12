package middleware

import (
	"fmt"
	"net/http"
)

func MaxBodySize(size int64) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > size {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "request body too large")
				return
			}

			r.Body = http.MaxBytesReader(w, r.Body, size)
			h.ServeHTTP(w, r)
		})
	}
}
