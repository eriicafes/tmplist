package internal

import (
	"net/http"
	"slices"
)

// Apply returns a Mux where middlewares will be applied to registered routes.
// Each middleware wraps the next and eventually the handler.
func Apply(mux Mux, middlewares ...Middleware) Mux {
	return &applyMux{mux, middlewares}
}

type applyMux struct {
	Mux
	middlewares []Middleware
}

func (am *applyMux) Handle(pattern string, handler http.Handler) {
	// Each middleware wraps the next and finally the handler.
	for _, mh := range slices.Backward(am.middlewares) {
		handler = mh(handler)
	}
	am.Mux.Handle(pattern, handler)
}

func (am *applyMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	am.Handle(pattern, http.HandlerFunc(handler))
}

func (am *applyMux) Route(pattern string, handler func(http.ResponseWriter, *http.Request) error) {
	am.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			ErrorHandler(am.Mux, err)(w, r)
		}
	})
}
