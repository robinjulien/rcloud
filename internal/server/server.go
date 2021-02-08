package server

import (
	"net/http"

	"github.com/robinjulien/rcloud/internal/api"
	"github.com/robinjulien/rcloud/internal/ui"
	"github.com/robinjulien/rcloud/pkg/angularhandler"
)

// Serve set up directorypath, routes and launch the server
func Serve(directorypath string, databasepath string, port string) {
	var ah *angularhandler.AngularHandler = angularhandler.New(ui.FS())

	api.SetUp(directorypath, databasepath)
	router := http.NewServeMux()
	router.Handle("/", ah)
	router.Handle("/api/", http.StripPrefix("/api", api.Handler()))

	http.ListenAndServe(":"+port, router)
}
