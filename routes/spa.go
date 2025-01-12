package routes

import (
	"net/http"

	"github.com/eriicafes/tmpl"
)

func (c Context) MountSpa(mux Mux) {
	mux.HandleFunc("/", combine(func(w http.ResponseWriter, r *http.Request) error {
		tr := c.Renderer()
		spas := map[string]string{
			"react":  "src/react/index.tsx",
			"svelte": "src/svelte/index.ts",
			"vue":    "src/vue/index.ts",
		}
		return tr.Render(w, tmpl.Tmpl("spa/index", tmpl.Map{
			"File": spas[r.URL.Query().Get("spa")],
		}))
	}))
}
