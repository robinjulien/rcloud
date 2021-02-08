package angularhandler

import (
	"net/http"
	"os"
)

// AngularHandler is the type of handler that serves angular apps
type AngularHandler struct {
	fs http.FileSystem
}

// New creates a new instance of angular handler
func New(fs http.FileSystem) *AngularHandler {
	return &AngularHandler{fs}
}

// ServeIndex serves index.html under the / path
func (ah *AngularHandler) ServeIndex(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "/"
	index, err := ah.fs.Open("index.html")

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	stats, err := index.Stat()

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	http.ServeContent(w, r, stats.Name(), stats.ModTime(), index)
}

// ServeHTTP implements the http.Handler interface
// It checks if the requested file exists, if not or if it is a directory, it serves the index
// If an error occurs, it responds with 403 Forbidden
func (ah *AngularHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, err := ah.fs.Open(r.URL.Path)
	if err != nil {
		if os.IsNotExist(err) {
			ah.ServeIndex(w, r)
			return
		}
		w.WriteHeader(http.StatusForbidden)
		return
	}

	stats, err := f.Stat()

	if stats.IsDir() {
		ah.ServeIndex(w, r)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	http.ServeContent(w, r, stats.Name(), stats.ModTime(), f)
}
