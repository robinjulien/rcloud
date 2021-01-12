//+build prod

package ui

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

func init() {
	fmt.Println("prod")
	guifs = func() http.FileSystem {
		//go:embed gui/*
		var fsys embed.FS
		guifsys, err := fs.Sub(fsys, "gui")

		if err != nil {
			panic(err)
		}

		return http.FS(guifsys)
	}()
}
