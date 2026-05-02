package main

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	dockeradapter "github.com/paknahad/docker-log/internal/docker"
	"github.com/paknahad/docker-log/internal/ui"
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

	client, err := dockeradapter.NewClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "docker-log: %v\n", err)
		os.Exit(1)
	}

	containers, err := client.ListRunningContainers(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "docker-log: %v\n", err)
		os.Exit(1)
	}

	if _, err := tea.NewProgram(ui.NewSelectionModel(containers)).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "docker-log: run selection UI: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Fprintln(os.Stdout, `docker-log streams and filters logs from running Docker containers.

Usage:
  docker-log [--help] [--version]

Options:
  -h, --help     Show help
      --version  Show version`)
}
