package api

import (
	"encoding/json"
	"log"
	"strconv"
	"net/http"
	"runtime"
	"time"

	"github.com/miresa-dev/miresa-srv/internal/conf"
)

// Version writes version information to the HTTP response and sends it back.
func Version(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{}

	v := conf.Config.VersionEndpoint

	if v.OS {
		data["os"] = runtime.GOOS
	}
	if v.Arch {
		data["arch"] = runtime.GOARCH
	}
	if v.GoVersion {
		data["go"] = runtime.Version()
	}
	if v.GoroutineCount {
		data["goroutines"] = strconv.Itoa(runtime.NumGoroutine())
	}
	if v.ServerTime {
		data["time"] = time.Now().String()
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// It's a map[string]string. If we can't marshal it, something
		// is wrong.
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to marshal JSON: %v\n", err)
	}
}
