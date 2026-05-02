## In progress

**Issue #4 — Multiplexed live log streaming:** Building the focused stream lifecycle and fan-in layer in `internal/stream` on branch `agent/4-multiplexed-live-log-streaming`.

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
