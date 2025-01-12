package routes

import (
	"fmt"
	"net/http"

	"github.com/eriicafes/tmplist/request"
)

func (c Context) authMiddleware() Middleware {
	return func(handler Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			rc := c.Session.Authenticate(w, r)
			_, user := request.User.FromContext(rc.Context())
			_, session := request.Session.FromContext(rc.Context())

			if !user || !session {
				http.Redirect(w, r, "/classic/login", http.StatusFound)
				return nil
			}
			return handler(w, rc)
		}
	}
}

func (c Context) guestMiddleware() Middleware {
	return func(handler Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			rc := c.Session.Authenticate(w, r)
			_, user := request.User.FromContext(rc.Context())
			_, session := request.Session.FromContext(rc.Context())

			fmt.Println(user, session)

			if user || session {
				http.Redirect(w, r, "/classic/", http.StatusFound)
				return nil
			}
			return handler(w, rc)
		}
	}
}
