package internal

import (
	"net/http"
)

// HandlerFunc is a handler func that returns an error
type HandlerFunc func(http.ResponseWriter, *http.Request) error

type ErrorHandlerFunc func(http.ResponseWriter, *http.Request, error)

// Middleware is a function that accepts a handler and returns another handler.
// Middlewares can run before or after the handler and they decide whether or not to call the handler.
type Middleware func(handler http.Handler) http.Handler

type HandleError interface {
	HandleError(http.ResponseWriter, *http.Request, error)
}

type withRequest struct{ *http.Request }

func (r withRequest) Error() string { return "" }

func WithRequest(r *http.Request) error { return withRequest{r} }

func routeErrorHandler(mux ServeMux, handlers []func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		for _, handler := range handlers {
			err = handler(w, r)
			if wr, ok := err.(withRequest); ok {
				r = wr.Request
				continue
			}
			if err != nil {
				ErrorHandler(mux, err)(w, r)
				break
			}
		}
	}
}

// ErrorHandler returns the error handler func for mux.
// If mux does not implement the HandleError interface the returned handler func will write a default error response.
func ErrorHandler(mux ServeMux, err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if h, ok := mux.(HandleError); ok {
			h.HandleError(w, r, err)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
