package routes

import (
	"net/http"

	"github.com/eriicafes/tmplist/internal"
)

func (c Context) authMiddleware() internal.Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, user, ok := c.Auth.Authenticate(w, r)
			if !ok {
				http.Redirect(w, r, "/classic/login", http.StatusFound)
				return
			}
			ctx := requestSession.Set(r.Context(), session)
			ctx = requestUser.Set(ctx, user)
			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (c Context) guestMiddleware() internal.Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _, ok := c.Auth.Authenticate(w, r)
			if ok {
				http.Redirect(w, r, "/classic/", http.StatusFound)
				return
			}
			handler.ServeHTTP(w, r)
		})
	}
}
