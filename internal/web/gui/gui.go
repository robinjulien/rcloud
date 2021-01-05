//+build !prod

package gui

import (
	"fmt"
	"net/http"
)

func init() {
	fmt.Println("dev")
	guifs = func() http.FileSystem {
		return http.Dir("./internal/web/gui/gui")
	}()
}
