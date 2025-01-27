package routes

import (
	"net/http"
	"time"

	"github.com/eriicafes/tmpl"
	"github.com/eriicafes/tmplist/internal"
	"github.com/eriicafes/tmplist/schemas"
)

type Tpl struct {
	tmpl.Templates
}

func (c Context) Settings(mux internal.ServeMux) {
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		if modeCookie, err := r.Cookie("mode"); err == nil {
			switch modeCookie.Value {
			case "classic", "enhanced", "spa":
				http.Redirect(w, r, "/"+modeCookie.Value, http.StatusFound)
				return
			}
		}
		if err := c.Render(w, tmpl.Tmpl("settings/index")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("POST /{$}", func(w http.ResponseWriter, r *http.Request) {
		form := schemas.SettingsData{
			Mode:  r.PostFormValue("mode"),
			Delay: r.PostFormValue("delay"),
		}
		if err := form.Validate(); err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "mode",
			Value:   form.Mode,
			Expires: time.Now().Add(time.Hour * 24 * 30), // 30 days
			Path:    "/",
		})
		http.SetCookie(w, &http.Cookie{
			Name:    "delay",
			Value:   form.Delay,
			Expires: time.Now().Add(time.Hour * 24 * 30), // 30 days
			Path:    "/",
		})
		if form.Mode == "none" {
			form.Mode = ""
		}
		http.Redirect(w, r, "/"+form.Mode, http.StatusFound)
	})
}
