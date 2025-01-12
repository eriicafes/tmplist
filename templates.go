package main

import (
	"html/template"
	"os"

	"github.com/eriicafes/tmpl"
	"github.com/eriicafes/tmpl/vite"
)

func setupTemplates(dev bool) (tmpl.Templates, *vite.Vite) {
	v, err := vite.New(vite.Config{
		Dev:    dev,
		Output: os.DirFS("frontend/dist"),
	})
	if err != nil {
		panic(err)
	}
	funcMap := template.FuncMap{
		"vite_is_dev": func() bool { return v.Dev },
	}
	tp := tmpl.New(os.DirFS("templates")).
		Funcs(v.Funcs(), funcMap).
		Autoload("components").
		Load("spa/index").
		LoadTree("classic/pages").
		LoadTree("enhanced/pages").
		MustParse()

	return tp, v
}
