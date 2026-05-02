# UI Module

## What it does

Provides Bubble Tea models for terminal interaction. The selection model lets users navigate running containers, see their image and status, and select one or more containers before log streaming starts. The log model renders buffered stream output and applies a bottom-input text filter without mutating live streams.

## Public API

- `NewSelectionModel(containers []domain.Container)`: creates a container selection model from discovered containers.
- `SelectionModel.Update(msg tea.Msg)`: handles keyboard navigation, multi-select toggling, and quit/start actions.
- `SelectionModel.View()`: renders the selectable container list or an empty state.
- `SelectionModel.SelectedContainers()`: returns selected containers in display order.
- `SelectionModel.Cursor()`: returns the active row index for tests and higher-level orchestration.
- `SelectionModel.Done()`: reports whether the model exited through Enter or a cancel key.
- `SelectionModel.Started()`: reports whether the model exited because the user pressed Enter to start streaming.
- `SelectionModel.Cancelled()`: reports whether the model exited because the user pressed `q` or Ctrl-C.
- `NewLogModel(events <-chan stream.Event)`: creates a log viewer for an existing stream event channel.
- `LogModel.Update(msg tea.Msg)`: consumes stream events, handles filter typing, and exits on quit keys.
- `LogModel.View()`: renders filtered buffered log lines followed by the filter prompt.
- `LogModel.Filter()`: returns the current filter query for tests and orchestration.

## Data tables

None.

## Pipeline steps

The UI receives normalized `domain.Container` values from the Docker adapter layer. It tracks selection state locally and returns the selected containers to later stream-management code. During log viewing, the UI receives normalized `stream.Event` values from the stream module, stores rendered lines in memory, and asks `internal/filter` which buffered lines are visible for the current query.

## Routes

None.

## Configuration

None.

## Notes

Keep Docker SDK access out of this module. UI models should consume domain or stream values and emit user intent so the terminal event loop stays responsive. Filter edits must only affect display state; they should not create new stream commands or reopen Docker readers.
