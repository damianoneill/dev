BIN     := dev
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/damianoneill/dev/cmd.Version=$(VERSION)"

.PHONY: build test lint fmt clean install release-dry update-golden

build:
	go build $(LDFLAGS) -o $(BIN) .

test:
	go test ./...

lint:
	golangci-lint run

fmt:
	gofmt -w .

clean:
	rm -f $(BIN)

install:
	go install $(LDFLAGS) .

# Regenerate golden test fixtures from current binary output.
update-golden: build
	./$(BIN) version               > testdata/golden/version.txt
	./$(BIN)                       > testdata/golden/help.txt
	./$(BIN) build --dry-run --cwd testdata/projects/go_project     > testdata/golden/build_dry_run_go.txt
	./$(BIN) ci    --dry-run --cwd testdata/projects/go_project     > testdata/golden/ci_dry_run_go.txt
	./$(BIN) build --dry-run --cwd testdata/projects/python_project > testdata/golden/build_dry_run_python.txt
	./$(BIN) ci    --dry-run --cwd testdata/projects/python_project > testdata/golden/ci_dry_run_python.txt

release-dry:
	goreleaser release --snapshot --clean
