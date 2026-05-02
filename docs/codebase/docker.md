# Docker Module

## What it does

Owns all access to the local Docker daemon. The adapter discovers running containers through the Docker Go SDK, converts SDK response types into shared domain values, and opens normalized live log streams for selected containers.

## Public API

- `NewClient()`: creates a Docker SDK-backed client from the local Docker environment.
- `NewClientWithAPI(api containerAPI)`: creates a client around a test or alternate container-list implementation.
- `Client.ListRunningContainers(ctx context.Context)`: returns normalized running containers or a wrapped discovery error.
- `Client.OpenContainerLogs(ctx context.Context, container domain.Container)`: inspects a selected container, opens a live stdout/stderr log stream, and returns plain log payload bytes or a wrapped Docker error.

## Data tables

None.

## Pipeline steps

The CLI constructs a Docker client, lists running containers, then passes normalized `domain.Container` values into the UI selection model. Streaming code should call `OpenContainerLogs` and pass the returned reader to the stream module instead of calling the Docker SDK directly from UI or stream packages. The returned reader is normalized at the adapter boundary: non-TTY Docker stdout/stderr frames are demultiplexed, while TTY streams are passed through unchanged.

## Routes

None.

## Configuration

The Docker SDK reads standard Docker environment variables such as `DOCKER_HOST` through `client.FromEnv`.

## Notes

Keep Docker SDK types inside this module. Tests should use the small internal interface rather than a live Docker daemon. Log streams are opened with follow enabled, stdout/stderr included, and no historical tail so downstream readers receive live incremental output. Downstream stream/UI code should treat Docker log readers as plain newline-delimited bytes and should not handle Docker multiplex headers itself.
