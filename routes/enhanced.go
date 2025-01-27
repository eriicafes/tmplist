package routes

import (
	"net/http"

	"github.com/eriicafes/tmpl"
	"github.com/eriicafes/tmplist/internal"
)

func (c Context) Enhanced(mux internal.Mux) {
	mux.Route("", func(w http.ResponseWriter, r *http.Request) error {
		return c.Render(w, tmpl.Tmpl("enhanced/pages/index"))
	})
	mux.Route("/{topicId}", func(w http.ResponseWriter, r *http.Request) error {
		return c.Render(w, tmpl.Tmpl("enhanced/pages/topic"))
	})
	mux.Route("/login", func(w http.ResponseWriter, r *http.Request) error {
		return c.Render(w, tmpl.Tmpl("enhanced/pages/login"))
	})
	mux.Route("/register", func(w http.ResponseWriter, r *http.Request) error {
		return c.Render(w, tmpl.Tmpl("enhanced/pages/register"))
	})
	mux.Route("/", func(w http.ResponseWriter, r *http.Request) error {
		return c.Render(w, tmpl.Tmpl("enhanced/pages/error"))
	})
}
