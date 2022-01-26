GITTAG ?= $(shell git describe --tags --always)
GITCOMMIT ?= $(shell git log -1 --pretty=format:"%H")
GOLDFLAGS ?= -s -w -extldflags '-zrelro -znow'
GOFLAGS ?= -trimpath
CGO_ENABLED ?= 0

.PHONY: all
all: build

.PHONY: build
build: fmtcheck go-lint
	go build -ldflags "$(GOLDFLAGS)" ./...

.PHONY: depscheck
depscheck:
	@echo "==> Checking source code dependencies..."
	@go mod tidy
	@git diff --exit-code -- go.mod go.sum || \
		(echo; echo "Found differences in go.mod/go.sum files. Run 'go mod tidy' or revert go.mod/go.sum changes."; exit 1)
	@# reset go.sum to state before checking if it is clean
	@git checkout -q go.sum

.PHONY: benchmark
benchmark:
	go test -bench=. -benchmem ./...

.PHONY: test
test: tools
	go test 				\
		-race 				\
		-coverprofile coverage.out 	\
		./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: go-lint
go-lint: tools
	@echo "==> Checking source code against linters..."
	@tools/golangci-lint run ./...
	@tools/golangci-lint run --build-tags integration ./...

.PHONY: docs-lint
docs-lint: tools
	@echo "==> Checking docs against linters..."
	@tools/misspell -error -source=text README.md || (echo; \
		echo "Unexpected misspelling found in docs files."; \
		echo "To automatically fix the misspelling, run 'make docs-lint-fix' and commit the changes."; \
		exit 1)
	@docker run -v $(PWD):/markdown 06kellyjac/markdownlint-cli README.md || (echo; \
		echo "Unexpected issues found in docs Markdown files."; \
		echo "To apply any automatic fixes, run 'make docs-lint-fix' and commit the changes."; \
		exit 1)

.PHONY: docs-lint-fix
docs-lint-fix: tools
	@echo "==> Applying automatic docs linter fixes..."
	@tools/misspell -w -source=text docs/
	@docker run -v $(PWD):/markdown 06kellyjac/markdownlint-cli --fix docs/

.PHONY: lint
lint: go-lint docs-lint


.PHONY: fmt
fmt:
	gofmt -s -w .

.PHONY: fmtcheck
fmtcheck:
	@./scripts/gofmtcheck.sh

.PHONY: tools
tools:
	cd tools && go build -o . github.com/client9/misspell/cmd/misspell
	cd tools && go build -o . github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: install-precommit-hook
install-precommit-hook: .git/hooks/pre-commit

.PHONY: .git/hooks/pre-commit
.git/hooks/pre-commit: scripts/pre-commit
	cp $< $@
	chmod +x $@
