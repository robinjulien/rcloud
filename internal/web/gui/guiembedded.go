//+build prod

package gui

import (
	"embed"
	"fmt"
	"net/http"
)

func init() {
	fmt.Println("prod")
	guifs = func() http.FileSystem {
		//go:embed gui
		var fs embed.FS
		return http.FS(fs)
	}()
}
