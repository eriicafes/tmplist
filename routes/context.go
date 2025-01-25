package routes

import (
	"github.com/eriicafes/tmpl"
	"github.com/eriicafes/tmplist/db"
	"github.com/eriicafes/tmplist/internal"
	"github.com/eriicafes/tmplist/internal/session"
)

var (
	requestUser    internal.ContextValue[db.User]    = "user"
	requestSession internal.ContextValue[db.Session] = "session"
)

type Context struct {
	tmpl.Templates
	DB   db.DB
	Auth *session.Auth[db.Session, db.User]
	Prod bool
}
