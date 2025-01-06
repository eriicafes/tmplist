package main

import (
	"log"
	"net/http"

	"github.com/eriicafes/tmpl"
)

type Entry struct {
	Title, File string
}

func main() {
	config := getConfig()
	templates := setupTemplates(!config.Prod)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tr := templates.Renderer()
		spas := map[string]Entry{
			"react":  {"Tmpl React", "src/react/index.tsx"},
			"svelte": {"Tmpl Svelte", "src/svelte/index.ts"},
			"vue":    {"Tmpl Vue", "src/vue/index.ts"},
		}
		entry, ok := spas[r.URL.Query().Get("spa")]
		if !ok {
			entry = Entry{"Tmpl Vanilla", ""}
		}

		if err := tr.Render(w, tmpl.Tmpl("spa", entry)); err != nil {
			log.Println(err)
		}
	})
	http.Handle("/", templates.Vite.ServePublic(handler))
	http.ListenAndServe(config.ListenAddr(), nil)
}
