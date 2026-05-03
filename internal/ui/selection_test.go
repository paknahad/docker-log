package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/paknahad/docker-log/internal/domain"
)

func TestSelectionModelSelectsMultipleContainers(t *testing.T) {
	model := NewSelectionModel([]domain.Container{
		{ID: "api-id", Name: "api", Image: "example/api:latest"},
		{ID: "db-id", Name: "db", Image: "postgres:16"},
	})

	model = applyKey(t, model, " ")
	model = applyKey(t, model, "down")
	model = applyKey(t, model, " ")

	selected := model.SelectedContainers()
	if len(selected) != 2 {
		t.Fatalf("len(SelectedContainers()) = %d, want 2", len(selected))
	}
	if selected[0].ID != "api-id" || selected[1].ID != "db-id" {
		t.Fatalf("SelectedContainers() = %#v, want api-id then db-id", selected)
	}
}

func TestSelectionModelKeyboardNavigationWraps(t *testing.T) {
	model := NewSelectionModel([]domain.Container{
		{ID: "api-id", Name: "api"},
		{ID: "db-id", Name: "db"},
	})

	model = applyKey(t, model, "up")
	if model.Cursor() != 1 {
		t.Fatalf("Cursor() after up from first row = %d, want 1", model.Cursor())
	}

	model = applyKey(t, model, "down")
	if model.Cursor() != 0 {
		t.Fatalf("Cursor() after down from last row = %d, want 0", model.Cursor())
	}
}

func TestSelectionModelEmptyStateRendersClearly(t *testing.T) {
	model := NewSelectionModel(nil)

	view := model.View()
	if !strings.Contains(view, "No running containers found") {
		t.Fatalf("View() = %q, want clear empty state", view)
	}
}

func TestSelectionModelRendersContainerStatus(t *testing.T) {
	model := NewSelectionModel([]domain.Container{
		{ID: "api-id", Name: "api", Image: "example/api:latest", Status: "Up 2 minutes"},
	})

	view := model.View()
	if !strings.Contains(view, "Up 2 minutes") {
		t.Fatalf("View() = %q, want container status", view)
	}
}

func TestSelectionModelEnterStartsSelection(t *testing.T) {
	model := NewSelectionModel([]domain.Container{{ID: "api-id", Name: "api"}})

	next, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	selection, ok := next.(SelectionModel)
	if !ok {
		t.Fatalf("Update() returned %T, want SelectionModel", next)
	}
	if !selection.Done() {
		t.Fatal("Done() = false, want true")
	}
	if !selection.Started() {
		t.Fatal("Started() = false, want true")
	}
	if selection.Cancelled() {
		t.Fatal("Cancelled() = true, want false")
	}
	if cmd == nil {
		t.Fatal("Update(enter) returned nil command, want quit command")
	}
}

func TestSelectionModelCtrlCCancelsSelection(t *testing.T) {
	model := NewSelectionModel([]domain.Container{{ID: "api-id", Name: "api"}})

	next, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	selection, ok := next.(SelectionModel)
	if !ok {
		t.Fatalf("Update() returned %T, want SelectionModel", next)
	}
	if !selection.Done() {
		t.Fatal("Done() = false, want true")
	}
	if selection.Started() {
		t.Fatal("Started() = true, want false")
	}
	if !selection.Cancelled() {
		t.Fatal("Cancelled() = false, want true")
	}
	if cmd == nil {
		t.Fatal("Update(ctrl+c) returned nil command, want quit command")
	}
}

func TestSelectionModelQDoesNotQuit(t *testing.T) {
	model := NewSelectionModel([]domain.Container{{ID: "api-id", Name: "api"}})

	next, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})

	selection, ok := next.(SelectionModel)
	if !ok {
		t.Fatalf("Update() returned %T, want SelectionModel", next)
	}
	if selection.Done() {
		t.Fatal("Done() = true, want false")
	}
	if selection.Started() {
		t.Fatal("Started() = true, want false")
	}
	if selection.Cancelled() {
		t.Fatal("Cancelled() = true, want false")
	}
	if cmd != nil {
		t.Fatal("Update(q) returned command, want nil")
	}
}

func applyKey(t *testing.T, model SelectionModel, key string) SelectionModel {
	t.Helper()

	next, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
	if key == "down" {
		next, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	}
	if key == "up" {
		next, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	}

	selection, ok := next.(SelectionModel)
	if !ok {
		t.Fatalf("Update() returned %T, want SelectionModel", next)
	}
	return selection
}
