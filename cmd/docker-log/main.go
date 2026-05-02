package main

import (
	"fmt"
	"os"
)

const version = "0.1.0-dev"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--help", "-h", "help":
			printHelp()
			return
		case "--version", "version":
			fmt.Fprintln(os.Stdout, version)
			return
		}
	}

	fmt.Fprintln(os.Stderr, "docker-log: interactive log viewer is not implemented yet")
	os.Exit(1)
}

func printHelp() {
	fmt.Fprintln(os.Stdout, `docker-log streams and filters logs from running Docker containers.

Usage:
  docker-log [--help] [--version]

Options:
  -h, --help     Show help
      --version  Show version`)
}
