package routes

import (
	"net/http"
	"time"

	"github.com/eriicafes/tmpl"
	"github.com/eriicafes/tmplist/db"
	"github.com/eriicafes/tmplist/internal"
	"github.com/eriicafes/tmplist/internal/session"
	v "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	requestUser    internal.ContextValue[db.User]    = "user"
	requestSession internal.ContextValue[db.Session] = "session"
)

type Context struct {
	tmpl.Templates
	DB   db.DB
	Auth *session.Auth[db.Session, db.User]
	Prod bool
}

func (c Context) Mount(mux internal.ServeMux) {
	c.Classic(internal.Prefix(mux, "/classic"))
	c.Enhanced(internal.Prefix(mux, "/enhanced"))
	c.Api(internal.Prefix(mux, "/api"))
	c.Spa(internal.Prefix(mux, "/spa"))

	// redirect if mode is set in cookies
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		if modeCookie, err := r.Cookie("mode"); err == nil {
			switch modeCookie.Value {
			case "classic", "enhanced", "spa":
				http.Redirect(w, r, "/"+modeCookie.Value, http.StatusFound)
				return
			}
		}
		if err := c.Render(w, tmpl.Tmpl("index", nil)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// set mode in cookies
	mux.HandleFunc("POST /{$}", func(w http.ResponseWriter, r *http.Request) {
		mode := r.PostFormValue("mode")
		err := v.Validate(mode, v.In("none", "classic", "enhanced", "spa"))
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "mode",
			Value:   mode,
			Expires: time.Now().Add(time.Hour * 24 * 30), // 30 days
			Path:    "/",
		})
		if mode == "none" {
			mode = ""
		}
		http.Redirect(w, r, "/"+mode, http.StatusFound)
	})
}
