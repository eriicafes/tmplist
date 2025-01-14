package classic_pages

import (
	"github.com/eriicafes/tmpl"
	"github.com/eriicafes/tmplist/db"
)

type Topic struct {
	Layout
	Topic db.Topic
	Todos []db.Todo
}

func (t Topic) Template() (string, any) {
	return tmpl.Tmpl("classic/pages/topic", t.Layout, t).Template()
}
