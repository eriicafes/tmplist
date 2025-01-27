package routes

import (
	"net/http"

	"github.com/eriicafes/tmpl"
	"github.com/eriicafes/tmplist/internal"
)

func (c Context) Spa(mux internal.Mux) {
	mux.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		c.Render(w, tmpl.Tmpl("spa/index"))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c.Render(w, tmpl.Tmpl("spa/index"))
	})
}
