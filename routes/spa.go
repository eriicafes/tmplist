package routes

import (
	"net/http"

	"github.com/eriicafes/tmpl"
	"github.com/eriicafes/tmplist/internal"
)

func (c Context) Spa(mux internal.Mux) {
	mux.HandleFunc("GET ", func(w http.ResponseWriter, r *http.Request) {
		c.Render(w, tmpl.Tmpl("spa/index", nil))
	})
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		c.Render(w, tmpl.Tmpl("spa/index", nil))
	})
}
