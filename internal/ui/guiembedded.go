//+build prod

package ui

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed gui/dist/gui/*
var fsys embed.FS

func FS() http.FileSystem {
	guifsys, err := fs.Sub(fsys, "gui/dist/gui")

	if err != nil {
		panic(err)
	}

	return http.FS(guifsys)
}
