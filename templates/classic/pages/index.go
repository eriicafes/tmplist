package classic_pages

import (
	"github.com/eriicafes/tmpl"
	"github.com/eriicafes/tmplist/db"
)

type Index struct {
	Layout
	User   db.User
	Topics []db.Topic
}

func (i Index) Template() (string, any) {
	return tmpl.Tmpl("classic/pages/index", i.Layout, i).Template()
}
