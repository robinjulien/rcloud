package api

import (
	"net/http"

	"github.com/robinjulien/rcloud/internal/web/api/auth"
)

// GetAPIMux generally used with http.stripprefix that is why no prefix on routes
func GetAPIMux() *http.ServeMux {
	var mux *http.ServeMux = new(http.ServeMux)

	mux.HandleFunc("/apitest", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok api"))
	})

	mux.Handle("/auth/", http.StripPrefix("/auth", auth.Handler()))

	return mux
}
