package server

import (
	"net/http"

	"github.com/robinjulien/rcloud/internal/web/ui"
)

func Serve() {
	http.Handle("/", http.FileServer(ui.GetGuiFS()))
	http.ListenAndServe(":80", nil)
}
