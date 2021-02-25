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
	router := http.NewServeMux()
	router.Handle("/ls", MethodMiddleware("GET", http.HandlerFunc(Ls)))
	router.Handle("/mkdir", MethodMiddleware("POST", http.HandlerFunc(Mkdir)))
	router.Handle("/touch", MethodMiddleware("POST", http.HandlerFunc(Touch)))
	router.Handle("/rm", MethodMiddleware("POST", http.HandlerFunc(Rm)))
	router.Handle("/cp", MethodMiddleware("POST", http.HandlerFunc(Cp)))
	router.Handle("/mv", MethodMiddleware("POST", http.HandlerFunc(Mv)))
	router.Handle("/cat", MethodMiddleware("GET", http.HandlerFunc(Cat)))
	router.Handle("/echo", MethodMiddleware("POST", http.HandlerFunc(Echo)))
	router.Handle("/download", MethodMiddleware("GET", http.HandlerFunc(Download)))
	router.Handle("/upload", MethodMiddleware("POST", http.HandlerFunc(Upload)))
	return router
}

// SanitizePath take an input path and return a sanitized version of it
// It filepath.Clean the path, remove all .. and ~
func SanitizePath(path string) string {
	res := "./" + path
	res = strings.ReplaceAll(res, "~", "")
	res = strings.ReplaceAll(res, "..", "")
	res = filepath.Clean(res)

	if filepath.IsAbs(res) {
		return "."
	}

	return res
}

type fileType struct {
	IsDir bool   `json:"isDir"`
	Name  string `json:"name"`
	Size  int64  `json:"size"`
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

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		RespondJSON(w, BaseResponse{Success: false, ErrorMessage: "directory already exists"})
		return
	}

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

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		RespondJSON(w, BaseResponse{Success: false, ErrorMessage: "file already exists"})
		return
	}

	file, err := os.OpenFile(path, os.O_CREATE, 0666)
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

	if stats, _ := fsrc.Stat(); stats.IsDir() {
		RespondJSON(w, BaseResponse{Success: false, ErrorMessage: "src is a directory"})
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
		RespondJSON(w, responseCat{BaseResponse: BaseResponse{Success: false, ErrorMessage: err.Error()}})
		return
	}

	if stats.IsDir() {
		RespondJSON(w, responseCat{BaseResponse: BaseResponse{Success: false, ErrorMessage: "requested ressource is a directory"}})
		return
	}

	/*if stats.Size() > { // This function is not currently safe as it loads an entire file in the RAM

	}*/

	content, err := os.ReadFile(path)

	if err != nil {
		RespondJSON(w, responseCat{BaseResponse: BaseResponse{Success: false, ErrorMessage: err.Error()}})
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
	// Parse our multipart form, 10 << 20 specifies a maximum of 10MB of its file parts stored in the memory, the remainder will be stored in temporary files
	// This does NOT specifies a maximum file or upload size. To do so, you have to use http.MaxByteReader.
	errMultipartForm := r.ParseMultipartForm(10 << 20)

	if errMultipartForm != nil {
		RespondJSON(w, BaseResponse{Success: false, ErrorMessage: errMultipartForm.Error()})
		return
	}

	path := SanitizePath(r.PostFormValue("path"))

	files := r.MultipartForm.File["files[]"]

	for i := range files { // loop through the files one by one
		file, errFile := files[i].Open()
		defer file.Close()

		if errFile != nil {
			RespondJSON(w, BaseResponse{Success: false, ErrorMessage: errFile.Error()})
			return
		}

		out, errOut := os.Create(path + "/" + files[i].Filename)
		defer out.Close()

		if errOut != nil {
			RespondJSON(w, BaseResponse{Success: false, ErrorMessage: errOut.Error()})
			return
		}

		_, errCopy := io.Copy(out, file)

		if errCopy != nil {
			RespondJSON(w, BaseResponse{Success: false, ErrorMessage: errCopy.Error()})
			return
		}
	}

	RespondJSON(w, BaseResponse{Success: true})
}
