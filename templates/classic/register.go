package classic

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

func (r Register) Tmpl() tmpl.Template {
	return tmpl.Wrap(&r.Layout, tmpl.Tmpl("classic/register", r))
}
