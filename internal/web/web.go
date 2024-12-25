package web

import (
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if err := w.Write([]byte("<!DOCTYPE html>This is in progress.")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
