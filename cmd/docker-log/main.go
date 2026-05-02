package main

import (
	"context"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	dockeradapter "github.com/paknahad/docker-log/internal/docker"
	"github.com/paknahad/docker-log/internal/domain"
	"github.com/paknahad/docker-log/internal/stream"
	"github.com/paknahad/docker-log/internal/ui"
)

const version = "0.1.0-dev"
const streamBuffer = 128

type dockerClient interface {
	ListRunningContainers(context.Context) ([]domain.Container, error)
	OpenContainerLogs(context.Context, domain.Container) (io.ReadCloser, error)
}

type selectionRunner func([]domain.Container) ([]domain.Container, error)
type logRunner func(<-chan stream.Event) error

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

	if err := run(context.Background(), client, runSelectionUI, runLogUI); err != nil {
		fmt.Fprintf(os.Stderr, "docker-log: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, client dockerClient, selectContainers selectionRunner, viewLogs logRunner) error {
	containers, err := client.ListRunningContainers(ctx)
	if err != nil {
		return err
	}

	selected, err := selectContainers(containers)
	if err != nil {
		return fmt.Errorf("run selection UI: %w", err)
	}
	if len(selected) == 0 {
		return nil
	}

	streamCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	sources := stream.SourcesForContainers(selected, client.OpenContainerLogs)
	events := stream.NewManager(streamBuffer).Start(streamCtx, sources)
	if err := viewLogs(events); err != nil {
		return fmt.Errorf("run log UI: %w", err)
	}
	return nil
}

func runSelectionUI(containers []domain.Container) ([]domain.Container, error) {
	model, err := tea.NewProgram(ui.NewSelectionModel(containers)).Run()
	if err != nil {
		return nil, err
	}

	selection, ok := model.(ui.SelectionModel)
	if !ok {
		return nil, fmt.Errorf("selection UI returned %T", model)
	}
	return selection.SelectedContainers(), nil
}

func runLogUI(events <-chan stream.Event) error {
	_, err := tea.NewProgram(ui.NewLogModel(events)).Run()
	return err
}

func printHelp() {
	fmt.Fprintln(os.Stdout, `docker-log streams and filters logs from running Docker containers.

Usage:
  docker-log [--help] [--version]

Options:
  -h, --help     Show help
      --version  Show version`)
}
