package scaffold

var gitignoreTmpl = `# dev binary
dev
dist/
`

var devYAMLTmpl = `language: {{.Language}}

tasks:
  ci:
    deps: [lint, test, build]
`

var readmeTmpl = `# {{.Name}}

> TODO: describe this project.

## Usage

` + "```" + `bash
dev build
dev test
dev lint
dev ci
` + "```" + `
`

var ciTmpl = `name: ci

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  ci:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: install dev
        run: go install github.com/damianoneill/dev@latest
{{if eq .Language "go"}}
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true
{{end}}{{if eq .Language "python"}}
      - uses: actions/setup-python@v5
        with:
          python-version: "3.12"
{{end}}
      - run: dev setup
      - run: dev ci
`

// Go-specific templates

// Go-specific templates (exported for use by the golang language package).

var GoMainTmpl = `package main

import "fmt"

func main() {
	fmt.Println("{{.Name}}")
}
`

var GolangciTmpl = `version: "2"

linters:
  default: none
  enable:
    - errcheck
    - govet
    - staticcheck
    - unused

formatters:
  enable:
    - gofmt

issues:
  max-issues-per-linter: 0
  max-same-issues: 0
`

// Python-specific templates (exported for use by the python language package).

var PyprojectTmpl = `[project]
name = "{{.Name}}"
version = "0.1.0"
requires-python = ">=3.11"

[dependency-groups]
dev = ["pytest", "ruff"]
`
