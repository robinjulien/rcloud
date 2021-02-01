package main

import (
	"fmt"
	"os"

	"github.com/robinjulien/rcloud/internal/server"
)

// usage ./rcloud <directory path> <database path> [-p port]
func main() {
	if len(os.Args) == 3 {
		server.Serve(os.Args[1], os.Args[2], "80")
	} else {
		fmt.Fprintf(os.Stderr, "Inavlid arguments.")
	}
}
