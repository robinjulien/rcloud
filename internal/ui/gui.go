//+build !prod

package ui

import (
	"net/http"
)

// FS is filesystem of gui, loaded dynamically
func FS() http.FileSystem {
	return http.Dir("./internal/ui/gui")
}
