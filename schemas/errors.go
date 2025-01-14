package schemas

import (
	"github.com/eriicafes/tmplist/internal/httperrors"
	v "github.com/go-ozzo/ozzo-validation/v4"
)

func FormErrors(err error) httperrors.Details {
	if err == nil {
		return nil
	}
	var details = make(httperrors.Details)
	for k, v := range err.(v.Errors) {
		details[k] = v.Error()
	}
	return details
}
