package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/robinjulien/rcloud/internal/server"
)

// usage ./rcloud <directory path> <database path> [-p port]
func main() {
	port := flag.Int("port", 80, "port http server will be running on")
	if len(os.Args) == 3 {
		server.Serve(os.Args[1], os.Args[2], strconv.Itoa(*port))
	} else {
		fmt.Fprintf(os.Stderr, "Inavlid arguments.")
	}
}
