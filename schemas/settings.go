package schemas

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
)

type SettingsData struct {
	Mode  string
	Delay string
}

func (d SettingsData) Validate() error {
	return v.ValidateStruct(&d,
		v.Field(&d.Mode, v.Required, v.In("none", "classic", "enhanced", "spa")),
		v.Field(&d.Delay, v.Required, v.In("normal", "slow", "very_slow")),
	)
}
