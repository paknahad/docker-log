# ADR 0001 — Multiplexed log streams per container

**Status:** Accepted
**Date:** 2026-05-02

## Context

docker-log needs to display live logs from multiple Docker containers simultaneously inside a single terminal UI.

The application must:
- Support concurrent log streams
- Keep the UI responsive under high log volume
- Preserve container identity for each log line
- Allow filtering without interrupting log collection
- Avoid coupling Docker SDK logic directly to the UI

A simple sequential model or direct UI-bound Docker streaming approach would make filtering, buffering, testing, and concurrency management harder as the application grows.

## Decision

The system will use a multiplexed log stream architecture.

Each selected container owns an independent live log stream managed by the runtime layer. These streams emit normalized log events into a shared fan-in pipeline. The UI consumes this multiplexed stream and applies display filtering separately from stream ingestion.

The architecture is divided into:
- Docker adapter layer
- Stream management layer
- Filtering layer
- UI rendering layer

Filtering operates on rendered visibility only and does not mutate or restart upstream streams.

Docker integration must remain isolated behind internal interfaces to allow testing without requiring a live Docker daemon.

## Consequences

**Positive:**
- Clean separation between Docker access and UI rendering
- Easier concurrency management
- Filtering can evolve independently
- Better testability through stream abstraction
- Scales better under high log throughput
- Allows future additions such as buffering or export without redesigning the pipeline

**Negative:**
- More internal complexity than direct stream-to-UI rendering
- Requires careful synchronization and backpressure handling
- Introduces buffering and lifecycle management concerns
- Slightly higher memory usage during heavy log streaming

## Alternatives considered

### Direct Docker stream rendering in UI

Rejected because UI responsiveness and testability would degrade as concurrency increases.

### Polling container logs repeatedly

Rejected because it is inefficient, introduces duplication risk, and does not provide true live streaming behavior.

### Single merged Docker stream without per-container ownership

Rejected because container lifecycle management and source attribution become unclear.

## Enforcement

The following changes require a new ADR:
- Replacing the multiplexed stream model
- Coupling Docker SDK calls directly into UI components
- Moving filtering upstream into stream ingestion
- Replacing container-owned stream lifecycle management