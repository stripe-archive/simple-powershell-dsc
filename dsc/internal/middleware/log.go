package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zenazn/goji/web/mutil"
)

// LogrusLogger is a middleware that will log each request recieved, along with
// some useful information, to the given logger.
func LogrusLogger(logger logrus.FieldLogger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// If we're running using apex/gateway, we won't have a
			// RequestURI, so just fall back to stringifying the
			// URL.
			requestURI := r.RequestURI
			if requestURI == "" {
				requestURI = r.URL.String()
			}

			entry := logger.WithFields(logrus.Fields{
				"request": requestURI,
				"method":  r.Method,
				"remote":  r.RemoteAddr,
			})

			// TODO: request ID?

			// Wrap the writer so we can track data information.
			neww := mutil.WrapWriter(w)

			// Dispatch to the underlying handler.
			entry.Info("started handling request")
			h.ServeHTTP(neww, r)
			duration := time.Since(start)

			// Log final information.
			entry.WithFields(logrus.Fields{
				"bytes_written": neww.BytesWritten(),
				"status":        neww.Status(),
				"duration":      float64(duration.Nanoseconds()) / float64(1000),
			}).Info("completed handling request")
		}
		return http.HandlerFunc(fn)
	}
}
