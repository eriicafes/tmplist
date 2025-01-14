package internal

import (
	"net/http"
	"strings"
)

// Prefix returns a Mux where routes will be registered under the given prefix.
func Prefix(mux Mux, prefix string) Mux {
	return &prefixMux{mux, prefix}
}

type prefixMux struct {
	Mux
	prefix string
}

func (pm *prefixMux) prefixPattern(pattern string) string {
	if strings.Contains(pattern, " ") {
		return strings.Replace(pattern, " ", " "+pm.prefix, 1)
	}
	return pm.prefix + pattern
}

func (pm *prefixMux) Handle(pattern string, handler http.Handler) {
	pm.Mux.Handle(pm.prefixPattern(pattern), handler)
}

func (pm *prefixMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	pm.Mux.HandleFunc(pm.prefixPattern(pattern), handler)
}

func (pm *prefixMux) Route(pattern string, handler func(http.ResponseWriter, *http.Request) error) {
	pm.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			ErrorHandler(pm.Mux, err)(w, r)
		}
	})
}
