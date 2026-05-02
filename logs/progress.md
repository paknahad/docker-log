## 2026-05-02 — Go project bootstrap opened for review

**What it does:** Prepares docker-log to build as a Go command-line app.
**How:** Adds the Go module, a minimal CLI entrypoint, a shared container domain type, Go Makefile checks, and project-specific README content.
**Why:** Gives the queued feature work a tested project foundation to build on.
**Status:** PR open; human merge required because the bootstrap changes CI and Makefile controls.

## In progress — Multi-container selection UI

**What it does:** Adds the terminal screen where users choose which running containers they want to inspect.
**How:** Builds and tests a Bubble Tea selection model that supports keyboard navigation and multi-select state.
**Why:** Container selection is the first interactive step before docker-log can stream logs from multiple containers.
**Status:** In progress for issue #3.
