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

	lineEvents := logLineEvents(got)
	want := map[string]bool{
		"api|api one|api: api one":             false,
		"api|api two|api: api two":             false,
		"worker|worker one|worker: worker one": false,
	}
	for _, event := range lineEvents {
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

func TestManagerReportsCleanStreamDisconnectWithoutStoppingOtherStreams(t *testing.T) {
	sources := []Source{
		{
			Container: "done",
			Open: func(context.Context) (io.ReadCloser, error) {
				return io.NopCloser(strings.NewReader("last line\n")), nil
			},
		},
		{
			Container: "other",
			Open: func(context.Context) (io.ReadCloser, error) {
				return io.NopCloser(strings.NewReader("still visible\n")), nil
			},
		},
	}

	events := NewManager(0).Start(context.Background(), sources)
	got := collectEvents(events)

	var sawDisconnect bool
	var sawOtherLine bool
	for _, event := range got {
		if event.Container == "done" && event.Disconnected {
			sawDisconnect = true
		}
		if event.Container == "other" && event.Line == "other: still visible" {
			sawOtherLine = true
		}
	}
	if !sawDisconnect {
		t.Fatalf("did not receive clean stream disconnect in %#v", got)
	}
	if !sawOtherLine {
		t.Fatalf("unrelated stream did not continue in %#v", got)
	}
}

func TestManagerStreamsLinesLongerThanScannerDefaultLimit(t *testing.T) {
	longLine := strings.Repeat("x", 70*1024)
	source := Source{
		Container: "api",
		Open: func(context.Context) (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader(longLine + "\n")), nil
		},
	}

	events := collectEvents(NewManager(0).Start(context.Background(), []Source{source}))

	lineEvents := logLineEvents(events)
	if len(lineEvents) != 1 {
		t.Fatalf("len(lineEvents) = %d, want 1: %#v", len(lineEvents), events)
	}
	if lineEvents[0].Err != nil {
		t.Fatalf("event.Err = %v, want nil", lineEvents[0].Err)
	}
	if lineEvents[0].Message != longLine {
		t.Fatalf("len(event.Message) = %d, want %d", len(lineEvents[0].Message), len(longLine))
	}
	if lineEvents[0].Line != "api: "+longLine {
		t.Fatalf("event.Line has prefix/message mismatch")
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

func TestManagerCapsRequestedEventBuffer(t *testing.T) {
	events := NewManager(maxEventBuffer+1).Start(context.Background(), nil)

	if got := cap(events); got != maxEventBuffer {
		t.Fatalf("cap(events) = %d, want %d", got, maxEventBuffer)
	}
}

func TestManagerAppliesBackpressureWhenEventBufferIsFull(t *testing.T) {
	reader, writer := io.Pipe()
	source := Source{
		Container: "api",
		Open: func(context.Context) (io.ReadCloser, error) {
			return reader, nil
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	events := NewManager(1).Start(ctx, []Source{source})

	if _, err := writer.Write([]byte("one\n")); err != nil {
		t.Fatalf("write first line: %v", err)
	}

	secondWritten := make(chan error, 1)
	go func() {
		_, err := writer.Write([]byte("two\n"))
		secondWritten <- err
	}()

	select {
	case err := <-secondWritten:
		if err != nil {
			t.Fatalf("write second line: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("second write did not reach the stream reader")
	}

	thirdWritten := make(chan error, 1)
	go func() {
		_, err := writer.Write([]byte("three\n"))
		thirdWritten <- err
	}()

	select {
	case err := <-thirdWritten:
		t.Fatalf("third write completed before consumer drained the full event buffer: %v", err)
	case <-time.After(50 * time.Millisecond):
	}

	if event := <-events; event.Line != "api: one" {
		t.Fatalf("first event = %#v, want api: one", event)
	}

	select {
	case err := <-thirdWritten:
		if err != nil {
			t.Fatalf("write third line after draining buffer: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("third write stayed blocked after consumer drained the event buffer")
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

	events := logLineEvents(collectEvents(NewManager(0).Start(context.Background(), sources)))
	if len(events) != 2 {
		t.Fatalf("len(log line events) = %d, want 2", len(events))
	}
}

func BenchmarkManagerFanInWithBoundedBuffer(b *testing.B) {
	var input strings.Builder
	for i := 0; i < b.N; i++ {
		input.WriteString("line\n")
	}
	source := Source{
		Container: "api",
		Open: func(context.Context) (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader(input.String())), nil
		},
	}

	b.ReportAllocs()
	b.ResetTimer()
	count := 0
	for range NewManager(maxEventBuffer).Start(context.Background(), []Source{source}) {
		count++
	}
	if count != b.N {
		b.Fatalf("count = %d, want %d", count, b.N)
	}
}

func collectEvents(events <-chan Event) []Event {
	var got []Event
	for event := range events {
		got = append(got, event)
	}
	return got
}

func logLineEvents(events []Event) []Event {
	var lines []Event
	for _, event := range events {
		if event.Line != "" {
			lines = append(lines, event)
		}
	}
	return lines
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
