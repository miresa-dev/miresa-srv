package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/miresa-dev/miresa-srv/internal/api"
	"github.com/miresa-dev/miresa-srv/internal/conf"
	"github.com/miresa-dev/miresa-srv/internal/db"
	"github.com/miresa-dev/miresa-srv/internal/middleware"
	"github.com/miresa-dev/miresa-srv/internal/web"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
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
	r.Use(chiMiddleware.StripSlashes)

	r.Get("/", web.Home)
	//r.Get("/u/{id}", web.User)
	//r.Get("/i/{id}", web.Item)

	a.Get("/v", api.Version)

	a.Get("/init", api.CaptchaAndSID)
	a.Post("/login", api.Login)
	a.Post("/logout", api.Logout)
	a.Post("/u", api.CreateUser)
	a.Get("/u", api.GetAllUsers)
	a.Get("/u/{id}", api.GetUser)
	a.Patch("/u/{id}", api.UpdateUser)

	a.Get("/i", api.GetItems)
	a.Get("/i/{id}", api.GetItem)
	a.Post("/i", api.CreateItem)
	a.Patch("/i/{id}", api.UpdateItem)

	r.Mount("/api/v0", a)

	// Starting the server.
	log.Printf("Listening on :%d\n", config.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), r))
}
