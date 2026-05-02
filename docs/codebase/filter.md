# Filter Module

## What it does

Provides display-only log line matching. The filter engine decides which buffered log lines should be visible for a query without changing stream state or Docker readers.

## Public API

- `NewMatcher(query string)`: creates a case-sensitive substring matcher.
- `Matcher.Matches(line string)`: reports whether a line should be visible.
- `Lines(lines []string, query string)`: returns matching lines in their original order.

## Data tables

None.

## Pipeline steps

The UI keeps the full buffered log output and calls this module during rendering. An empty query matches everything. Non-empty queries use Go's case-sensitive substring matching.

## Routes

None.

## Configuration

None.

## Notes

Filtering must remain downstream from live stream readers. Changing the query should never restart Docker streams or discard buffered lines.
