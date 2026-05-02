package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/paknahad/docker-log/internal/filter"
	"github.com/paknahad/docker-log/internal/stream"
)

// LogModel renders live stream output with a display-only text filter.
type LogModel struct {
	events <-chan stream.Event
	lines  []string
	query  string
	done   bool
}

type streamEventMsg stream.Event
type streamClosedMsg struct{}

func NewLogModel(events <-chan stream.Event) LogModel {
	return LogModel{events: events}
}

func (m LogModel) Init() tea.Cmd {
	return waitForStreamEvent(m.events)
}

func (m LogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			m.done = true
			return m, tea.Quit
		case tea.KeyBackspace:
			if m.query != "" {
				runes := []rune(m.query)
				m.query = string(runes[:len(runes)-1])
			}
			return m, nil
		case tea.KeyEsc:
			m.query = ""
			return m, nil
		}

		switch msg.String() {
		case "q":
			m.done = true
			return m, tea.Quit
		default:
			if len(msg.Runes) > 0 {
				m.query += string(msg.Runes)
			}
			return m, nil
		}
	case streamEventMsg:
		event := stream.Event(msg)
		m.lines = append(m.lines, renderStreamEvent(event))
		return m, waitForStreamEvent(m.events)
	case streamClosedMsg:
		return m, nil
	}

	return m, nil
}

func (m LogModel) View() string {
	var b strings.Builder

	visible := filter.Lines(m.lines, m.query)
	if len(visible) == 0 {
		b.WriteString("No log lines match the current filter.\n")
	} else {
		for _, line := range visible {
			b.WriteString(line)
			b.WriteByte('\n')
		}
	}

	fmt.Fprintf(&b, "\nFilter: %s", m.query)
	return b.String()
}

func (m LogModel) Done() bool {
	return m.done
}

func (m LogModel) Filter() string {
	return m.query
}

func waitForStreamEvent(events <-chan stream.Event) tea.Cmd {
	if events == nil {
		return nil
	}
	return func() tea.Msg {
		event, ok := <-events
		if !ok {
			return streamClosedMsg{}
		}
		return streamEventMsg(event)
	}
}

func renderStreamEvent(event stream.Event) string {
	if event.Err != nil {
		if event.Container == "" {
			return event.Err.Error()
		}
		return fmt.Sprintf("%s: %v", event.Container, event.Err)
	}
	return event.Line
}
