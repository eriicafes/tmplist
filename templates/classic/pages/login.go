package classic_pages

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

func (l Login) Template() (string, any) {
	return tmpl.Tmpl("classic/pages/login", l.Layout, l).Template()
}
