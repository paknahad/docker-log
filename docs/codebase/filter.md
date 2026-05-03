# Filter Module

## What it does

Provides display-only log line matching. The filter engine decides which buffered log lines should be visible for a query without changing stream state or Docker readers. Filter state is represented in this module so callers do not need to duplicate matching-mode details.

## Public API

- `State`: stores filter text, whether regex mode is enabled, and whether matching is case-sensitive.
- `NewState(text string)`: creates default state for plain-text, case-sensitive matching.
- `NewMatcher(query string)`: creates a default plain-text, case-sensitive matcher.
- `NewMatcherForState(state State)`: creates a matcher for plain-text or regex state and returns validation errors for invalid regex patterns.
- `Matcher.Matches(line string)`: reports whether a line should be visible.
- `Lines(lines []string, query string)`: returns matching lines in their original order.
- `LinesWithState(lines []string, state State)`: returns matching lines or a validation error for invalid state.
- `ValidationError`: reports invalid filter input, currently invalid regex patterns.

## Data tables

None.

## Pipeline steps

The UI keeps the full buffered log output and calls this module during rendering. An empty query matches everything. Default state uses Go's case-sensitive substring matching to preserve the original behavior. Case-insensitive plain-text matching compares lowercased text, while regex matching compiles the pattern once when the matcher is created.

## Routes

None.

## Configuration

None.

## Notes

Filtering must remain downstream from live stream readers. Changing the query should never restart Docker streams or discard buffered lines.
