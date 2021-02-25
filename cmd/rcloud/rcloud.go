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
	https := flag.Bool("https", false, "enable https mode")
	flag.Parse()

	if len(flag.Args()) == 2 {
		server.Serve(flag.Arg(0), flag.Arg(1), strconv.Itoa(*port), *https)
	} else {
		fmt.Fprintf(os.Stderr, "Inavlid arguments.")
	}
}
