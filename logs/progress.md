## In progress

**Issue #44 — Add filter option switches to the TUI:** Adding visible regex and case-sensitivity controls to the log view while keeping the filter input usable.

## 2026-05-03 — Filter case sensitivity can be toggled

**What it does:** docker-log now lets users switch the live log filter between exact-case matching and case-insensitive matching.
**How:** The log view keeps the existing case-sensitive default and uses Ctrl+T to toggle the filter state without restarting streams.
**Why:** Developers can quickly find messages when log casing varies across services while preserving precise matching by default.
**Status:** Merged.
**PR:** #46
STATUS: Advanced filtering → ✅ shipped

## 2026-05-03 — Regex filtering is available in the log view

**What it does:** docker-log can now filter live log messages with Go regular expressions and shows invalid regex feedback instead of failing silently.
**How:** The log view toggles regex mode while still filtering only message content, so container names do not create false matches.
**Why:** Developers can narrow noisy multiplexed logs with more precise patterns while keeping the existing plain-text filter behavior.
**Status:** Merged.
**PR:** #45
STATUS: Advanced filtering → ✅ shipped

## 2026-05-03 — Filters now ignore container names

**What it does:** docker-log now filters live logs by the actual message text instead of matching container name prefixes.
**How:** The log view keeps rendered output separate from the text used for filter matching, so container names stay visible without affecting results.
**Why:** Users can filter multiplexed logs by content without accidentally showing every line from a similarly named container.
**Status:** Merged.
**PR:** #39
STATUS: Interactive log filtering -> ✅ shipped

## 2026-05-03 — README reflects the shipped log workflow

**What it does:** The README now tells users that docker-log can list running containers, select one or more, stream live logs, and filter visible output.
**How:** Updated the main project description and quick start without changing the app or broadening its local Docker scope.
**Why:** The project entry point should match the workflow users can test today.
**Status:** PR open.
**PR:** #36

## 2026-05-03 — Ctrl+C is the only exit shortcut

**What it does:** docker-log no longer exits when the user presses `q`, so `q` can be typed into the live log filter like normal text.
**How:** The selection and log views now reserve application exit for Ctrl+C and the help text reflects that shortcut.
**Why:** Filter input should accept ordinary letters without accidentally closing the app.
**Status:** Merged.
**PR:** #35
STATUS: Interactive log filtering -> ✅ shipped

## 2026-05-03 — Container prefixes are colorized

**What it does:** docker-log now gives each container name prefix a stable color in the live log view so multiplexed output is easier to scan.
**How:** The UI assigns prefix colors during a log-view session while keeping message text and filter matching plain.
**Why:** Users can distinguish containers faster without changing the underlying log content.
**Status:** Merged.
**PR:** #34
STATUS: Multiplexed log readability -> ✅ shipped

## 2026-05-03 — Container lifecycle disconnects are visible

**What it does:** docker-log now shows when an individual container log stream disconnects instead of silently ending that container's output.
**How:** The stream manager emits a lifecycle event on clean stream closure and the log view renders that status alongside log lines.
**Why:** Stopped or restarted containers should be visible to the user without crashing the application or hiding the failure mode.
**Status:** Merged.
**PR:** #32
STATUS: Live log streaming resilience -> ✅ shipped

## 2026-05-03 — Stream buffering is bounded

**What it does:** docker-log now has an explicit cap on the stream event queue so heavy log output cannot request an unlimited in-memory buffer.
**How:** The stream manager caps requested event buffers at 4,096 events and blocks producers when the UI falls behind.
**Why:** This keeps memory growth bounded under high log volume while preserving log lines instead of silently dropping them.
**Status:** Merged.
**PR:** #31
STATUS: Live log streaming resilience -> ✅ shipped

## 2026-05-03 — Selection quit exits cleanly

**What it does:** docker-log now exits from the container selection screen when the user presses `q` or Ctrl-C, even if containers were already selected.
**How:** The CLI only starts streams when the selection UI reports that Enter started the session.
**Why:** Cancel keys should never surprise users by opening live log streams.
**Status:** Merged.
**PR:** #30
STATUS: Selection cancel handling -> ✅ shipped

## 2026-05-02 — Docker log frame artifacts removed

**What it does:** docker-log now shows clean stdout and stderr text for common non-TTY containers instead of leaking Docker framing bytes into the log view.
**How:** The Docker adapter inspects TTY mode and demultiplexes Docker stdout/stderr frames before the stream layer reads them.
**Why:** Live logs should look like readable container output no matter how Docker transports stdout and stderr internally.
**Status:** Merged.
**PR:** #27
STATUS: Live log streaming resilience -> ✅ shipped

## 2026-05-03 — Filter modes can be represented safely

**What it does:** docker-log now has the internal filter settings needed for plain text, regex, and case-sensitive matching modes.
**How:** The filter engine stores matching mode in one state value and validates regex patterns before using them.
**Why:** This prepares advanced filtering without spreading matching rules through the UI or risking crashes on bad patterns.
**Status:** Merged.
**PR:** #41
STATUS: Advanced filtering → foundation shipped

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
