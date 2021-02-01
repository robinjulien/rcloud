/*
In the filemanager api, permissions errors are treated as server configuration errors, thus 500.
A file not found is treated as not found error, thus 404
A file that already exists when it should be created is 200 whereas it is 201 if it is created.
*/
package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Permissions for new files/folders
var (
	FilePerm os.FileMode = 664 // rw-rw-r--
	DirPerm  os.FileMode = 775 // rwxrwxr-x
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

type responseLs struct {
	BaseResponse
	Dir []fileType `json:"dir"`
}

// Ls /fm/ls handler, list files of directory
func Ls(w http.ResponseWriter, r *http.Request) {
	path := SanitizePath(r.FormValue("path"))

	files, err := os.ReadDir(path)

	if err != nil {
		RespondJSON(w, BaseResponse{Success: false, ErrorMessage: err.Error()})
		return
	}

	res := responseLs{}
	res.Success = true
	res.Dir = make([]fileType, 0, len(files))

	for _, file := range files {
		info, err := file.Info()

		if err != nil {
			RespondJSON(w, BaseResponse{Success: false, ErrorMessage: err.Error()})
			return
		}

		res.Dir = append(res.Dir, fileType{Name: file.Name(), IsDir: file.IsDir(), Size: info.Size()})
	}

	RespondJSON(w, res)
}

// Mkdir /fm/mkdir create a directory and all of its parents if necessary
func Mkdir(w http.ResponseWriter, r *http.Request) {
	path := SanitizePath(r.PostFormValue("path"))

	err := os.MkdirAll(path, DirPerm)

	if err != nil {
		RespondJSON(w, BaseResponse{Success: false, ErrorMessage: err.Error()})
		return
	}

	RespondJSON(w, BaseResponse{Success: true})

}

// Touch /fm/touch create a file given its path, if the file doesn't already exists
func Touch(w http.ResponseWriter, r *http.Request) {
	path := SanitizePath(r.PostFormValue("path"))

	file, err := os.OpenFile(path, os.O_CREATE, FilePerm)
	file.Close()

	if err != nil {
		RespondJSON(w, BaseResponse{Success: false, ErrorMessage: err.Error()})
		return
	}

	RespondJSON(w, BaseResponse{Success: true})
}

// Rm /fm/rm remove a file or a directory and all of its children
func Rm(w http.ResponseWriter, r *http.Request) {
	path := SanitizePath(r.PostFormValue("path"))

	err := os.RemoveAll(path)

	if err != nil {
		RespondJSON(w, BaseResponse{Success: false, ErrorMessage: err.Error()})
		return
	}

	RespondJSON(w, BaseResponse{Success: true})
}

// Cp /fm/cp copy src to dest
func Cp(w http.ResponseWriter, r *http.Request) {
	src := SanitizePath(r.PostFormValue("src"))
	dest := SanitizePath(r.PostFormValue("dest"))

	fsrc, err := os.Open(src)
	defer fsrc.Close()

	if err != nil {
		RespondJSON(w, BaseResponse{Success: false, ErrorMessage: err.Error()})
		return
	}

	fdest, err := os.Create(dest)
	defer fdest.Close()

	if err != nil {
		RespondJSON(w, BaseResponse{Success: false, ErrorMessage: err.Error()})
		return
	}

	_, err = io.Copy(fdest, fsrc)

	if err != nil {
		RespondJSON(w, BaseResponse{Success: false, ErrorMessage: err.Error()})
		return
	}

	RespondJSON(w, BaseResponse{Success: true})
}

// Mv /fm/mv renames or moves a file or a directory
func Mv(w http.ResponseWriter, r *http.Request) {
	src := SanitizePath(r.PostFormValue("src"))
	dest := SanitizePath(r.PostFormValue("dest"))

	if _, err := os.Stat(dest); err == nil {
		RespondJSON(w, BaseResponse{Success: false, ErrorMessage: "dest already exists"})
		return
	}

	err := os.Rename(src, dest)

	if err != nil {
		RespondJSON(w, BaseResponse{Success: false, ErrorMessage: err.Error()})
		return
	}

	RespondJSON(w, BaseResponse{Success: true})
}

type responseCat struct {
	BaseResponse
	Content string `json:"content"`
}

// Cat /fm/cat returns the content of a file as characters
func Cat(w http.ResponseWriter, r *http.Request) {
	path := SanitizePath(r.FormValue("path"))

	stats, err := os.Stat(path)

	if err != nil {
		RespondJSON(w, BaseResponse{Success: false, ErrorMessage: err.Error()})
		return
	}

	if stats.IsDir() {
		RespondJSON(w, BaseResponse{Success: false, ErrorMessage: "requested ressource is a directory"})
		return
	}

	/*if stats.Size() > { // This function is not currently safe as it loads an entire file in the RAM

	}*/

	content, err := os.ReadFile(path)

	if err != nil {
		RespondJSON(w, BaseResponse{Success: false, ErrorMessage: err.Error()})
		return
	}

	res := responseCat{}
	res.Success = true
	res.Content = string(content)

	RespondJSON(w, res)
}

// Echo /fm/echo write content to a file
func Echo(w http.ResponseWriter, r *http.Request) {
	path := SanitizePath(r.PostFormValue("path"))
	content := r.PostFormValue("content")

	err := os.WriteFile(path, []byte(content), FilePerm)

	if err != nil {
		RespondJSON(w, BaseResponse{Success: false, ErrorMessage: err.Error()})
		return
	}

	RespondJSON(w, BaseResponse{Success: true})
}

// Download /fm/download downloads a file
func Download(w http.ResponseWriter, r *http.Request) {
	path := SanitizePath(r.FormValue("path"))

	stats, err := os.Stat(path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			fmt.Fprintln(os.Stderr, err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if stats.IsDir() {
		RespondJSON(w, BaseResponse{Success: false, ErrorMessage: "requested ressource is a directory"})
		return
	}

	r.URL.Path = ""
	http.ServeFile(w, r, path)
}

// Upload /fm/upload uploads multiple files
func Upload(w http.ResponseWriter, r *http.Request) {

}
