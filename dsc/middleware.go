package dsc

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func protocolVersion(inner http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if ver := r.Header.Get("ProtocolVersion"); ver != "2.0" {
			w.WriteHeader(http.StatusNotImplemented)

			// TODO: real error
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintf(w, "protocol version %q not supported", ver)
			return
		}

		// Responses are protocol version 2.0 as well.
		w.Header()["ProtocolVersion"] = []string{"2.0"}

		// Pass to the underlying HTTP handler
		inner.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func checkRegistration(keys []string) func(http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Get needed headers; if they aren't present, it's a
			// bad request.
			msDate := r.Header.Get("x-ms-date")
			if msDate == "" {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `missing "x-ms-date" header`)
				return
			}
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `missing "Authorization" header`)
				return
			}

			authHeader = strings.TrimPrefix(authHeader, "Shared ")

			// Read the entire body of the request; we need to do
			// this in order to verify the signature.
			body, err := ioutil.ReadAll(r.Body)
			r.Body.Close()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "error reading request body: %s", err)
				return
			}

			// Replace the body of the request so it can be read again.
			r.Body = ioutil.NopCloser(bytes.NewReader(body))

			// Hash the request body
			h := sha256.New()
			h.Write(body)
			bodyHash := h.Sum(nil)

			// For each registration key, attempt to validate the signature
			match := false
			for _, rkey := range keys {
				// Verify against the `Authorization` header
				expectedAuth := calculateDSCSignature(rkey, msDate, bodyHash)
				if hmac.Equal([]byte(expectedAuth), []byte(authHeader)) {
					match = true
					break
				}
			}

			// If we didn't match, return a "Unauthorized" response.
			if !match {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, `bad signature on request`)
				return
			}

			// Otherwise, we're good; call our underlying handler.
			inner.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func calculateDSCSignature(key, date string, bodyHash []byte) string {
	bodyEnc := base64.StdEncoding.EncodeToString(bodyHash)

	// The string to sign is the body of the request
	// concatenated with the `x-ms-date` header.
	stringToSign := bodyEnc + "\n" + date

	// Create HMAC signature using this key.
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(stringToSign))
	expectedAuth := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return expectedAuth
}
