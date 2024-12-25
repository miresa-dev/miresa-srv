package main

import (
	"net/http"
	"strconv"
	"log"

	"github.com/miresa-dev/miresa-srv/internal/api"
	"github.com/miresa-dev/miresa-srv/internal/conf"
	"github.com/miresa-dev/miresa-srv/internal/db"
	"github.com/miresa-dev/miresa-srv/internal/middleware"
	"github.com/miresa-dev/miresa-srv/internal/web"

	"github.com/go-chi/chi/v5"
)

func main() {
	// Configuration for the server.
	config, err := conf.LoadConf()
	if err != nil {
		panic(err)
	}

	// Initializing the database.
	if err = db.Init(config.DatabaseURL); err != nil {
		panic(err)
	}
	defer db.Close()

	// Initializing the routers.
	r := chi.NewRouter()
	a := chi.NewRouter()

	r.Use(middleware.Log)

	r.Get("/", web.Home)

	a.Get("/v", api.Version)

	a.Post("/u",       api.CreateUser)
	a.Get("/u/{id}",   api.GetUser)
	//a.Patch("/u/{id}", api.UpdateUser)

	r.Mount("/api/v0", a)

	// Starting the server.
	log.Printf("Listening on :%d\n", config.Port)
	http.ListenAndServe(":"+strconv.Itoa(config.Port), r)
}
