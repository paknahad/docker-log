# Stack Applied

`docs/stack.md` was present but empty during bootstrap. The applied stack was inferred from `CLAUDE.md`, `docs/product.md`, and `docs/architecture.md`.

## Applied stack

- Go CLI application
- Bubble Tea terminal UI planned for interactive views
- Lip Gloss planned for terminal styling
- Docker Go SDK planned for Docker access
- `go test`, `go vet`, `staticcheck`, and `gofmt` for local verification

## Notes

The first bootstrap adds the Go toolchain, module, Makefile checks, and a minimal compiling CLI skeleton. The GitHub Actions workflow still needs a token with `workflow` scope before it can be updated for Go checks. Feature work remains queued in GitHub issues.
