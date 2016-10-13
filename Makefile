# go-structmapper Makefile
# This Makefile is not required for using go-structmapper and only
# provides functionality that is used during development

GO_BUILDFLAGS?=-v
GO_TEST_FLAGS?=-v
GO_BENCHMARK_FLAGS?=-run="INVALID" -v
ROOT_PACKAGE_NAME=github.com/anexia-it/go-structmapper

PACKAGES=$(shell go list ./... | \
	grep -v '/vendor/' | \
	sed \
		-e 's|^$(ROOT_PACKAGE_NAME)/|./|g' \
		-e 's|^$(ROOT_PACKAGE_NAME)$$|./|g' \
	)
SRC_FILES=$(shell find . -name '*.go' | egrep -v '^./vendor/')

all: build

format:
	@echo "Running gofmt..."
	@gofmt -w -s -r '(a) -> a' $(SRC_FILES)
	@echo "Running goimports..."
	@goimports -e=true -w=true $(SRC_FILES)

generate:
	@echo "Generating..."
	@go generate $(PACKAGES)

build: format generate
	@echo "Building packages..."
	@go install $(GO_BUILDFLAGS) $(PACKAGES)

test: build
	@echo "Running unit tests..."
	@rm -rf ./coverage.out
	@mkdir -p ./coverage.out
	@$(foreach pkg,$(PACKAGES), \
		go test $(GO_TEST_FLAGS) -cover -coverprofile \
			./coverage.out/$(shell echo $(pkg) | \
		sed -e 's|/|_|g' -e 's|^.||g').coverage $(pkg);)
	@echo "mode: set" > ./coverage.out/combined
	@cat ./coverage.out/*.coverage | \
		sed -e '/mode:.*/d' >> ./coverage.out/combined
	@go tool cover -html ./coverage.out/combined -o ./coverage.html
	@go tool cover -func=./coverage.out/combined | grep 'total' | \
		awk '{print "coverage (all): " $$3 " of statements"; }'

benchmark:
	@echo "Running benchmarks..."
	@$(foreach pkg,$(PACKAGES), \
		go test -bench="." $(GO_BENCHMARK_FLAGS) $(pkg);)
