GOLANGCI_VERSION = latest
BIN_DIR ?= $(shell pwd)/bin

ci: lint test
.PHONY: ci

$(BIN_DIR)/golangci-lint: $(BIN_DIR)/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} $(BIN_DIR)/golangci-lint
$(BIN_DIR)/golangci-lint-${GOLANGCI_VERSION}:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s $(GOLANGCI_VERSION)
	@mv $(BIN_DIR)/golangci-lint $@

lint: $(BIN_DIR)/golangci-lint
	@echo "--- lint all the things"
	@$(BIN_DIR)/golangci-lint run ./...
	@cd ./standard; $(BIN_DIR)/golangci-lint run
	@cd ./middleware/raw; $(BIN_DIR)/golangci-lint run
	@cd ./middleware/zerolog; $(BIN_DIR)/golangci-lint run
.PHONY: lint

lint-fix: $(BIN_DIR)/golangci-lint
	@echo "--- lint all the things"
	@$(BIN_DIR)/golangci-lint run --fix ./...
	@cd ./standard; $(BIN_DIR)/golangci-lint run --fix ./...
	@cd ./middleware/raw; $(BIN_DIR)/golangci-lint run --fix ./...
	@cd ./middleware/zerolog; $(BIN_DIR)/golangci-lint run --fix ./...
.PHONY: lint-fix

test:
	@echo "--- test all the things"
	@go test -v -cover ./lambdaextras/... ./middleware/...
	@cd ./standard; go test -v -cover ./...
	@cd ./middleware/raw; go test -v -cover ./...
	@cd ./middleware/zerolog; go test -v -cover ./...
.PHONY: test
