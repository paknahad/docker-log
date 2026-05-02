# Docker Module

## What it does

Owns all access to the local Docker daemon. The current adapter discovers running containers through the Docker Go SDK and converts SDK response types into shared domain values.

## Public API

- `NewClient()`: creates a Docker SDK-backed client from the local Docker environment.
- `NewClientWithAPI(api containerAPI)`: creates a client around a test or alternate container-list implementation.
- `Client.ListRunningContainers(ctx context.Context)`: returns normalized running containers or a wrapped discovery error.

## Data tables

None.

## Pipeline steps

The CLI constructs a Docker client, lists running containers, then passes normalized `domain.Container` values into the UI selection model. Future streaming code should extend this module instead of calling the Docker SDK directly from UI or stream packages.

## Routes

None.

## Configuration

The Docker SDK reads standard Docker environment variables such as `DOCKER_HOST` through `client.FromEnv`.

## Notes

Keep Docker SDK types inside this module. Tests should use the small internal interface rather than a live Docker daemon.
