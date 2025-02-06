package enhanced

import (
	"github.com/eriicafes/tmpl"
	"github.com/eriicafes/tmplist/internal/httperrors"
	"github.com/eriicafes/tmplist/schemas"
)

type LoginForm struct {
	Form   schemas.LoginData
	Error  string
	Errors httperrors.Details
}

func (l LoginForm) Tmpl() tmpl.Template {
	return tmpl.Associated("enhanced/login", "form", l)
}

type Login struct {
	Layout
	LoginForm
}

func (l Login) Tmpl() tmpl.Template {
	return tmpl.Wrap(&l.Layout, tmpl.Tmpl("enhanced/login", l))
}
