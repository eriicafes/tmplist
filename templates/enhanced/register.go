package enhanced

import (
	"github.com/eriicafes/tmpl"
	"github.com/eriicafes/tmplist/internal/httperrors"
	"github.com/eriicafes/tmplist/schemas"
)

type RegisterForm struct {
	Form   schemas.RegisterData
	Error  string
	Errors httperrors.Details
}

func (r RegisterForm) Tmpl() tmpl.Template {
	return tmpl.Associated("enhanced/register", "form", r)
}

type Register struct {
	Layout
	RegisterForm
}

func (r Register) Tmpl() tmpl.Template {
	return tmpl.Wrap(&r.Layout, tmpl.Tmpl("enhanced/register", r))
}
