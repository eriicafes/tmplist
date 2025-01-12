package routes

import (
	"net/http"

	"github.com/eriicafes/tmpl"
)

func (c Context) MountEnhanced(mux Mux) {
	mux.HandleFunc("/{$}", combine(func(w http.ResponseWriter, r *http.Request) error {
		tr := c.Renderer()
		return tr.Render(w, tmpl.Tmpl("enhanced/pages/index"))
	}))
	mux.HandleFunc("/{topicId}", combine(func(w http.ResponseWriter, r *http.Request) error {
		tr := c.Renderer()
		return tr.Render(w, tmpl.Tmpl("enhanced/pages/topic"))
	}))
	mux.HandleFunc("/login", combine(func(w http.ResponseWriter, r *http.Request) error {
		tr := c.Renderer()
		return tr.Render(w, tmpl.Tmpl("enhanced/pages/login"))
	}))
	mux.HandleFunc("/register", combine(func(w http.ResponseWriter, r *http.Request) error {
		tr := c.Renderer()
		return tr.Render(w, tmpl.Tmpl("enhanced/pages/register"))
	}))
	mux.HandleFunc("/", combine(func(w http.ResponseWriter, r *http.Request) error {
		tr := c.Renderer()
		return tr.Render(w, tmpl.Tmpl("enhanced/pages/404"))
	}))
}
