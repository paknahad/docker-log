package stream

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/paknahad/docker-log/internal/domain"
)

type Source struct {
	ContainerID string
	Container   string
	Open        func(context.Context) (io.ReadCloser, error)
}

type Event struct {
	Container string
	Message   string
	Line      string
	Err       error
}

type Manager struct {
	buffer int
}

const (
	maxLogLineBytes = 1024 * 1024
	maxEventBuffer  = 4096
)

func NewManager(buffer int) Manager {
	if buffer < 0 {
		buffer = 0
	}
	if buffer > maxEventBuffer {
		buffer = maxEventBuffer
	}
	return Manager{buffer: buffer}
}

func SourcesForContainers(containers []domain.Container, open func(context.Context, domain.Container) (io.ReadCloser, error)) []Source {
	sources := make([]Source, 0, len(containers))
	for _, container := range containers {
		container := container
		sources = append(sources, Source{
			ContainerID: container.ID,
			Container:   container.DisplayName(),
			Open: func(ctx context.Context) (io.ReadCloser, error) {
				if open == nil {
					return nil, fmt.Errorf("open stream for %s: missing container opener", container.DisplayName())
				}
				return open(ctx, container)
			},
		})
	}
	return sources
}

func (m Manager) Start(ctx context.Context, sources []Source) <-chan Event {
	events := make(chan Event, m.buffer)
	var wg sync.WaitGroup

	for _, source := range sources {
		source := source
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.stream(ctx, source, events)
		}()
	}

	go func() {
		wg.Wait()
		close(events)
	}()

	return events
}

func (m Manager) stream(ctx context.Context, source Source, events chan<- Event) {
	if source.Open == nil {
		send(ctx, events, Event{
			Container: source.Container,
			Err:       fmt.Errorf("open stream for %s: missing source opener", source.Container),
		})
		return
	}

	reader, err := source.Open(ctx)
	if err != nil {
		send(ctx, events, Event{Container: source.Container, Err: err})
		return
	}
	if reader == nil {
		send(ctx, events, Event{
			Container: source.Container,
			Err:       fmt.Errorf("open stream for %s: nil reader", source.Container),
		})
		return
	}

	var closeOnce sync.Once
	closeReader := func() {
		closeOnce.Do(func() {
			_ = reader.Close()
		})
	}
	defer closeReader()

	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			closeReader()
		case <-done:
		}
	}()
	defer close(done)

	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, 64*1024), maxLogLineBytes)
	for scanner.Scan() {
		message := scanner.Text()
		event := Event{
			Container: source.Container,
			Message:   message,
			Line:      fmt.Sprintf("%s: %s", source.Container, message),
		}
		if !send(ctx, events, event) {
			return
		}
	}
	if err := scanner.Err(); err != nil && ctx.Err() == nil {
		send(ctx, events, Event{Container: source.Container, Err: err})
	}
}

func send(ctx context.Context, events chan<- Event, event Event) bool {
	// A full channel blocks producer goroutines on purpose. This applies
	// backpressure to Docker readers instead of dropping lines or allowing an
	// unbounded in-memory queue to grow when the UI cannot keep up.
	select {
	case events <- event:
		return true
	case <-ctx.Done():
		return false
	}
}
