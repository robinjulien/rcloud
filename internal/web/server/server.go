package server

import (
	"net/http"

	"github.com/robinjulien/rcloud/internal/web/api"
	"github.com/robinjulien/rcloud/internal/web/ui"
)

// Serve set up directorypath, routes and launch the server
func Serve(directorypath string, databasepath string, port string) {
	api.SetUp(directorypath, databasepath)
	http.Handle("/", http.FileServer(ui.GetGuiFS()))
	http.Handle("/api/", http.StripPrefix("/api", api.GetAPIMux()))
	http.ListenAndServe(":"+port, nil)
}
