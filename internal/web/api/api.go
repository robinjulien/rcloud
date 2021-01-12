package api

import "net/http"

// GETAPIMux generally used with http.stripprefix that is why no prefix on routes
func GetAPIMux() *http.ServeMux {
	var mux *http.ServeMux = new(http.ServeMux)

	mux.HandleFunc("/apitest", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok api"))
	})

	return mux
}
