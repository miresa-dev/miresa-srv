package middleware

import (
	"log"
	"net/http"
)

// Log logs HTTP request methods and paths.
func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: get the HTTP status that was given to `next` and print
		// it as well.
		next.ServeHTTP(w, r)

		log.Printf("%s %s", r.Method, r.URL.Path)
	})
}
