package classic

import (
	"github.com/eriicafes/tmpl"
	"github.com/eriicafes/tmplist/internal/httperrors"
	"github.com/eriicafes/tmplist/schemas"
)

type Login struct {
	Layout
	Form   schemas.LoginData
	Error  string
	Errors httperrors.Details
}

func (l Login) Tmpl() tmpl.Template {
	return tmpl.Wrap(&l.Layout, tmpl.Tmpl("classic/login", l))
}
