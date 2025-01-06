package main

import (
	"os"

	"github.com/eriicafes/tmpl"
	"github.com/eriicafes/tmpl/vite"
)

type Templates struct {
	tmpl.Templates
	Vite *vite.Vite
}

func setupTemplates(dev bool) Templates {
	v, err := vite.New(vite.Config{
		Dev:    dev,
		Output: os.DirFS("frontend/dist"),
	})
	if err != nil {
		panic(err)
	}
	tp := tmpl.New(os.DirFS("templates")).
		Funcs(v.Funcs()).
		Load("spa").
		MustParse()

	return Templates{tp, v}
}
