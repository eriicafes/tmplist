package internal

import (
	"net/http"
)

// ServeMux is the generic interface for route registering mux.
type ServeMux interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

// Mux is a specialized mux that has a Route method.
// The Route method registers routes with handlers that return an error
type Mux interface {
	ServeMux
	Route(pattern string, handler func(http.ResponseWriter, *http.Request) error)
}

func New() Mux {
	return Use(http.NewServeMux())
}

func MuxHandler(mux Mux) http.Handler {
	if h, ok := mux.(http.Handler); ok {
		return h
	}
	return nil
}
