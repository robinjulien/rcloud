package server

import (
	"net/http"

	"github.com/robinjulien/rcloud/internal/web/api"
	"github.com/robinjulien/rcloud/internal/web/ui"
)

func Serve() {
	http.Handle("/", http.FileServer(ui.GetGuiFS()))
	http.Handle("/api/", http.StripPrefix("/api", api.GetAPIMux()))
	http.ListenAndServe(":80", nil)
}
