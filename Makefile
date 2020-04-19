GOLANGCI_VERSION = 1.24.0

ci: lint test
.PHONY: ci

LDFLAGS := -ldflags="-s -w"

bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint
bin/golangci-lint-${GOLANGCI_VERSION}:
	@curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | BINARY=golangci-lint bash -s -- v${GOLANGCI_VERSION}
	@mv bin/golangci-lint $@

bin/mockgen:
	@env GOBIN=$$PWD/bin GO111MODULE=on go install github.com/golang/mock/mockgen

bin/gcov2lcov:
	@env GOBIN=$$PWD/bin GO111MODULE=on go install github.com/jandelgado/gcov2lcov

mocks: bin/mockgen
	@echo "--- build all the mocks"
	@bin/mockgen -destination=mocks/handler.go -package=mocks github.com/aws/aws-lambda-go/lambda Handler
.PHONY: mocks

lint: bin/golangci-lint
	@echo "--- lint all the things"
	@bin/golangci-lint run
.PHONY: lint

test: bin/gcov2lcov
	@echo "--- test all the things"
	@go test -v -covermode=count -coverprofile=coverage.txt ./lambdaextras/... ./middleware/...
	@bin/gcov2lcov -infile=coverage.txt -outfile=coverage.lcov
.PHONY: test
