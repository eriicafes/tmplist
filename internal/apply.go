package internal

import (
	"net/http"
	"slices"
)

// Apply returns a Mux which applies middlewares to the handler.
// Apply with an empty middlewares list can be used to convert a *http.ServeMux to a Mux.
func Apply(mux ServeMux, middlewares ...Middleware) Mux {
	return &applyMux{mux, middlewares}
}

type applyMux struct {
	ServeMux
	middlewares []Middleware
}

func (am *applyMux) Handle(pattern string, handler http.Handler) {
	// Each middleware wraps the next and finally the handler.
	for _, mh := range slices.Backward(am.middlewares) {
		handler = mh(handler)
	}
	am.ServeMux.Handle(pattern, handler)
}

func (am *applyMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	am.Handle(pattern, http.HandlerFunc(handler))
}

func (am *applyMux) Route(pattern string, handler func(http.ResponseWriter, *http.Request) error) {
	am.Handle(pattern, withErrorHandler(am, handler))
}
