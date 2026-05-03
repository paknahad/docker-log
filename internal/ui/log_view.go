package ui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/paknahad/docker-log/internal/filter"
	"github.com/paknahad/docker-log/internal/stream"
)

// LogModel renders live stream output with a display-only text filter.
type LogModel struct {
	events           <-chan stream.Event
	lines            []renderedLogLine
	filterState      filter.State
	done             bool
	colorizePrefixes bool
	containerColors  map[string]string
	nextColor        int
}

type streamEventMsg stream.Event
type streamClosedMsg struct{}

type renderedLogLine struct {
	filterText string
	display    string
}

func NewLogModel(events <-chan stream.Event) LogModel {
	return LogModel{
		events:           events,
		filterState:      filter.NewState(""),
		colorizePrefixes: terminalSupportsANSIPrefixColors(),
		containerColors:  make(map[string]string),
	}
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
		case tea.KeyCtrlR:
			m.filterState.Regex = !m.filterState.Regex
			return m, nil
		case tea.KeyCtrlT:
			m.filterState.CaseSensitive = !m.filterState.CaseSensitive
			return m, nil
		case tea.KeyBackspace:
			if m.filterState.Text != "" {
				runes := []rune(m.filterState.Text)
				m.filterState.Text = string(runes[:len(runes)-1])
			}
			return m, nil
		case tea.KeyEsc:
			m.filterState.Text = ""
			return m, nil
		}

		if len(msg.Runes) > 0 {
			m.filterState.Text += string(msg.Runes)
		}
		return m, nil
	case streamEventMsg:
		event := stream.Event(msg)
		m.lines = append(m.lines, m.renderStreamEvent(event))
		return m, waitForStreamEvent(m.events)
	case streamClosedMsg:
		return m, nil
	}

	return m, nil
}

func (m LogModel) View() string {
	var b strings.Builder

	visible := m.visibleLines()
	if len(visible) == 0 {
		b.WriteString("No log lines match the current filter.\n")
	} else {
		for _, line := range visible {
			b.WriteString(line)
			b.WriteByte('\n')
		}
	}

	fmt.Fprintf(&b, "\nFilter%s: %s", m.filterModeLabel(), m.filterState.Text)
	fmt.Fprintf(&b, "\nOptions: Regex [%s] Ctrl+R | Case-sensitive [%s] Ctrl+T", onOff(m.filterState.Regex), onOff(m.filterState.CaseSensitive))
	if err := m.FilterError(); err != nil {
		fmt.Fprintf(&b, "\nInvalid regex: %v", err)
	}
	return b.String()
}

func (m LogModel) Done() bool {
	return m.done
}

func (m LogModel) Filter() string {
	return m.filterState.Text
}

func (m LogModel) Regex() bool {
	return m.filterState.Regex
}

func (m LogModel) FilterError() error {
	_, err := filter.NewMatcherForState(m.filterState)
	return err
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

func (m LogModel) visibleLines() []string {
	matcher, err := filter.NewMatcherForState(m.filterState)
	if err != nil {
		return nil
	}
	visible := make([]string, 0, len(m.lines))
	for _, line := range m.lines {
		if matcher.Matches(line.filterText) {
			visible = append(visible, line.display)
		}
	}
	return visible
}

func (m LogModel) filterModeLabel() string {
	labels := make([]string, 0, 2)
	if m.filterState.Regex {
		labels = append(labels, "regex")
	}
	if !m.filterState.CaseSensitive {
		labels = append(labels, "case-insensitive")
	}
	if len(labels) == 0 {
		return ""
	}
	return " (" + strings.Join(labels, ", ") + ")"
}

func (m *LogModel) renderStreamEvent(event stream.Event) renderedLogLine {
	if event.Err != nil {
		if event.Container == "" {
			return plainLogLine(event.Err.Error())
		}
		return m.renderPrefixedLine(event.Container, event.Err.Error())
	}
	if event.Disconnected {
		if event.Container == "" {
			return plainLogLine("stream disconnected")
		}
		return m.renderPrefixedLine(event.Container, "stream disconnected")
	}
	if event.Container == "" {
		return plainLogLine(event.Line)
	}
	message := event.Message
	if message == "" {
		message = strings.TrimPrefix(event.Line, event.Container+": ")
	}
	return m.renderPrefixedLine(event.Container, message, message)
}

func (m *LogModel) renderPrefixedLine(container, message string, filterText ...string) renderedLogLine {
	textForFilter := message
	if len(filterText) > 0 {
		textForFilter = filterText[0]
	}
	display := fmt.Sprintf("%s: %s", container, message)
	if !m.colorizePrefixes {
		return renderedLogLine{filterText: textForFilter, display: display}
	}

	color := m.colorForContainer(container)
	return renderedLogLine{
		filterText: textForFilter,
		display:    fmt.Sprintf("\x1b[%sm%s\x1b[0m: %s", color, container, message),
	}
}

func (m *LogModel) colorForContainer(container string) string {
	if m.containerColors == nil {
		m.containerColors = make(map[string]string)
	}
	if color, ok := m.containerColors[container]; ok {
		return color
	}

	color := prefixColorPalette[m.nextColor%len(prefixColorPalette)]
	m.containerColors[container] = color
	m.nextColor++
	return color
}

func plainLogLine(line string) renderedLogLine {
	return renderedLogLine{filterText: line, display: line}
}

func onOff(enabled bool) string {
	if enabled {
		return "on"
	}
	return "off"
}

func terminalSupportsANSIPrefixColors() bool {
	if _, disabled := os.LookupEnv("NO_COLOR"); disabled {
		return false
	}
	if os.Getenv("CLICOLOR") == "0" {
		return false
	}
	term := os.Getenv("TERM")
	return term != "" && term != "dumb"
}

var prefixColorPalette = []string{"32", "33", "36", "35", "34", "31", "92", "94"}
