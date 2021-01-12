package ui

import "net/http"

var guifs http.FileSystem

func GetGuiFS() http.FileSystem {
	return guifs
}
