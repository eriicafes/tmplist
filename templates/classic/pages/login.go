package classic_pages

import (
	"github.com/eriicafes/tmplist/httperrors"
	"github.com/eriicafes/tmplist/schemas"
)

type Login struct {
	Form   schemas.LoginData
	Error  string
	Errors httperrors.Details
}

func (l Login) Template() (string, any) {
	return "classic/pages/login", l
}
