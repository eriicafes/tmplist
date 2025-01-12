package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/eriicafes/tmplist/db"
)

type AuthService struct {
	DB db.DB
}

func (a *AuthService) GenerateSessionToken() (string, error) {
	// generate random bytes
	b := make([]byte, 20)
	_, err := rand.Read(b)
	// encode with base32 no padding lowercase
	encoding := base32.StdEncoding.WithPadding(base32.NoPadding)
	encoded := encoding.EncodeToString(b)
	return strings.ToLower(encoded), err
}

func (a *AuthService) CreateSession(token string, userId int) (db.Session, error) {
	hashedToken := hashSessionToken([]byte(token))
	// insert session to db
	s := db.Session{
		Id:        hashedToken,
		UserId:    userId,
		ExpiresAt: time.Now().Add(sessionDuration),
	}
	err := a.DB.InsertSession(s)
	if err == db.ErrNotFound {
		return s, fmt.Errorf("user not found")
	}
	return s, nil
}

func (a *AuthService) ValidateSessionToken(token string) (db.Session, db.User, error) {
	hashedToken := hashSessionToken([]byte(token))
	// get session from db
	session, user, err := a.DB.GetSessionAndUser(hashedToken)
	if err != nil {
		return session, user, fmt.Errorf("session not found")
	}
	// check if session is expired
	if time.Now().After(session.ExpiresAt) {
		a.DB.DeleteSession(session.Id)
		return session, user, fmt.Errorf("session expired")
	}
	// refresh session
	halfTime := session.ExpiresAt.Add(-sessionDuration / 2)
	if time.Now().After(halfTime) {
		session.ExpiresAt = time.Now().Add(sessionDuration)
		if err := a.DB.UpdateSession(session); err != nil {
			return session, user, fmt.Errorf("session refresh failed")
		}
	}
	return session, user, nil
}

func (a *AuthService) InvalidateSession(sessionId string) error {
	return a.DB.DeleteSession(sessionId)
}

func (a *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (a *AuthService) ComparePassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func hashSessionToken(b []byte) string {
	h := sha256.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}
