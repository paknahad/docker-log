# Stream Module

## What it does

Owns live log stream lifecycle management. The stream manager starts one reader per selected container source, scans log lines incrementally with a bounded 1 MiB per-line maximum, prefixes each emitted line with the container name, and fans data or source errors into one shared event channel.

## Public API

- `Source`: describes one selected container stream and its reader opener.
- `Event`: normalized stream output with container name, raw message, prefixed line, or an error.
- `NewManager(buffer int)`: creates a stream manager with the requested event channel buffer.
- `Manager.Start(ctx, sources)`: starts all sources concurrently and returns the fan-in event channel.
- `SourcesForContainers(containers, open)`: converts selected domain containers into one stream source each.

## Data tables

None.

## Pipeline steps

The UI or runtime layer supplies selected containers and an opener function from the Docker adapter. The manager opens every source concurrently, scans lines as they arrive, emits prefixed events into a shared channel, and closes the channel after all streams stop. Context cancellation closes live readers so streaming can end on user exit.

## Routes

None.

## Configuration

None.

## Notes

Stream failures are emitted as events for the affected container and do not stop unrelated streams. Individual log lines up to 1 MiB are supported so large structured payloads do not trip the scanner default; lines beyond that bound are treated as stream errors. Filtering remains downstream and must not restart or mutate these streams.
