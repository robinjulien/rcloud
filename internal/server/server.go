package server

import (
	"net/http"

	"github.com/robinjulien/rcloud/internal/api"
	"github.com/robinjulien/rcloud/internal/ui"
)

// Serve set up directorypath, routes and launch the server
func Serve(directorypath string, databasepath string, port string) {
	api.SetUp(directorypath, databasepath)
	router := http.NewServeMux()
	router.Handle("/", http.FileServer(ui.FS()))
	router.Handle("/api/", http.StripPrefix("/api", api.Handler()))
	http.ListenAndServe(":"+port, router)
}
