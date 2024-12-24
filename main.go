package main

import (
	"net/http"
	"strconv"
	"fmt"

	"github.com/miresa-dev/miresa-srv/internal/api"
	"github.com/miresa-dev/miresa-srv/internal/conf"
	"github.com/miresa-dev/miresa-srv/internal/middleware"
	"github.com/miresa-dev/miresa-srv/internal/web"

	"github.com/go-chi/chi/v5"
)

func main() {
	config, err := conf.LoadConf()
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	a := chi.NewRouter()

	r.Use(middleware.Log)

	r.Get("/", web.Home)

	a.Get("/v", api.Version)
	a.Get("/c", api.Conf)

	r.Mount("/api/v0", a)
	
	http.ListenAndServe(":"+strconv.Itoa(config.Port), r)
}
