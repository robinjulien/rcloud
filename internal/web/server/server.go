package server

import (
	"net/http"

	"github.com/robinjulien/rcloud/internal/web/gui"
)

func Serve() {
	http.Handle("/", http.FileServer(gui.GetGuiFS()))
	http.ListenAndServe(":80", nil)
}
