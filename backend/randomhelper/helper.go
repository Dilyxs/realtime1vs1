/*
Package randomhelper that has helper functions that doesn't belong to a particular destination!
*/
package randomhelper

import (
	"log"
	"net/http"
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
