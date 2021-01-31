package filemanager

import (
	"errors"
	"net/http"
	"os"

	"github.com/robinjulien/rcloud/internal/web/api/auth"
	"github.com/robinjulien/rcloud/internal/web/api/common"
)

// Mkdir /fm/mkdir create a directory and all of its parents if necessary
func Mkdir(w http.ResponseWriter, r *http.Request) {
	if !common.CheckMethod(w, r, "POST") {
		return
	}

	path := SanitizePath(r.PostFormValue("path"))

	u := auth.GetUserByCookies(r)

	if u == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			w.WriteHeader(http.StatusForbidden)
		} else {
			if f, err2 := os.Stat(path); err2 == nil && !f.IsDir() {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}
