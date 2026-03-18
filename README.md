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
dev init <language>  # generate dev.yaml with all defaults
dev build            # build the project
dev test             # run tests
dev lint             # run linters
dev fmt              # format source code
dev clean            # remove build artifacts
dev run              # run the project
dev sync             # sync dependencies (go mod tidy / uv sync)
dev scan             # run security scans
dev ci               # lint → test → build pipeline
dev setup            # install dependencies
dev doctor           # validate the environment
dev version          # print version
```

All commands support:

```bash
--dry-run       # print commands without executing
--verbose       # verbose output
--cwd <path>    # run as if invoked from <path>
```

## Configuration

Run `dev init <language>` to generate a `dev.yaml` pre-populated with all defaults:

```bash
cd my-project
dev init go
```

This writes a `dev.yaml` you can tune:

```yaml
language: go

tasks:
  build:
    cmd: go build ./...
  test:
    cmd: go test -race ./...
  lint:
    cmd: golangci-lint run
  ci:
    deps: [lint, test, build]
```

If no `dev.yaml` is present, language is auto-detected from `go.mod` / `pyproject.toml` and built-in defaults apply.

### Override model

1. Built-in defaults per language
2. Project `dev.yaml` overrides defaults
3. CLI flags override everything

### Composing tools with deps

Use `deps` to compose multiple tools under a single command. Sub-tasks are named entries in `dev.yaml` — they are **not** top-level `dev` commands, only reachable through the task runner.

```yaml
tasks:
  ruff:
    cmd: ruff check .
  pre-commit:
    cmd: pre-commit run --all-files
  lint:
    deps: [ruff, pre-commit]
```

```yaml
tasks:
  trivy:
    cmd: trivy fs .
  opengrep:
    cmd: opengrep scan .
  scan:
    deps: [trivy, opengrep]
```

Remove a tool from `deps` or change its `cmd` without touching anything else.

### pre-commit

`pre-commit` integrates naturally through the task model. Add it to `setup` so hooks are installed on bootstrap, and to `lint` as a composed dep:

```yaml
tasks:
  setup:
    cmd: pre-commit install
  pre-commit:
    cmd: pre-commit run --all-files
  lint:
    deps: [golangci-lint, pre-commit]
  golangci-lint:
    cmd: golangci-lint run
```

### Global defaults

Set defaults for all projects in `~/.dev/config.yaml`:

```yaml
defaults:
  lint: golangci-lint run --fix
```

## Language support

| Language | Detection | Build | Test | Lint | Fmt | Sync |
|---|---|---|---|---|---|---|
| Go | `go.mod` | `go build ./...` | `go test ./...` | `golangci-lint run` | `gofmt -w .` | `go mod tidy` |
| Python | `pyproject.toml` / `setup.py` | `python -m build` | `pytest` | `ruff check .` | `ruff format .` | `uv sync` |

Adding a new language requires implementing a small interface and registering it — no plugin system needed.

## Development

```bash
make build          # build ./dev binary
make test           # run all tests
make lint           # run golangci-lint
make update-golden  # regenerate golden test fixtures
make release-dry    # local goreleaser snapshot
```
