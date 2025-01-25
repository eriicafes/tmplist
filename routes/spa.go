package routes

import (
	"net/http"

	"github.com/eriicafes/tmpl"
	"github.com/eriicafes/tmplist/internal"
)

func (c Context) Spa(mux internal.Mux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tr := c.Renderer()
		tr.Render(w, tmpl.Tmpl("spa/index"))
	})
}
