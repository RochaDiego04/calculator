package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	health := flag.Bool("health", false, "probe the local health endpoint and exit")
	flag.Parse()

	if *health {
		if err := probeHealth(); err != nil {
			fmt.Fprintln(os.Stderr, "health probe failed:", err)
			os.Exit(1)
		}
		return
	}

	run()
}
