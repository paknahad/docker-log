package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/paknahad/docker-log/internal/stream"
)

func TestLogModelFiltersBufferedLinesLive(t *testing.T) {
	model := NewLogModel(nil)
	model = updateLogWithEvent(t, model, stream.Event{Line: "api: started"})
	model = updateLogWithEvent(t, model, stream.Event{Line: "worker: ready"})

	view := model.View()
	if !strings.Contains(view, "api: started") || !strings.Contains(view, "worker: ready") {
		t.Fatalf("View() = %q, want all buffered lines before filtering", view)
	}

	model, cmd := updateLogWithKey(t, model, "api")
	if cmd != nil {
		t.Fatal("typing a filter returned a command, want nil so streams are not restarted")
	}

	view = model.View()
	if !strings.Contains(view, "api: started") {
		t.Fatalf("View() = %q, want matching line", view)
	}
	if strings.Contains(view, "worker: ready") {
		t.Fatalf("View() = %q, want non-matching line hidden", view)
	}
}

func TestLogModelClearingFilterRestoresBufferedLines(t *testing.T) {
	model := NewLogModel(nil)
	model = updateLogWithEvent(t, model, stream.Event{Line: "api: started"})
	model = updateLogWithEvent(t, model, stream.Event{Line: "worker: ready"})
	model, _ = updateLogWithKey(t, model, "api")

	for range "api" {
		model, _ = updateLogWithBackspace(t, model)
	}

	view := model.View()
	if !strings.Contains(view, "api: started") || !strings.Contains(view, "worker: ready") {
		t.Fatalf("View() = %q, want clearing filter to restore buffered lines", view)
	}
}

func TestLogModelRendersStreamErrorsAsBufferedLines(t *testing.T) {
	model := NewLogModel(nil)
	model = updateLogWithEvent(t, model, stream.Event{Container: "api", Err: errForTest("stream closed")})

	view := model.View()
	if !strings.Contains(view, "api: stream closed") {
		t.Fatalf("View() = %q, want stream error line", view)
	}
}

func updateLogWithEvent(t *testing.T, model LogModel, event stream.Event) LogModel {
	t.Helper()

	next, _ := model.Update(streamEventMsg(event))
	logModel, ok := next.(LogModel)
	if !ok {
		t.Fatalf("Update() returned %T, want LogModel", next)
	}
	return logModel
}

func updateLogWithKey(t *testing.T, model LogModel, text string) (LogModel, tea.Cmd) {
	t.Helper()

	next, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(text)})
	logModel, ok := next.(LogModel)
	if !ok {
		t.Fatalf("Update() returned %T, want LogModel", next)
	}
	return logModel, cmd
}

func updateLogWithBackspace(t *testing.T, model LogModel) (LogModel, tea.Cmd) {
	t.Helper()

	next, cmd := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	logModel, ok := next.(LogModel)
	if !ok {
		t.Fatalf("Update() returned %T, want LogModel", next)
	}
	return logModel, cmd
}

type errForTest string

func (e errForTest) Error() string {
	return string(e)
}
