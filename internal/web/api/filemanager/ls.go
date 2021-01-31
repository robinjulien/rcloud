package filemanager

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/robinjulien/rcloud/internal/web/api/auth"
	"github.com/robinjulien/rcloud/internal/web/api/common"
)

type fileType struct {
	IsDir bool   `json:"isDir"`
	Name  string `json:"name"`
	Size  int64  `json:""`
}

type responseLs []fileType

// Ls /fm/ls handler, list files of directory
func Ls(w http.ResponseWriter, r *http.Request) {
	if !common.CheckMethod(w, r, "GET") {
		return
	}

	path := SanitizePath(r.FormValue("path"))

	u := auth.GetUserByCookies(r)

	if u == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	files, err := ioutil.ReadDir(path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			w.WriteHeader(http.StatusNotFound)
		} else if errors.Is(err, os.ErrPermission) {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	res := make(responseLs, 0, len(files))

	for _, file := range files {
		res = append(res, fileType{Name: file.Name(), IsDir: file.IsDir(), Size: file.Size()})
	}

	common.RespondJSON(w, res)
}
