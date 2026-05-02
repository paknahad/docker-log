package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/paknahad/docker-log/internal/domain"
)

// SelectionModel is the Bubble Tea model for choosing containers to stream.
type SelectionModel struct {
	containers []domain.Container
	cursor     int
	selected   map[string]struct{}
	result     selectionResult
}

type selectionResult int

const (
	selectionPending selectionResult = iota
	selectionStarted
	selectionCancelled
)

func NewSelectionModel(containers []domain.Container) SelectionModel {
	return SelectionModel{
		containers: append([]domain.Container(nil), containers...),
		selected:   make(map[string]struct{}),
	}
}

func (m SelectionModel) Init() tea.Cmd {
	return nil
}

func (m SelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch key.String() {
	case "enter":
		m.result = selectionStarted
		return m, tea.Quit
	case "ctrl+c", "q":
		m.result = selectionCancelled
		return m, tea.Quit
	case "up", "k":
		if len(m.containers) == 0 {
			return m, nil
		}
		m.cursor--
		if m.cursor < 0 {
			m.cursor = len(m.containers) - 1
		}
	case "down", "j":
		if len(m.containers) == 0 {
			return m, nil
		}
		m.cursor = (m.cursor + 1) % len(m.containers)
	case " ":
		if len(m.containers) == 0 {
			return m, nil
		}
		id := m.containers[m.cursor].ID
		if _, exists := m.selected[id]; exists {
			delete(m.selected, id)
		} else {
			m.selected[id] = struct{}{}
		}
	}

	return m, nil
}

func (m SelectionModel) View() string {
	if len(m.containers) == 0 {
		return "No running containers found.\n\nPress q to quit.\n"
	}

	var b strings.Builder
	b.WriteString("Select containers to stream\n\n")

	for i, container := range m.containers {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[container.ID]; ok {
			checked = "x"
		}

		fmt.Fprintf(&b, "%s [%s] %s", cursor, checked, container.DisplayName())
		if container.Image != "" {
			fmt.Fprintf(&b, "  %s", container.Image)
		}
		if container.Status != "" {
			fmt.Fprintf(&b, "  %s", container.Status)
		}
		b.WriteByte('\n')
	}

	b.WriteString("\nSpace selects, enter starts, q quits.\n")
	return b.String()
}

func (m SelectionModel) Cursor() int {
	return m.cursor
}

func (m SelectionModel) Done() bool {
	return m.result != selectionPending
}

func (m SelectionModel) Started() bool {
	return m.result == selectionStarted
}

func (m SelectionModel) Cancelled() bool {
	return m.result == selectionCancelled
}

func (m SelectionModel) SelectedContainers() []domain.Container {
	selected := make([]domain.Container, 0, len(m.selected))
	for _, container := range m.containers {
		if _, ok := m.selected[container.ID]; ok {
			selected = append(selected, container)
		}
	}
	return selected
}
