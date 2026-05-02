# Architecture

## Core architectural decision

The system is built around a multiplexed log stream per selected container.

Each selected Docker container owns an independent live log stream managed by the runtime layer. These streams are normalized into a shared event pipeline that the UI consumes concurrently. Filtering and rendering operate on the multiplexed stream without mutating the underlying sources.

This abstraction separates:
- Docker integration
- Stream management
- Filtering
- Terminal rendering

Changes to this model require an ADR.

## Stack

| Area | Technology | Notes |
|---|---|---|
| Language | Go | Primary implementation language |
| TUI Framework | Bubble Tea | Interactive terminal UI |
| Styling | Lip Gloss | Terminal styling/layout |
| Docker Integration | Docker Go SDK | Direct Docker API access |
| Testing | Go test | Unit and integration tests |
| Linting | golangci-lint | Enforced in CI |
| Formatting | gofmt | Required before merge |

## Module structure

### `/cmd/docker-log`
Application entrypoint and CLI bootstrap.

### `/internal/docker`
Docker SDK integration, container discovery, and log stream creation.

### `/internal/stream`
Stream lifecycle management, fan-in multiplexing, buffering, and backpressure handling.

### `/internal/filter`
Log filtering engine and match logic.

### `/internal/ui`
Bubble Tea models, views, update loop, keyboard handling, and rendering.

### `/internal/domain`
Shared domain models such as containers, log events, and stream state.

### `/test`
Integration and end-to-end tests.

## Data flow

1. User launches the CLI.
2. Docker adapter queries running containers.
3. User selects one or more containers.
4. Runtime creates one live log stream per container.
5. Stream manager multiplexes events into a unified pipeline.
6. Filter engine applies local display filtering.
7. UI renders filtered live output with container prefixes.
8. User updates filters or selection interactively.

## Security model

- No external network services.
- No cloud dependencies.
- Docker access occurs through the local Docker daemon.
- No credentials or secrets stored in code or config.
- No telemetry collection by default.
- All filtering and processing occur locally on-device.
- The application may require Docker socket access depending on host configuration.

## Deployment

The application runs on local developer machines only.

Supported targets:
- Linux
- macOS

Distribution format:
- Single compiled CLI binary

Update strategy:
- Manual binary replacement initially
- Automated release distribution may be added later

No hosted infrastructure is required.