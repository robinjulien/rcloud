package api

import (
	"net/http"
	"os"

	"github.com/robinjulien/rcloud/internal/web/api/auth"
	"github.com/robinjulien/rcloud/internal/web/api/filemanager"
)

// SetUp sets up all the ressources needed for the use of the api
func SetUp(directorypath string, databasepath string) {
	os.Chdir(directorypath)
	auth.SetUp(databasepath)
}

// GetAPIMux generally used with http.stripprefix that is why no prefix on routes
func GetAPIMux() *http.ServeMux {
	var mux *http.ServeMux = new(http.ServeMux)

	mux.HandleFunc("/apitest", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok api"))
	})

	mux.Handle("/auth/", http.StripPrefix("/auth", auth.Handler()))
	mux.Handle("/fm/", http.StripPrefix("/fm", filemanager.Handler()))

	return mux
}
