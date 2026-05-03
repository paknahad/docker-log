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

func TestLogModelRendersStreamDisconnectsAsBufferedLines(t *testing.T) {
	model := NewLogModel(nil)
	model = updateLogWithEvent(t, model, stream.Event{Container: "api", Disconnected: true})

	view := model.View()
	if !strings.Contains(view, "api: stream disconnected") {
		t.Fatalf("View() = %q, want stream disconnect line", view)
	}
}

func TestLogModelColorizesOnlyContainerPrefixes(t *testing.T) {
	model := NewLogModel(nil)
	model.colorizePrefixes = true
	model = updateLogWithEvent(t, model, stream.Event{Container: "api", Message: "started", Line: "api: started"})

	view := model.View()
	if !strings.Contains(view, "\x1b[32mapi\x1b[0m: started") {
		t.Fatalf("View() = %q, want colorized container prefix only", view)
	}
	if strings.Contains(view, "\x1b[32mstarted") || strings.Contains(view, "started\x1b[0m") {
		t.Fatalf("View() = %q, want message content left uncolored", view)
	}
}

func TestLogModelKeepsContainerColorsStableDuringSession(t *testing.T) {
	model := NewLogModel(nil)
	model.colorizePrefixes = true
	model = updateLogWithEvent(t, model, stream.Event{Container: "api", Message: "one", Line: "api: one"})
	model = updateLogWithEvent(t, model, stream.Event{Container: "worker", Message: "two", Line: "worker: two"})
	model = updateLogWithEvent(t, model, stream.Event{Container: "api", Message: "three", Line: "api: three"})

	view := model.View()
	if !strings.Contains(view, "\x1b[32mapi\x1b[0m: one") || !strings.Contains(view, "\x1b[32mapi\x1b[0m: three") {
		t.Fatalf("View() = %q, want api to keep its assigned color", view)
	}
	if !strings.Contains(view, "\x1b[33mworker\x1b[0m: two") {
		t.Fatalf("View() = %q, want worker to receive a distinct color", view)
	}
}

func TestLogModelFallsBackToPlainPrefixesWhenColorDisabled(t *testing.T) {
	model := NewLogModel(nil)
	model.colorizePrefixes = false
	model = updateLogWithEvent(t, model, stream.Event{Container: "api", Message: "started", Line: "api: started"})

	view := model.View()
	if strings.Contains(view, "\x1b[") {
		t.Fatalf("View() = %q, want no ANSI escape sequences when color is disabled", view)
	}
	if !strings.Contains(view, "api: started") {
		t.Fatalf("View() = %q, want plain prefixed line", view)
	}
}

func TestLogModelFiltersColorizedLinesUsingPlainText(t *testing.T) {
	model := NewLogModel(nil)
	model.colorizePrefixes = true
	model = updateLogWithEvent(t, model, stream.Event{Container: "api", Message: "started", Line: "api: started"})
	model = updateLogWithEvent(t, model, stream.Event{Container: "worker", Message: "ready", Line: "worker: ready"})
	model, _ = updateLogWithKey(t, model, "worker")

	view := model.View()
	if strings.Contains(view, "api") {
		t.Fatalf("View() = %q, want non-matching line hidden", view)
	}
	if !strings.Contains(view, "\x1b[33mworker\x1b[0m: ready") {
		t.Fatalf("View() = %q, want matching colorized line visible", view)
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
