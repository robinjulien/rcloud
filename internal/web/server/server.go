package server

import (
	"net/http"

	"github.com/robinjulien/rcloud/internal/web/gui"
)

func Serve() {
	http.HandleFunc("/gui/", func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(gui.GetGuiFS()).ServeHTTP(w, r)
	})
	http.ListenAndServe(":80", nil)
}
