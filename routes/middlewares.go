package routes

import (
	"net/http"
	"net/url"

	"github.com/eriicafes/tmplist/internal"
)

func (c Context) authMiddleware() internal.Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// prevent non safe cross-subdomain requests
			if !c.allowOriginForNonSafeRequests(r) {
				http.Redirect(w, r, "/classic/login", http.StatusFound)
				return
			}

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

// allowOriginForNonSafeRequests returns a boolean indicating if the request
// should be allowed for non-safe http methods based on its origin header.
func (c Context) allowOriginForNonSafeRequests(r *http.Request) bool {
	switch r.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
		// skip for safe methods
		return true
	default:
		origin, err := url.Parse(r.Header.Get("Origin"))
		if err != nil {
			return false
		}
		if origin.Host != r.Host {
			return false
		}
		return true
	}
}
