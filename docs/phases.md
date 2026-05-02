# Build phases

## Sequencing rule

Container discovery and selection must land before log streaming, filtering, or UI polish.

## Phase 1 — Container discovery and selection

**Goal:** Show running Docker containers and allow selecting one or more.

**Done when:**
- The CLI starts successfully on Linux and macOS.
- Running containers are listed with name, ID, and image.
- The user can select multiple containers.
- Empty Docker states are handled clearly.
- Container listing has unit or integration test coverage.

## Phase 2 — Live log streaming

**Goal:** Stream logs from selected containers.

**Done when:**
- Each selected container has an independent live log stream.
- Log lines are prefixed with the container name.
- Streaming continues until the user exits.
- Container stop/removal during streaming does not crash the app.
- Stream fan-in behavior is covered by tests.

## Phase 3 — Interactive filtering

**Goal:** Add a bottom filter input that filters visible log output.

**Done when:**
- The user can type and edit a plain-text filter.
- Matching is applied to visible log lines.
- Filtering does not restart or interrupt log streams.
- Clearing the filter restores visible buffered output.
- Filter behavior is covered by tests.

## Phase 4 — Performance and reliability

**Goal:** Keep the TUI responsive under high log volume.

**Done when:**
- The app handles thousands of log lines per second without blocking input.
- Bounded buffering prevents unbounded memory growth.
- Backpressure behavior is explicit and documented.
- Slow containers or noisy containers do not freeze the UI.
- Performance-sensitive stream logic has tests or benchmarks.

## Phase 5 — Release readiness

**Goal:** Prepare a usable open-source v1.

**Done when:**
- Installation and usage docs are complete.
- Linux and macOS builds are produced.
- Basic troubleshooting docs exist.
- CI passes formatter, linter, and tests.
- The README describes scope, limitations, and examples.