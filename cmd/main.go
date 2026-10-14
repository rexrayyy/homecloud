// Command homecloud runs either the headless server or native desktop client.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"homecloud/internal/client"
	"homecloud/internal/server"
)

func main() {
	defaultMode := ""
	if runtime.GOOS == "windows" {
		defaultMode = "client"
	}
	mode := flag.String("mode", defaultMode, "operating mode: server or client")
	addr := flag.String("addr", ":8080", "server listen address")
	flag.Parse()
	switch *mode {
	case "server":
		if err := server.Run(*addr); err != nil {
			log.Fatal(err)
		}
	case "client":
		if err := client.Run(); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Fprintln(os.Stderr, "usage: homecloud --mode=server|client")
		os.Exit(2)
	}
}
