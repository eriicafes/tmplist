package main

import (
	"net/http"

	"github.com/eriicafes/tmplist/db"
	"github.com/eriicafes/tmplist/routes"
	"github.com/eriicafes/tmplist/services"
)

func main() {
	config := getConfig()
	templates, vite := setupTemplates(!config.Prod)
	database := db.Connect(config.DbURL)
	auth := &services.AuthService{DB: database}
	session := &services.SessionService{Prod: config.Prod, Auth: auth}
	rc := routes.Context{
		Templates: templates,
		DB:        database,
		Auth:      auth,
		Session:   session,
	}

	mux := http.NewServeMux()
	// mount routes under prefixes
	rc.MountClassic(routes.Prefix(mux, "/classic"))
	rc.MountEnhanced(routes.Prefix(mux, "/enhanced"))
	rc.MountSpa(routes.Prefix(mux, "/spa"))
	rc.MountApi(routes.Prefix(mux, "/api"))

	// serve vite static assets
	mux.Handle("/", vite.ServePublic(http.NotFoundHandler()))

	http.ListenAndServe(config.ListenAddr(), mux)
}
