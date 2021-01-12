//+build !prod

package ui

import (
	"fmt"
	"net/http"
)

func init() {
	fmt.Println("dev")
	guifs = func() http.FileSystem {
		return http.Dir("./internal/web/ui/gui")
	}()
}
