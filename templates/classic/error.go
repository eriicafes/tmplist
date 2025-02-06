package classic

import "github.com/eriicafes/tmpl"

type Error struct {
	Layout
	Title string
}

func (e Error) Tmpl() tmpl.Template {
	return tmpl.Wrap(&e.Layout, tmpl.Tmpl("classic/error", e))
}
