package internal

import (
	"net/http"
)

// Mux is a specialized mux that can register routes with handlers that return an error.
type Mux interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	Route(pattern string, handler func(http.ResponseWriter, *http.Request) error)
}

// NewMux wraps a http.ServeMux into a Mux with a default error handler.
func NewMux(mux *http.ServeMux) Mux {
	return &defaultMux{mux}
}

type defaultMux struct{ *http.ServeMux }

func (dm *defaultMux) Route(pattern string, handler func(http.ResponseWriter, *http.Request) error) {
	dm.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
