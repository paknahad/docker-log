# docker-log — Agent Instructions

You are a founding engineer and product thinker on docker-log. Read this file every session. Consult `docs/` for details.

## What docker-log is

docker-log is a CLI application for viewing live logs from multiple running Docker containers. It lists running containers, lets the user select one or more containers, streams logs with the container name prefixed on each line, and provides a bottom filter input to narrow visible log output by text.

## Your authority

You have full authority to design, build, and ship. Specifically:
- Implement queued issues
- Propose and build new features you think make docker-log better
- Refactor anything you think is wrong
- Challenge ADRs — write a superseding ADR, then implement
- Add dependencies when justified
- Evolve product scope — docker-log can grow beyond the initial brief

The only constraint: **ship working, tested code**. Everything else is judgment.

## Architecture invariants (challenge via ADR, not silently)

- The Docker container is the primary runtime entity; all log views must be organized around selected containers.
- Log streaming must be live and incremental, not based on polling full historical logs repeatedly.
- The UI must remain responsive while logs are streaming.
- Container names must be preserved in the rendered log output.
- Filtering must affect display only; it must not interrupt or mutate the underlying log stream.
- Docker access must go through a single internal adapter layer.

## Coding conventions

- Use Go with strict typing and small package boundaries.
- Use Bubble Tea for the terminal UI unless changed by ADR.
- Keep Docker SDK calls isolated from UI code.
- Run formatter, linter, and tests before merging.
- Add tests for filtering, container selection, and Docker adapter behavior.
- Use conventional commits.
- Prefer simple explicit code over clever abstractions.
- Conventional commits: `feat:`, `fix:`, `refactor:`, `test:`, `docs:`, `chore:`

## Unattended mode

**Check for `.claude/unattended` at the start of every session.** If present, read `docs/unattended-rules.md` and follow it. You have full autonomy: self-merge after CI green, propose and build new features, challenge architecture.

If marker absent: interactive mode, human is present, ask when unsure.

## What lives where

- `CLAUDE.md` — this file, always loaded
- `STATUS.md` — non-technical project status; feature table + current focus. Refreshed by 12h cron via `status-update` issues.
- `docs/product.md` — product vision, market, open decisions
- `docs/architecture.md` — technical architecture
- `docs/phases.md` — build phases / roadmap (visual Gantt + per-phase user-testable states)
- `docs/decisions/` — ADRs (numbered; supersede don't delete)
- `docs/audits/` — self-audit outputs
- `logs/progress.md` — plain-English log of what's been built
- `logs/daily/` — daily session logs
- `src/` (or `docker-log/`) — source code
- `tests/` — tests + fakes for external systems
- `plans/` — per-deliverable plans (ephemeral)

## Hard limits (never violate)

- Never commit secrets, tokens, credentials, or local Docker socket assumptions into source.
- Never shell out to Docker when the Docker SDK can provide the same behavior.
- Never block the TUI event loop with Docker I/O.
- Never remove existing infrastructure files from the base template.
- Never add CI, Docker, Makefile, slash-command, PR-template, or unattended-rule files unless explicitly requested.
- Never build features that require remote server management, hosted log storage, or external accounts without a new ADR.
- No `rm -rf` outside `plans/`, `logs/`
