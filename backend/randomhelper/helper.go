/*
Package randomhelper that has helper functions that doesn't belong to a particular destination!
*/
package randomhelper

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

func CheckIfAllEnvValid(variables ...string) {
	for i, pass := range variables {
		if pass == "" {
			log.Fatalf("could not load variable at index: %d", i)
		}
	}
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Vary", "Origin")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

const DefaultTokenLength = 32

func Generate(length int) string {
	bytes := make([]byte, length)
	//:NOTE: this almost never errors, it reads from /dev/urandom
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}

func GetMessageID() string {
	id, _ := uuid.NewV7()
	return id.String()
}
