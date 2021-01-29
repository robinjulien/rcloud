package main

import (
	"fmt"
	"os"

	"github.com/robinjulien/rcloud/internal/web/server"
)

// usage ./rcloud <directory path> <database path> [-p port]
func main() {
	if len(os.Args) == 3 {
		server.Serve()
	} else {
		fmt.Fprintf(os.Stderr, "Inavlid arguments.")
	}
}
