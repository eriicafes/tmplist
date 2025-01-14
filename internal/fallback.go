package internal

import "net/http"

// Fallback returns a Mux which executes errorHandler when the handler returns a non-nil error.
func Fallback(mux Mux, errorHandler ErrorHandlerFunc) Mux {
	return &fallbackMux{mux, errorHandler}
}

type fallbackMux struct {
	Mux
	errorHandler func(http.ResponseWriter, *http.Request, error)
}

func (fm *fallbackMux) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	fm.errorHandler(w, r, err)
}

func (fm *fallbackMux) Route(pattern string, handler func(http.ResponseWriter, *http.Request) error) {
	fm.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			fm.HandleError(w, r, err)
		}
	})
}
