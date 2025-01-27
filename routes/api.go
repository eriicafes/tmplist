package routes

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/eriicafes/tmplist/internal"
	"github.com/eriicafes/tmplist/internal/httperrors"
)

func (c Context) Api(mux internal.Mux) {
	mux = internal.Fallback(mux, c.ApiErrorHandler())

	mux.Handle("", internal.ErrorHandler(mux, httperrors.New("api routes", http.StatusOK)))
	mux.Handle("/", internal.ErrorHandler(mux, httperrors.New("route not found", http.StatusNotFound)))
}

func (c Context) ApiErrorHandler() internal.ErrorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		var herr httperrors.HTTPError
		if !errors.As(err, &herr) {
			log.Println("Unexpected error:", err)
			herr = httperrors.New("Something went wrong", http.StatusInternalServerError)
		}
		statusCode, msg, details := herr.HTTPError()

		// return json error response
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(statusCode)
		jsonResp, _ := json.Marshal(map[string]any{
			"message": msg,
			"errors":  details,
		})
		w.Write(jsonResp)
	}
}
