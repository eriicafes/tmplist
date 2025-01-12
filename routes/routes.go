package routes

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/eriicafes/tmpl"
	"github.com/eriicafes/tmplist/db"
	"github.com/eriicafes/tmplist/httperrors"
	"github.com/eriicafes/tmplist/services"
)

type Context struct {
	tmpl.Templates
	DB      db.DB
	Auth    *services.AuthService
	Session *services.SessionService
}

type Mux interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

// Prefix returns a wrapped Mux where routes will be mounted under the given prefix.
func Prefix(mux Mux, prefix string) Mux {
	return &prefixMux{mux, prefix}
}

type prefixMux struct {
	mux    Mux
	prefix string
}

func (pm *prefixMux) prefixPattern(pattern string) string {
	if strings.Contains(pattern, " ") {
		return strings.Replace(pattern, " ", " "+pm.prefix, 1)
	}
	return pm.prefix + pattern
}

func (pm *prefixMux) Handle(pattern string, handler http.Handler) {
	pm.mux.Handle(pm.prefixPattern(pattern), handler)
}

func (pm *prefixMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	pm.mux.HandleFunc(pm.prefixPattern(pattern), handler)
}

// Route returns a Mux where routes can be grouped with middlewares.
// Each middleware wraps the next and eventually the handler.
// When a non-nil error is returned an error response is returned.
// Route exposes the On(pattern, handler) method for route handlers that returns an error.
func Route(mux Mux, middlewares ...Middleware) route { return route{mux, middlewares} }

type route struct {
	mux         Mux
	middlewares []Middleware
}

func (r *route) Handle(pattern string, handler http.Handler) {
	r.mux.Handle(pattern, combine(ToHandler(handler), r.middlewares...))
}

func (r *route) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	r.mux.HandleFunc(pattern, combine(ToHandler(http.HandlerFunc(handler)), r.middlewares...))
}

func (r *route) On(pattern string, handler func(http.ResponseWriter, *http.Request) error) {
	r.mux.HandleFunc(pattern, combine(handler, r.middlewares...))
}

// Handler is a handler func that returns an error
type Handler func(http.ResponseWriter, *http.Request) error

// Middleware is a handler func that wraps another handler.
type Middleware func(handler Handler) Handler

// combine chains a group of middlewares with the route handler.
// Each middleware wraps the next and eventually the handler.
// When a non-nil error is returned an error response is returned.
func combine(handler Handler, middlewares ...Middleware) http.HandlerFunc {
	for _, mh := range slices.Backward(middlewares) {
		handler = mh(handler)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err == nil {
			return
		}
		var herr httperrors.HTTPError
		if !errors.As(err, &herr) {
			log.Println("Unexpected error:", err)
			herr = httperrors.New("Something went wrong", http.StatusInternalServerError)
		}
		statusCode, msg, details := herr.HTTPError()

		// return json error response
		if r.Header.Get("accept") == "application/json" {
			w.WriteHeader(statusCode)
			jsonResp, _ := json.Marshal(map[string]any{
				"message": msg,
				"errors":  details,
			})
			w.Write(jsonResp)
			return
		}

		// return fallback error response
		http.Error(w, msg, statusCode)
	})
}

func ToHandler(handler http.Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		handler.ServeHTTP(w, r)
		return nil
	}
}
