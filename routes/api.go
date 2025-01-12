package routes

import (
	"fmt"
	"net/http"
)

func (c Context) MountApi(mux Mux) {
	mux.HandleFunc("/", combine(func(w http.ResponseWriter, r *http.Request) error {
		fmt.Fprint(w, "api routes")
		return nil
	}))
}
