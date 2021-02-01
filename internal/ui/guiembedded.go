//+build prod

package ui

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed gui/*
var fsys embed.FS

func FS() http.FileSystem {
	guifsys, err := fs.Sub(fsys, "gui")

	if err != nil {
		panic(err)
	}

	return http.FS(guifsys)
}
