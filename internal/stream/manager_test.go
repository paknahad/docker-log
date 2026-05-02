package stream

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/paknahad/docker-log/internal/domain"
)

func TestManagerFansInLinesWithContainerPrefixes(t *testing.T) {
	sources := []Source{
		{
			Container: "api",
			Open: func(context.Context) (io.ReadCloser, error) {
				return io.NopCloser(strings.NewReader("api one\napi two\n")), nil
			},
		},
		{
			Container: "worker",
			Open: func(context.Context) (io.ReadCloser, error) {
				return io.NopCloser(strings.NewReader("worker one\n")), nil
			},
		},
	}

	events := NewManager(0).Start(context.Background(), sources)
	got := collectEvents(events)

	want := map[string]bool{
		"api|api one|api: api one": false,
		"api|api two|api: api two": false,
		"worker|worker one|worker: worker one": false,
	}
	for _, event := range got {
		if event.Err != nil {
			t.Fatalf("event.Err = %v, want nil", event.Err)
		}
		key := event.Container + "|" + event.Message + "|" + event.Line
		if _, ok := want[key]; !ok {
			t.Fatalf("unexpected event: %#v", event)
		}
		want[key] = true
	}
	for key, seen := range want {
		if !seen {
			t.Fatalf("missing event %s in %#v", key, got)
		}
	}
}

func TestManagerReportsSourceFailureWithoutStoppingOtherStreams(t *testing.T) {
	sourceErr := errors.New("stream failed")
	sources := []Source{
		{
			Container: "bad",
			Open: func(context.Context) (io.ReadCloser, error) {
				return nil, sourceErr
			},
		},
		{
			Container: "good",
			Open: func(context.Context) (io.ReadCloser, error) {
				return io.NopCloser(strings.NewReader("still running\n")), nil
			},
		},
	}

	events := NewManager(0).Start(context.Background(), sources)
	got := collectEvents(events)

	var sawError bool
	var sawGoodLine bool
	for _, event := range got {
		if errors.Is(event.Err, sourceErr) && event.Container == "bad" {
			sawError = true
		}
		if event.Container == "good" && event.Line == "good: still running" {
			sawGoodLine = true
		}
	}
	if !sawError {
		t.Fatalf("did not receive source error in %#v", got)
	}
	if !sawGoodLine {
		t.Fatalf("unrelated stream did not continue in %#v", got)
	}
}

func TestManagerClosesLiveStreamsOnContextCancel(t *testing.T) {
	reader, writer := io.Pipe()
	closed := make(chan struct{}, 1)
	source := Source{
		Container: "api",
		Open: func(context.Context) (io.ReadCloser, error) {
			return closeNotifyingReadCloser{ReadCloser: reader, closed: closed}, nil
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	events := NewManager(1).Start(ctx, []Source{source})

	if _, err := writer.Write([]byte("first\n")); err != nil {
		t.Fatalf("write log line: %v", err)
	}
	event := <-events
	if event.Line != "api: first" {
		t.Fatalf("event.Line = %q, want api: first", event.Line)
	}

	cancel()

	select {
	case <-closed:
	case <-time.After(time.Second):
		t.Fatal("stream reader was not closed after context cancellation")
	}
	_ = writer.Close()
}

func TestSourcesForContainersCreatesOneSourcePerSelectedContainer(t *testing.T) {
	containers := []domain.Container{
		{ID: "abc123", Name: "api"},
		{ID: "def456", Name: "worker"},
	}

	sources := SourcesForContainers(containers, func(ctx context.Context, container domain.Container) (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader(container.ID + "\n")), nil
	})

	if len(sources) != 2 {
		t.Fatalf("len(sources) = %d, want 2", len(sources))
	}
	if sources[0].ContainerID != "abc123" || sources[0].Container != "api" {
		t.Fatalf("sources[0] = %#v, want api source", sources[0])
	}
	if sources[1].ContainerID != "def456" || sources[1].Container != "worker" {
		t.Fatalf("sources[1] = %#v, want worker source", sources[1])
	}

	events := collectEvents(NewManager(0).Start(context.Background(), sources))
	if len(events) != 2 {
		t.Fatalf("len(events) = %d, want 2", len(events))
	}
}

func collectEvents(events <-chan Event) []Event {
	var got []Event
	for event := range events {
		got = append(got, event)
	}
	return got
}

type closeNotifyingReadCloser struct {
	io.ReadCloser
	closed chan<- struct{}
}

func (c closeNotifyingReadCloser) Close() error {
	err := c.ReadCloser.Close()
	select {
	case c.closed <- struct{}{}:
	default:
	}
	return err
}
