package classic_pages

import (
	"github.com/eriicafes/tmplist/httperrors"
	"github.com/eriicafes/tmplist/schemas"
)

type Register struct {
	Form   schemas.RegisterData
	Error  string
	Errors httperrors.Details
}

func (l Register) Template() (string, any) {
	return "classic/pages/register", l
}
