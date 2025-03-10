package schemas

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
)

type TopicData struct {
	Topic string     `json:"topic"`
	Todos []TodoData `json:"todos"`
}

func (d TopicData) Validate() error {
	return v.ValidateStruct(&d,
		v.Field(&d.Topic, v.Required),
		v.Field(&d.Todos),
	)
}

type TodoData struct {
	Text    string `json:"text"`
	Checked bool   `json:"checked"`
}

func (d TodoData) Validate() error {
	return v.ValidateStruct(&d,
		v.Field(&d.Text, v.Required),
	)
}
