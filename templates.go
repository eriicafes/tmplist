package main

import (
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
	tp := tmpl.New(os.DirFS("templates")).
		Funcs(v.Funcs()).
		Autoload("components").
		Load("index").
		Load("spa/index").
		LoadTree("classic").
		LoadTree("enhanced").
		MustParse()

	return tp, v
}
