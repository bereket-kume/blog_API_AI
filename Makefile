# Makefile for Blog API Tests

# Variables
GO=go
TEST_FLAGS=-v -race -cover
COVERAGE_DIR=coverage
TEST_TIMEOUT=30s

# Test targets
.PHONY: test test-unit test-integration test-all coverage clean

# Run all tests
test-all: test-unit test-integration

# Run unit tests (with mocks)
test-unit:
	@echo "Running unit tests..."
	$(GO) test $(TEST_FLAGS) -timeout $(TEST_TIMEOUT) ./usecases/... ./Delivery/controllers/...

# Run integration tests (with real database)
test-integration:
	@echo "Running integration tests..."
	$(GO) test $(TEST_FLAGS) -timeout $(TEST_TIMEOUT) ./Infrastructure/repositories/...

# Run specific test files
test-usecase:
	@echo "Running blog usecase tests..."
	$(GO) test $(TEST_FLAGS) -timeout $(TEST_TIMEOUT) ./usecases/blog_usecase_test.go

test-controller:
	@echo "Running blog controller tests..."
	$(GO) test $(TEST_FLAGS) -timeout $(TEST_TIMEOUT) ./Delivery/controllers/blog_controller_test.go

test-repository:
	@echo "Running blog repository tests..."
	$(GO) test $(TEST_FLAGS) -timeout $(TEST_TIMEOUT) ./Infrastructure/repositories/blog_repository_test.go

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	$(GO) test $(TEST_FLAGS) -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated at $(COVERAGE_DIR)/coverage.html"

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	$(GO) test -v -timeout $(TEST_TIMEOUT) ./...

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	$(GO) test -race -timeout $(TEST_TIMEOUT) ./...

# Run tests with short timeout (for CI/CD)
test-short:
	@echo "Running tests with short timeout..."
	$(GO) test $(TEST_FLAGS) -timeout 10s ./...

# Clean up test artifacts
clean:
	@echo "Cleaning up test artifacts..."
	rm -rf $(COVERAGE_DIR)
	go clean -testcache

# Install test dependencies
test-deps:
	@echo "Installing test dependencies..."
	go get github.com/stretchr/testify/assert
	go get github.com/stretchr/testify/mock
	go get github.com/stretchr/testify/suite
	go mod tidy

# Check if MongoDB is running (for integration tests)
check-mongo:
	@echo "Checking local MongoDB connection..."
	@if ! pgrep -x "mongod" > /dev/null; then \
		echo "Local MongoDB is not running. Please start it manually:"; \
		echo "  sudo systemctl start mongod"; \
		exit 1; \
	fi
	@if command -v mongosh > /dev/null 2>&1; then \
		if ! mongosh --eval "db.runCommand('ping')" > /dev/null 2>&1; then \
			echo "Cannot connect to MongoDB. Please check your configuration."; \
			exit 1; \
		fi; \
	elif command -v mongo > /dev/null 2>&1; then \
		if ! mongo --eval "db.runCommand('ping')" > /dev/null 2>&1; then \
			echo "Cannot connect to MongoDB. Please check your configuration."; \
			exit 1; \
		fi; \
	elif ! nc -z localhost 27017 2>/dev/null; then \
		echo "Cannot connect to MongoDB port 27017. Please check your configuration."; \
		exit 1; \
	fi
	@echo "MongoDB is running and accessible"

# Run integration tests with MongoDB check
test-integration-mongo: check-mongo test-integration

# Run all tests with MongoDB check
test-all-mongo: check-mongo test-all

# Help target
help:
	@echo "Available test targets:"
	@echo "  test-all          - Run all tests (unit + integration)"
	@echo "  test-unit         - Run unit tests only"
	@echo "  test-integration  - Run integration tests only"
	@echo "  test-usecase      - Run blog usecase tests"
	@echo "  test-controller   - Run blog controller tests"
	@echo "  test-repository   - Run blog repository tests"
	@echo "  coverage          - Run tests with coverage report"
	@echo "  test-verbose      - Run tests with verbose output"
	@echo "  test-race         - Run tests with race detection"
	@echo "  test-short        - Run tests with short timeout"
	@echo "  test-deps         - Install test dependencies"
	@echo "  check-mongo       - Check/start MongoDB container"
	@echo "  clean             - Clean up test artifacts"
	@echo "  help              - Show this help message" 