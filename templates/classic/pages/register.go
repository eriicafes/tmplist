package classic_pages

import (
	"github.com/eriicafes/tmpl"
	"github.com/eriicafes/tmplist/internal/httperrors"
	"github.com/eriicafes/tmplist/schemas"
)

type Register struct {
	Layout
	Form   schemas.RegisterData
	Error  string
	Errors httperrors.Details
}

func (r Register) Template() (string, any) {
	return tmpl.Tmpl("classic/pages/register", r.Layout, r).Template()
}
