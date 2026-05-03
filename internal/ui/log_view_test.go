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

func TestLogModelPlainTextFilterDefaultsToCaseSensitive(t *testing.T) {
	model := NewLogModel(nil)
	model = updateLogWithEvent(t, model, stream.Event{Container: "api", Message: "Error opening stream", Line: "api: Error opening stream"})
	model = updateLogWithEvent(t, model, stream.Event{Container: "worker", Message: "error opening stream", Line: "worker: error opening stream"})
	model, _ = updateLogWithKey(t, model, "Error")

	view := model.View()
	if !strings.Contains(view, "api: Error opening stream") {
		t.Fatalf("View() = %q, want exact-case match visible", view)
	}
	if strings.Contains(view, "worker: error opening stream") {
		t.Fatalf("View() = %q, want different-case match hidden by default", view)
	}
	if !strings.Contains(view, "Filter: Error") {
		t.Fatalf("View() = %q, want default case-sensitive prompt", view)
	}
}

func TestLogModelCtrlTTogglesCaseInsensitivePlainTextFiltering(t *testing.T) {
	model := NewLogModel(nil)
	model = updateLogWithEvent(t, model, stream.Event{Container: "api", Message: "Error opening stream", Line: "api: Error opening stream"})
	model = updateLogWithEvent(t, model, stream.Event{Container: "worker", Message: "error opening stream", Line: "worker: error opening stream"})
	model, _ = updateLogWithKey(t, model, "Error")
	model, cmd := updateLogWithCtrlT(t, model)

	if cmd != nil {
		t.Fatal("ctrl+t returned a command, want nil so streams are not restarted")
	}

	view := model.View()
	if !strings.Contains(view, "api: Error opening stream") || !strings.Contains(view, "worker: error opening stream") {
		t.Fatalf("View() = %q, want both case variants visible", view)
	}
	if !strings.Contains(view, "Filter (case-insensitive): Error") {
		t.Fatalf("View() = %q, want case-insensitive prompt", view)
	}
}

func TestLogModelShowsFilterOptionSwitches(t *testing.T) {
	model := NewLogModel(nil)
	model, _ = updateLogWithKey(t, model, "Error")

	view := model.View()
	if !strings.Contains(view, "Options: Regex [off] Ctrl+R | Case-sensitive [on] Ctrl+T") {
		t.Fatalf("View() = %q, want visible default filter option switches", view)
	}

	model, _ = updateLogWithCtrlR(t, model)
	model, _ = updateLogWithCtrlT(t, model)

	view = model.View()
	if !strings.Contains(view, "Options: Regex [on] Ctrl+R | Case-sensitive [off] Ctrl+T") {
		t.Fatalf("View() = %q, want visible toggled filter option switches", view)
	}
	if model.Filter() != "Error" {
		t.Fatalf("Filter() = %q, want toggle shortcuts to leave filter input intact", model.Filter())
	}
}

func TestLogModelTreatsQAsFilterText(t *testing.T) {
	model := NewLogModel(nil)

	model, cmd := updateLogWithKey(t, model, "q")

	if cmd != nil {
		t.Fatal("typing q returned a command, want nil")
	}
	if model.Done() {
		t.Fatal("Done() = true, want false")
	}
	if model.Filter() != "q" {
		t.Fatalf("Filter() = %q, want q", model.Filter())
	}
}

func TestLogModelCtrlCQuits(t *testing.T) {
	model := NewLogModel(nil)

	next, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	logModel, ok := next.(LogModel)
	if !ok {
		t.Fatalf("Update() returned %T, want LogModel", next)
	}
	if !logModel.Done() {
		t.Fatal("Done() = false, want true")
	}
	if cmd == nil {
		t.Fatal("Update(ctrl+c) returned nil command, want quit command")
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
	model, _ = updateLogWithKey(t, model, "ready")

	view := model.View()
	if strings.Contains(view, "api") {
		t.Fatalf("View() = %q, want non-matching line hidden", view)
	}
	if !strings.Contains(view, "\x1b[33mworker\x1b[0m: ready") {
		t.Fatalf("View() = %q, want matching colorized line visible", view)
	}
}

func TestLogModelFilterIgnoresContainerNamePrefix(t *testing.T) {
	model := NewLogModel(nil)
	model = updateLogWithEvent(t, model, stream.Event{Container: "api", Message: "started", Line: "api: started"})
	model, _ = updateLogWithKey(t, model, "api")

	view := model.View()
	if strings.Contains(view, "api: started") {
		t.Fatalf("View() = %q, want container-name-only match hidden", view)
	}
	if !strings.Contains(view, "No log lines match the current filter.") {
		t.Fatalf("View() = %q, want empty filtered state", view)
	}
}

func TestLogModelFiltersMultiplexedStreamsByMessageOnly(t *testing.T) {
	model := NewLogModel(nil)
	model = updateLogWithEvent(t, model, stream.Event{Container: "api", Message: "ready", Line: "api: ready"})
	model = updateLogWithEvent(t, model, stream.Event{Container: "worker", Message: "api request handled", Line: "worker: api request handled"})
	model = updateLogWithEvent(t, model, stream.Event{Container: "api", Message: "idle", Line: "api: idle"})
	model, _ = updateLogWithKey(t, model, "api")

	view := model.View()
	if strings.Contains(view, "api: ready") || strings.Contains(view, "api: idle") {
		t.Fatalf("View() = %q, want api container lines without message matches hidden", view)
	}
	if !strings.Contains(view, "worker: api request handled") {
		t.Fatalf("View() = %q, want message match from worker stream visible", view)
	}
}

func TestLogModelRegexModeMatchesMessageContent(t *testing.T) {
	model := NewLogModel(nil)
	model = updateLogWithEvent(t, model, stream.Event{Container: "api", Message: "status=200", Line: "api: status=200"})
	model = updateLogWithEvent(t, model, stream.Event{Container: "worker", Message: "status=500", Line: "worker: status=500"})
	model, _ = updateLogWithCtrlR(t, model)
	model, _ = updateLogWithKey(t, model, `status=5\d\d`)

	if !model.Regex() {
		t.Fatal("Regex() = false, want true")
	}
	if err := model.FilterError(); err != nil {
		t.Fatalf("FilterError() = %v, want nil", err)
	}

	view := model.View()
	if strings.Contains(view, "api: status=200") {
		t.Fatalf("View() = %q, want regex non-match hidden", view)
	}
	if !strings.Contains(view, "worker: status=500") {
		t.Fatalf("View() = %q, want regex match visible", view)
	}
}

func TestLogModelRegexModeDoesNotMatchContainerNames(t *testing.T) {
	model := NewLogModel(nil)
	model = updateLogWithEvent(t, model, stream.Event{Container: "api-1", Message: "started", Line: "api-1: started"})
	model, _ = updateLogWithCtrlR(t, model)
	model, _ = updateLogWithKey(t, model, `api-\d`)

	view := model.View()
	if strings.Contains(view, "api-1: started") {
		t.Fatalf("View() = %q, want container-name-only regex match hidden", view)
	}
	if !strings.Contains(view, "No log lines match the current filter.") {
		t.Fatalf("View() = %q, want empty filtered state", view)
	}
}

func TestLogModelInvalidRegexIsExposedWithoutCrashing(t *testing.T) {
	model := NewLogModel(nil)
	model = updateLogWithEvent(t, model, stream.Event{Container: "api", Message: "started", Line: "api: started"})
	model, _ = updateLogWithCtrlR(t, model)
	model, _ = updateLogWithKey(t, model, `[`)

	if err := model.FilterError(); err == nil {
		t.Fatal("FilterError() = nil, want invalid regex error")
	}

	view := model.View()
	if !strings.Contains(view, "Invalid regex:") {
		t.Fatalf("View() = %q, want invalid regex feedback", view)
	}
	if strings.Contains(view, "api: started") {
		t.Fatalf("View() = %q, want invalid regex to hide filtered lines", view)
	}
}

func TestLogModelRegexFiltersMultiplexedStreamsByMessageOnly(t *testing.T) {
	model := NewLogModel(nil)
	model = updateLogWithEvent(t, model, stream.Event{Container: "api", Message: "ready", Line: "api: ready"})
	model = updateLogWithEvent(t, model, stream.Event{Container: "worker", Message: "api request 42 handled", Line: "worker: api request 42 handled"})
	model = updateLogWithEvent(t, model, stream.Event{Container: "api", Message: "idle", Line: "api: idle"})
	model, _ = updateLogWithCtrlR(t, model)
	model, _ = updateLogWithKey(t, model, `api request \d+`)

	view := model.View()
	if strings.Contains(view, "api: ready") || strings.Contains(view, "api: idle") {
		t.Fatalf("View() = %q, want api container lines without message regex matches hidden", view)
	}
	if !strings.Contains(view, "worker: api request 42 handled") {
		t.Fatalf("View() = %q, want message regex match from worker stream visible", view)
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

func updateLogWithCtrlR(t *testing.T, model LogModel) (LogModel, tea.Cmd) {
	t.Helper()

	next, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlR})
	logModel, ok := next.(LogModel)
	if !ok {
		t.Fatalf("Update() returned %T, want LogModel", next)
	}
	return logModel, cmd
}

func updateLogWithCtrlT(t *testing.T, model LogModel) (LogModel, tea.Cmd) {
	t.Helper()

	next, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlT})
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
