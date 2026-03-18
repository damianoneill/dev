# dev

A consistent developer CLI across all your projects.

`dev` provides a unified command interface for build, test, lint, format, and more — regardless of language or project structure. Drop it in any Go or Python repo and get the same commands everywhere.

## Install

```bash
go install github.com/damianoneill/dev@latest
```

Or download a binary from the [releases page](https://github.com/damianoneill/dev/releases).

## Usage

```bash
dev build       # build the project
dev test        # run tests
dev lint        # run linters
dev fmt         # format source code
dev clean       # remove build artifacts
dev run         # run the project
dev ci          # lint → test → build pipeline
dev setup       # install dependencies
dev doctor      # validate the environment
dev version     # print version
```

All commands support:

```bash
--dry-run       # print commands without executing
--verbose       # verbose output
--cwd <path>    # run as if invoked from <path>
```

## Configuration

Add a `dev.yaml` to your project root to customise commands or define tasks:

```yaml
language: go   # optional — detected from go.mod / pyproject.toml

tasks:
  build:
    cmd: go build -o bin/app ./cmd/app
  test:
    cmd: go test -race ./...
  ci:
    deps: [lint, test, build]
```

If no `dev.yaml` is present, language is auto-detected and sensible defaults apply.

### Override model

1. Built-in defaults per language
2. Project `dev.yaml` overrides defaults
3. CLI flags override everything

### Global defaults

Set defaults for all projects in `~/.dev/config.yaml`:

```yaml
defaults:
  lint: golangci-lint run --fix
```

## Language support

| Language | Detection    | Build          | Test    | Lint          | Fmt            |
|----------|--------------|----------------|---------|---------------|----------------|
| Go       | `go.mod`     | `go build ./…` | `go test ./…` | `golangci-lint run` | `gofmt -w .` |
| Python   | `pyproject.toml` / `setup.py` | `python -m build` | `pytest` | `ruff check .` | `ruff format .` |

Adding a new language requires implementing a small interface and registering it — no plugin system needed.

## Development

```bash
make build          # build ./dev binary
make test           # run all tests
make lint           # run golangci-lint
make update-golden  # regenerate golden test fixtures
make release-dry    # local goreleaser snapshot
```
