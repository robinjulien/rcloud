package api

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// FmHandler returns handler of the filemanager part of the api
func FmHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/ls", MethodMiddleware("GET", http.HandlerFunc(Ls)))
	mux.Handle("/mkdir", MethodMiddleware("POST", http.HandlerFunc(Mkdir)))
	return mux
}

// SanitizePath take an input path and return a sanitized version of it
// It filepath.Clean the path, remove all .. and ~
func SanitizePath(path string) string {
	res := "./" + path
	res = strings.ReplaceAll(res, "~", "")
	res = strings.ReplaceAll(res, "..", "")
	res = filepath.Clean(res)

	return res
}

type fileType struct {
	IsDir bool   `json:"isDir"`
	Name  string `json:"name"`
	Size  int64  `json:""`
}

type responseLs []fileType

// Ls /fm/ls handler, list files of directory
func Ls(w http.ResponseWriter, r *http.Request) {
	path := SanitizePath(r.FormValue("path"))

	files, err := os.ReadDir(path)

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
		info, err := file.Info()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		res = append(res, fileType{Name: file.Name(), IsDir: file.IsDir(), Size: info.Size()})
	}

	RespondJSON(w, res)
}

// Mkdir /fm/mkdir create a directory and all of its parents if necessary
func Mkdir(w http.ResponseWriter, r *http.Request) {
	path := SanitizePath(r.PostFormValue("path"))

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
