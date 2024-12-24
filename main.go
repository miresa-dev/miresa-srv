package main

import (
	"net/http"

	"github.com/miresa-dev/miresa-srv/internal/api"
	"github.com/miresa-dev/miresa-srv/internal/middleware"
	"github.com/miresa-dev/miresa-srv/internal/web"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	a := chi.NewRouter()

	r.Use(middleware.Log)

	r.Get("/", web.Home)

	a.Get("/v", api.Version)

	r.Mount("/api/v0", a)
	
	http.ListenAndServe(":8000", r)
}
