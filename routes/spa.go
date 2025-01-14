package routes

import (
	"net/http"

	"github.com/eriicafes/tmpl"
	"github.com/eriicafes/tmplist/internal"
)

func (c Context) Spa(mux internal.Mux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tr := c.Renderer()
		spas := map[string]string{
			"react":  "src/react/index.tsx",
			"svelte": "src/svelte/index.ts",
			"vue":    "src/vue/index.ts",
		}
		tr.Render(w, tmpl.Tmpl("spa/index", tmpl.Map{
			"File": spas[r.URL.Query().Get("spa")],
		}))
	})
}
