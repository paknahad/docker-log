## 2026-05-02 — Selection start and cancel are distinct

**What it does:** docker-log can now tell whether the selection screen exited because the user started streaming or cancelled.
**How:** The selection UI records a focused result state and exposes helpers for started versus cancelled outcomes.
**Why:** Later orchestration can handle Enter differently from `q` or Ctrl-C without guessing from selected containers.
**Status:** Merged.
**PR:** #26
STATUS: Selection cancel handling -> ✅ shipped

## 2026-05-02 — Long log lines keep streaming

**What it does:** docker-log can keep reading a container when it emits a single unusually large log line.
**How:** The stream reader now allows log lines up to a documented 1 MiB bound instead of stopping at the default scanner limit.
**Why:** Large JSON payloads and stack traces should not make a healthy container disappear from the live log view.
**Status:** Merged.
**PR:** #25
STATUS: Live log streaming resilience -> ✅ shipped

## 2026-05-02 — Selected containers wired into live log view

**What it does:** docker-log can now move from selecting containers into the live log viewer for the chosen containers.
**How:** The CLI reads the selection result, creates Docker-backed stream sources, starts the stream manager, and launches the log view.
**Why:** This completes the core product workflow from choosing containers to watching live logs.
**Status:** PR open.
**PR:** #20
STATUS: End-to-end live log viewing → PR open

## 2026-05-02 — Docker log stream opener added

**What it does:** docker-log can now ask Docker for a live log stream from a selected container through the project’s adapter layer.
**How:** Adds a Docker client method that opens followed stdout/stderr logs without loading historical output, plus adapter-boundary tests.
**Why:** This keeps real log access out of UI/runtime code and clears the next step for wiring selected containers into the log viewer.
**Status:** PR open.
**PR:** #19
STATUS: Docker log stream opener -> foundation shipped

## 2026-05-02 — Interactive log filtering added

**What it does:** Adds the log-view pieces needed for users to type a filter and narrow visible log output while keeping buffered lines available.
**How:** Adds a case-sensitive filter module and a Bubble Tea log model that buffers stream events and filters only what is rendered.
**Why:** This makes filtering a display concern, so live streams can keep running while users refine or clear the filter.
**Status:** Merged.
**PR:** #15
STATUS: Interactive log filtering → foundation shipped

## 2026-05-02 — Multiplexed stream manager added

**What it does:** docker-log now has the internal machinery to read multiple selected container log streams at the same time and combine them into one feed.
**How:** Adds a stream manager that runs one reader per container source, prefixes each log line with the container name, and reports stream-specific errors without stopping other streams.
**Why:** This creates the live streaming foundation needed before the terminal UI can show logs from several containers together.
**Status:** PR open.
**PR:** #12
STATUS: Multiplexed live log streaming → foundation shipped

## 2026-05-02 — Docker container discovery added

**What it does:** docker-log now finds running local Docker containers and shows their name, image, and status in the selection screen.
**How:** Adds a Docker SDK adapter that normalizes container data before passing it to the terminal UI.
**Why:** Real container discovery is the foundation for selecting containers and streaming their logs.
**Status:** Merged.
**PR:** #11
STATUS: Docker container discovery → ✅ shipped

## 2026-05-02 — Go project bootstrap opened for review

**What it does:** Prepares docker-log to build as a Go command-line app.
**How:** Adds the Go module, a minimal CLI entrypoint, a shared container domain type, Go Makefile checks, and project-specific README content.
**Why:** Gives the queued feature work a tested project foundation to build on.
**Status:** PR open; human merge required because the bootstrap changes CI and Makefile controls.

## 2026-05-02 — Multi-container selection UI added

**What it does:** Adds the terminal screen where users choose which running containers they want to inspect.
**How:** Builds and tests a Bubble Tea selection model that supports keyboard navigation and multi-select state.
**Why:** Container selection is the first interactive step before docker-log can stream logs from multiple containers.
**Status:** Merged.
**PR:** #10
STATUS: Multi-container selection UI → ✅ shipped
