package api

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"time"
)

// Version writes version information to the HTTP response and sends it back.
func Version(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{}

	data["os"] = runtime.GOOS
	data["arch"] = runtime.GOARCH
	data["go"] = runtime.Version()
	data["goroutines"] = "not yet implemented"
	data["time"] = time.Now().String()

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// It's a map[string]string. If we can't marshal it, something
		// is wrong.
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to marshal JSON: %v\n", err)
	}
}
