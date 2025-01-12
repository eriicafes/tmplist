package services

import (
	"net/http"
	"time"

	"github.com/eriicafes/tmplist/request"
)

const sessionCookie = "auth_session"
const sessionDuration = time.Hour * 24 * 30 // expire in 30 days

type SessionService struct {
	Prod bool
	Auth *AuthService
}

// Authenticate gets the session token from the request cookies and populates the
func (s *SessionService) Authenticate(w http.ResponseWriter, r *http.Request) *http.Request {
	// get session token from cookies
	token, ok := s.GetCookie(r)
	if !ok {
		return r
	}
	// validate session token
	session, user, err := s.Auth.ValidateSessionToken(token)
	if err != nil {
		// delete invalid session
		if _, err := r.Cookie(sessionCookie); err != nil {
			s.DeleteCookie(w)
		}
		return r
	}
	// set cookie to potentially refresh session
	s.SetCookie(w, token, session.ExpiresAt)

	// set user and session in request context
	ctx := r.Context()
	ctx = request.User.SetContext(ctx, user)
	ctx = request.Session.SetContext(ctx, session)
	return r.WithContext(ctx)
}

func (s *SessionService) GetCookie(r *http.Request) (string, bool) {
	token, err := r.Cookie(sessionCookie)
	if err != nil {
		return "", false
	}
	return token.Value, err == nil
}

func (s *SessionService) SetCookie(w http.ResponseWriter, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		HttpOnly: true,
		Secure:   s.Prod,
		SameSite: http.SameSiteLaxMode,
		Value:    token,
		Expires:  expiresAt.UTC(),
	})
}

func (s *SessionService) DeleteCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		HttpOnly: true,
		Secure:   s.Prod,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}
