package session

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var ErrSessionExpired = errors.New("session expired")

// Auth provides the useful methods for working with user sessions.
type Auth[S SessionObject, U any] struct {
	storage Storage[S, U]
	options Options
}

type Storage[S SessionObject, U any] interface {
	GetSessionAndUser(id string) (S, U, error)
	CreateSession(id string, user U, expiresAt time.Time) (S, error)
	UpdateSession(session S, expiresAt time.Time) (S, error)
	DeleteSession(id string) error
}

type SessionObject interface{ GetExpirestAt() time.Time }

type Options struct {
	// Cookie is the session cookie name. Default is "auth_session".
	Cookie string

	// Duration is the session duration. Default is 30 days.
	Duration time.Duration

	// Secure sends cookies over HTTPS only. Default is false.
	Secure bool
}

// New returns a new auth provider.
func NewAuth[S SessionObject, U any](storage Storage[S, U], options Options) *Auth[S, U] {
	if options.Cookie == "" {
		options.Cookie = "auth_session"
	}
	if options.Duration == 0 {
		options.Duration = time.Hour * 24 * 30
	}
	return &Auth[S, U]{storage, options}
}

// Authenticate gets and validates the session token from the request cookies.
// Authenticate returns the session, user and a bool to indicate a valid session token.
// On an invalid session token Authenticate deletes the session cookie.
func (a *Auth[S, U]) Authenticate(w http.ResponseWriter, r *http.Request) (S, U, bool) {
	var session S
	var user U
	// get session token from cookies
	token, ok := a.GetCookie(r)
	if !ok {
		return session, user, false
	}
	// validate session token
	session, user, err := a.ValidateSessionToken(token)
	if err != nil {
		// delete invalid session
		a.DeleteCookie(w)
		return session, user, false
	}
	// set cookie to potentially refresh session
	a.SetCookie(w, token, session.GetExpirestAt())
	return session, user, true
}

// GetCookie gets the session cookie.
func (a *Auth[S, U]) GetCookie(r *http.Request) (string, bool) {
	token, err := r.Cookie(a.options.Cookie)
	if err != nil {
		return "", false
	}
	return token.Value, err == nil
}

// SetCookie sets the session cookie.
func (a *Auth[S, U]) SetCookie(w http.ResponseWriter, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     a.options.Cookie,
		HttpOnly: true,
		Secure:   a.options.Secure,
		SameSite: http.SameSiteLaxMode,
		Value:    token,
		Expires:  expiresAt.UTC(),
	})
}

// DeleteCookie deletes the session cookie.
func (a *Auth[S, U]) DeleteCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     a.options.Cookie,
		HttpOnly: true,
		Secure:   a.options.Secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

// GenerateSessionToken generates a secure random string the can be used as a session token.
func (a *Auth[S, U]) GenerateSessionToken() (string, error) {
	// generate random bytes
	b := make([]byte, 20)
	_, err := rand.Read(b)
	// encode with base32 no padding lowercase
	encoding := base32.StdEncoding.WithPadding(base32.NoPadding)
	encoded := encoding.EncodeToString(b)
	return strings.ToLower(encoded), err
}

// CreateSession creates a new session and persists it to storage.
func (a *Auth[S, U]) CreateSession(token string, user U) (S, error) {
	hashedToken := hashSessionToken([]byte(token))
	// persist in storage
	return a.storage.CreateSession(hashedToken, user, time.Now().Add(a.options.Duration))
}

// ValidateSessionToken validates a session token and returns the session and user.
// ValidateSessionToken may update the session's expiresAt timestamp so make sure to set new cookies.
// ValidateSessionToken also deletes expired sessions.
func (a *Auth[S, U]) ValidateSessionToken(token string) (S, U, error) {
	hashedToken := hashSessionToken([]byte(token))
	// get session from db
	session, user, err := a.storage.GetSessionAndUser(hashedToken)
	if err != nil {
		return session, user, err
	}
	// check if session is expired
	expiresAt := session.GetExpirestAt()
	if time.Now().After(expiresAt) {
		a.storage.DeleteSession(hashedToken)
		return session, user, ErrSessionExpired
	}
	// refresh session
	halfTime := expiresAt.Add(-a.options.Duration / 2)
	if time.Now().After(halfTime) {
		newExpiresAt := time.Now().Add(a.options.Duration)
		if session, err = a.storage.UpdateSession(session, newExpiresAt); err != nil {
			return session, user, err
		}
	}
	return session, user, nil
}

// InvalidateSession deletes a session from storage.
func (a *Auth[S, U]) InvalidateSession(id string) error {
	return a.storage.DeleteSession(id)
}

// HashPassword hashes a password using bcrypt.
func (a *Auth[S, U]) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// ComparePassword compares a bcrypt hashed password with a possible plaintext equivalent.
func (a *Auth[S, U]) ComparePassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func hashSessionToken(b []byte) string {
	h := sha256.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}
