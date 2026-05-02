# Domain Module

## What it does

Defines shared domain types used across Docker discovery, stream management, filtering, and terminal UI code.

## Public API

- `Container`: a normalized running Docker container with ID, name, image, and status fields.
- `Container.DisplayName()`: returns the preferred label for UI and log prefixes.

## Data tables

None.

## Pipeline steps

Domain values are produced by adapter modules and consumed by stream and UI modules. This package does not perform I/O.

## Routes

None.

## Configuration

None.

## Notes

Keep this module free of Docker SDK and Bubble Tea imports so it can remain the stable boundary between infrastructure and UI packages.
