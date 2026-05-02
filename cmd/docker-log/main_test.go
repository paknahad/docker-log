package main

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/paknahad/docker-log/internal/domain"
	"github.com/paknahad/docker-log/internal/stream"
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
	if len(gotEvents) != 1 {
		t.Fatalf("len(events) = %d, want 1: %#v", len(gotEvents), gotEvents)
	}
	if gotEvents[0].Container != "worker" || gotEvents[0].Line != "worker: ready" {
		t.Fatalf("event = %#v, want worker prefixed ready line", gotEvents[0])
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
