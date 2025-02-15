package schemas

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (d LoginData) Validate() error {
	return v.ValidateStruct(&d,
		v.Field(&d.Email, v.Required, is.EmailFormat),
		v.Field(&d.Password, v.Required),
	)
}

type RegisterData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (d RegisterData) Validate() error {
	return v.ValidateStruct(&d,
		v.Field(&d.Email, v.Required, is.EmailFormat),
		v.Field(&d.Password, v.Required, v.Length(8, 0)),
	)
}
