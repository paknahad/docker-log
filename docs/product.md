# Product

## One-line pitch

docker-log is a terminal UI for viewing and filtering live logs from multiple Docker containers in one place.

## Vision

Developers running local Docker environments often need to inspect logs across multiple services simultaneously. Existing workflows usually involve multiple terminal windows, repeated `docker logs` commands, or heavy observability stacks that are excessive for local development.

docker-log provides a focused local-first workflow: select running containers, stream logs in real time, and filter output interactively without leaving the terminal. The goal is to reduce friction during debugging and service orchestration work.

The project is intentionally narrow in scope. It is not a hosted logging platform, observability system, or infrastructure management tool. It exists to improve the local Docker developer experience with a fast and dependable CLI workflow.

## Problems it solves

- Watching logs from multiple containers in a single interface
- Quickly identifying which container emitted a log line
- Filtering noisy log output during debugging
- Reducing terminal-window sprawl in Docker-based development
- Providing a lightweight alternative to full observability stacks for local use

## Target users

- Backend developers using Docker Compose locally
- Full-stack developers running multi-service environments
- Platform engineers debugging local container stacks
- Developers working in terminal-first workflows

## Business model

Free and open-source software.

No hosted service, subscription, or paid tier is planned.

## Open decisions

- Whether to support log persistence/export
- Whether to support regex filtering or plain-text only
- Whether to support saved container groups/profiles
- Whether to support remote Docker daemons
- Whether to support Kubernetes in a separate project or not at all

## Out of scope

- Kubernetes log aggregation
- Hosted log storage
- Cloud observability integrations
- Metrics dashboards
- Authentication and user management
- Multi-user collaboration
- Remote infrastructure orchestration
- Browser-based UI
- AI-assisted log analysis