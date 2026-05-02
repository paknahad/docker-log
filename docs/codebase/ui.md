# UI Module

## What it does

Provides Bubble Tea models for terminal interaction. The current model lets users navigate a list of running containers and select one or more containers before log streaming starts.

## Public API

- `NewSelectionModel(containers []domain.Container)`: creates a container selection model from discovered containers.
- `SelectionModel.Update(msg tea.Msg)`: handles keyboard navigation, multi-select toggling, and quit/start actions.
- `SelectionModel.View()`: renders the selectable container list or an empty state.
- `SelectionModel.SelectedContainers()`: returns selected containers in display order.
- `SelectionModel.Cursor()`: returns the active row index for tests and higher-level orchestration.
- `SelectionModel.Done()`: reports whether the model exited through enter or quit.

## Data tables

None.

## Pipeline steps

The UI receives normalized `domain.Container` values from the Docker adapter layer. It tracks selection state locally and returns the selected containers to later stream-management code. Filtering and streaming are handled by separate modules.

## Routes

None.

## Configuration

None.

## Notes

Keep Docker SDK access out of this module. UI models should consume domain values and emit user intent so the terminal event loop stays responsive.
