package filemanager

import "net/http"

// Handler returns handler of the filemanager part of the api
func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ls", Ls)
	return mux
}
