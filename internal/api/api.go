package api

import (
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/miresa-dev/miresa-srv/internal/conf"
)

// Version writes version information to the HTTP response and sends it back.
func Version(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{}

	data["os"] = runtime.GOOS
	data["arch"] = runtime.GOARCH
	data["go"] = runtime.Version()
	data["goroutines"] = "not yet implemented"

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// It's a map[string]string. If we can't marshal it, something
		// is wrong.
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Conf(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(conf.Config); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
