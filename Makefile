GOLANGCI_VERSION = 1.34.0
BIN_DIR ?= $(shell pwd)/bin

ci: lint test
.PHONY: ci

$(BIN_DIR)/golangci-lint: $(BIN_DIR)/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} $(BIN_DIR)/golangci-lint
$(BIN_DIR)/golangci-lint-${GOLANGCI_VERSION}:
	@curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | BINARY=golangci-lint bash -s -- v${GOLANGCI_VERSION}
	@mv $(BIN_DIR)/golangci-lint $@

$(BIN_DIR)/mockgen:
	@env GOBIN=$$PWD/bin GO111MODULE=on go install github.com/golang/mock/mockgen

mocks: $(BIN_DIR)/mockgen
	@echo "--- build all the mocks"
	@$(BIN_DIR)/mockgen -destination=mocks/handler.go -package=mocks github.com/aws/aws-lambda-go/lambda Handler
.PHONY: mocks

lint: $(BIN_DIR)/golangci-lint
	@echo "--- lint all the things"
	@$(BIN_DIR)/golangci-lint run ./...
	@cd ./middleware/raw; $(BIN_DIR)/golangci-lint run
	@cd ./middleware/zerolog; $(BIN_DIR)/golangci-lint run
.PHONY: lint

test:
	@echo "--- test all the things"
	@go test -v -cover ./lambdaextras/... ./middleware/...
	@cd ./middleware/raw; go test -v -cover ./...
	@cd ./middleware/zerolog; go test -v -cover ./...
.PHONY: test
