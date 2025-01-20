package classic_pages

import "github.com/eriicafes/tmpl"

type Error struct {
	Layout
	Title string
}

func (e Error) Template() (string, any) {
	return tmpl.Tmpl("classic/pages/error", e.Layout, e).Template()
}
