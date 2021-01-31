package filemanager

import (
	"net/http"
	"path/filepath"
	"strings"
)

// Handler returns handler of the filemanager part of the api
func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ls", Ls)
	mux.HandleFunc("/mkdir", Mkdir)
	return mux
}

// SantizePath take an input path and return a sanitized version of it
// It filepath.Clean the path, remove all .. and ~
func SanitizePath(path string) string {
	res := "./" + path
	res = strings.ReplaceAll(res, "~", "")
	res = strings.ReplaceAll(res, "..", "")
	res = filepath.Clean(res)

	return res
}
