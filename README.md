# docker-log

docker-log is a terminal UI for viewing and filtering live logs from multiple running Docker containers.

It is built for local Docker development, where debugging usually means juggling several `docker logs` commands or terminal panes. docker-log will let you select running containers, stream their logs together, preserve each container name on every line, and narrow the visible output with an interactive filter.

The project is intentionally local-first. It is not a hosted logging platform, metrics system, cloud integration, or Kubernetes tool. Docker access goes through the Docker Go SDK, and filtering is display-only so the underlying live streams keep running.

## Tech stack

- Go 1.22+
- Bubble Tea for the terminal UI
- Lip Gloss for terminal styling
- Docker Go SDK for Docker daemon access
- `go test`, `go vet`, `staticcheck`, and `gofmt` for verification

## Quick start

```bash
git clone https://github.com/paknahad/docker-log.git
cd docker-log
make build
make ci
go run ./cmd/docker-log --help
```

The current bootstrap includes a compiling CLI skeleton. The queued implementation issues add Docker discovery, multi-container selection, multiplexed log streaming, and interactive filtering.

## Agent Workflow

This repository is maintained by an unattended coding agent. Project-specific operating instructions live in `CLAUDE.md`; the autonomous work loop and merge rules live in `docs/unattended-rules.md`. For the broader setup guide, see `GETTING_STARTED.md`.

## License

See `LICENSE`.
