package main

import (
	"context"
	"io"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/paknahad/docker-log/internal/domain"
	"github.com/paknahad/docker-log/internal/stream"
	"github.com/paknahad/docker-log/internal/ui"
)

func TestRunReturnsWithoutStreamingWhenNoContainersSelected(t *testing.T) {
	client := &fakeDockerClient{
		containers: []domain.Container{
			{ID: "api-id", Name: "api"},
		},
	}

	err := run(
		context.Background(),
		client,
		func(containers []domain.Container) ([]domain.Container, error) {
			if len(containers) != 1 || containers[0].ID != "api-id" {
				t.Fatalf("selection containers = %#v, want discovered api container", containers)
			}
			return nil, nil
		},
		func(<-chan stream.Event) error {
			t.Fatal("log viewer started with no selected containers")
			return nil
		},
	)
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}
	if client.opened != nil {
		t.Fatalf("opened streams = %#v, want none", client.opened)
	}
}

func TestRunStreamsSelectedContainersIntoLogView(t *testing.T) {
	client := &fakeDockerClient{
		containers: []domain.Container{
			{ID: "api-id", Name: "api"},
			{ID: "worker-id", Name: "worker"},
		},
		logs: map[string]string{
			"worker-id": "ready\n",
		},
	}

	var gotEvents []stream.Event
	err := run(
		context.Background(),
		client,
		func(containers []domain.Container) ([]domain.Container, error) {
			return []domain.Container{containers[1]}, nil
		},
		func(events <-chan stream.Event) error {
			for event := range events {
				gotEvents = append(gotEvents, event)
			}
			return nil
		},
	)
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}
	if len(client.opened) != 1 || client.opened[0] != "worker-id" {
		t.Fatalf("opened streams = %#v, want worker-id", client.opened)
	}
	logEvents := logLineEvents(gotEvents)
	if len(logEvents) != 1 {
		t.Fatalf("len(log events) = %d, want 1: %#v", len(logEvents), gotEvents)
	}
	if logEvents[0].Container != "worker" || logEvents[0].Line != "worker: ready" {
		t.Fatalf("event = %#v, want worker prefixed ready line", logEvents[0])
	}
}

func TestSelectionResultContainersReturnsNoneWhenSelectionCancelled(t *testing.T) {
	model := ui.NewSelectionModel([]domain.Container{{ID: "api-id", Name: "api"}})
	model = updateSelectionKey(t, model, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})
	model = updateSelectionKey(t, model, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})

	selected := selectionResultContainers(model)
	if len(selected) != 0 {
		t.Fatalf("selectionResultContainers(cancelled) = %#v, want none", selected)
	}
}

func TestSelectionResultContainersReturnsSelectedWhenSelectionStarted(t *testing.T) {
	model := ui.NewSelectionModel([]domain.Container{{ID: "api-id", Name: "api"}})
	model = updateSelectionKey(t, model, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})
	model = updateSelectionKey(t, model, tea.KeyMsg{Type: tea.KeyEnter})

	selected := selectionResultContainers(model)
	if len(selected) != 1 || selected[0].ID != "api-id" {
		t.Fatalf("selectionResultContainers(started) = %#v, want api-id", selected)
	}
}

type fakeDockerClient struct {
	containers []domain.Container
	logs       map[string]string
	opened     []string
}

func (f *fakeDockerClient) ListRunningContainers(context.Context) ([]domain.Container, error) {
	return append([]domain.Container(nil), f.containers...), nil
}

func (f *fakeDockerClient) OpenContainerLogs(_ context.Context, container domain.Container) (io.ReadCloser, error) {
	f.opened = append(f.opened, container.ID)
	return io.NopCloser(strings.NewReader(f.logs[container.ID])), nil
}

func updateSelectionKey(t *testing.T, model ui.SelectionModel, key tea.KeyMsg) ui.SelectionModel {
	t.Helper()

	next, _ := model.Update(key)
	selection, ok := next.(ui.SelectionModel)
	if !ok {
		t.Fatalf("Update() returned %T, want ui.SelectionModel", next)
	}
	return selection
}

func logLineEvents(events []stream.Event) []stream.Event {
	var lines []stream.Event
	for _, event := range events {
		if event.Line != "" {
			lines = append(lines, event)
		}
	}
	return lines
}
