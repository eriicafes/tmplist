package httperrors

import (
	"errors"
	"net/http"
)

type HTTPError interface {
	error
	Unwrap() error
	HTTPError() (statusCode int, message string, details Details)
}

type Details = map[string]string

// New returns a new http error with the provided status code.
func New(err string, statusCode int) HTTPError {
	return &httpError{errors.New(err), statusCode, nil}
}

// WithStatus wraps an error into an http error with the provided status code.
// If err is an http error it's underlying error and details will be preserved.
func WithStatus(err error, statusCode int) HTTPError {
	if herr, ok := err.(*httpError); ok {
		return &httpError{herr.error, statusCode, herr.details}
	}
	return &httpError{err, statusCode, nil}
}

// WithDetails wraps an error into an http error with the provided details.
// If err is an http error it's underlying error and statusCode will be preserved.
// If err is not an http error the status code will be set to http.StatusBadRequest.
func WithDetails(err error, details Details) HTTPError {
	if herr, ok := err.(*httpError); ok {
		return &httpError{herr.error, herr.statusCode, details}
	}
	return &httpError{err, http.StatusBadRequest, details}
}

type httpError struct {
	error
	statusCode int
	details    Details
}

func (e *httpError) Unwrap() error {
	return e.error
}

func (e *httpError) HTTPError() (int, string, Details) {
	statusCode := e.statusCode
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}
	return statusCode, e.error.Error(), e.details
}

// Opaque returns an error which hides the underlying error from the http response.
// The returned error returns the provided error message but unrwaps to the underlying error.
func Opaque(msg string, cause error) error {
	return &opaqueError{msg, cause}
}

type opaqueError struct {
	msg   string
	cause error
}

func (e *opaqueError) Unwrap() error {
	return e.cause
}

func (e *opaqueError) Error() string {
	return e.msg
}
