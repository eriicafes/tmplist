package internal

import (
	"net/http"
	"slices"
)

// Use returns a Mux which applies middlewares to the handler.
// Use with an empty middlewares list can be used to convert a *http.ServeMux to a Mux.
func Use(mux ServeMux, middlewares ...Middleware) Mux {
	return &useMux{mux, middlewares}
}

type useMux struct {
	ServeMux
	middlewares []Middleware
}

func (am *useMux) Handle(pattern string, handler http.Handler) {
	for _, mh := range slices.Backward(am.middlewares) {
		handler = mh(handler)
	}
	am.ServeMux.Handle(pattern, handler)
}

func (am *useMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	am.Handle(pattern, http.HandlerFunc(handler))
}

func (fm *useMux) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	ErrorHandler(fm.ServeMux, err)(w, r)
}

func (am *useMux) Route(pattern string, handlers ...func(http.ResponseWriter, *http.Request) error) {
	am.Handle(pattern, routeErrorHandler(am, handlers))
}
