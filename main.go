package main

import (
	"net/http"

	"github.com/eriicafes/tmplist/db"
	"github.com/eriicafes/tmplist/internal"
	"github.com/eriicafes/tmplist/internal/session"
	"github.com/eriicafes/tmplist/routes"
)

func main() {
	config := getConfig()
	templates, vite := setupTemplates(!config.Prod)
	database := db.Connect(config.DbURL)
	auth := session.NewAuth(
		db.SessionStorage{DB: database},
		session.SessionOptions{Secure: config.Prod, Path: "/"},
	)
	rc := routes.Context{
		Templates: templates,
		DB:        database,
		Auth:      auth,
		Prod:      config.Prod,
	}
	mux := http.NewServeMux()
	rc.Mount(mux)
	// serve vite static assets
	mux.Handle("/", vite.ServePublic(http.NotFoundHandler()))

	h := internal.RewriteTrailingSlash(mux)
	http.ListenAndServe(config.ListenAddr(), h)
}
