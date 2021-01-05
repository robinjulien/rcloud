package angularhandler

import (
	"net/http"
	"os"
	"path"
)

// AngularHandler a commenter
type AngularHandler struct {
	fs http.FileSystem
}

// New a commenter
func New(fs http.FileSystem) *AngularHandler {
	return &AngularHandler{fs}
}

func (ah *AngularHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, err := ah.fs.Open("." + path.Clean(r.URL.Path))
	if err != nil {
		println("Erreur1 :", err.Error())
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 pas trouve"))
			return
		}
	}

	stats, err := f.Stat()
	if stats.IsDir() {
		w.Write([]byte("Dossier"))
		return
	}
	if err != nil {
		println("Erreur2 :", err.Error())
	}

	http.ServeContent(w, r, stats.Name(), stats.ModTime(), f)
}
