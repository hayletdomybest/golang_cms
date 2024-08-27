# Variables
APP_NAME := weex_admin
CMD_DIR := cmd/http
MAIN_FILE := main.go
BIN_DIR := bin

# Go commands
BUILD := go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)/$(MAIN_FILE)
RUN := go run $(MAIN_FILE)
FMT := go fmt ./...
VET := go vet ./...
TEST := go test ./...
LINT := golangci-lint run

# Default target
.PHONY: all
all: fmt vet build

# Build the project
.PHONY: build
build:
	@echo "Building the application..."
	$(BUILD)

# Run the project
.PHONY: run
run:
	@echo "Running the application..."
	$(RUN)

# Format the code
.PHONY: fmt
fmt:
	@echo "Formatting the code..."
	$(FMT)

# Lint the code
.PHONY: lint
lint:
	@echo "Linting the code..."
	$(LINT)

# Run go vet
.PHONY: vet
vet:
	@echo "Running go vet..."
	$(VET)

# Test the project
.PHONY: test
test:
	@echo "Running tests..."
	$(TEST)

# Clean the build
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -f $(APP_NAME)

# Generate wire files
.PHONY: wire
wire:
	@echo "Generating wire files..."
	cd $(CMD_DIR)/di && wire

# Regenerate wire files and build
.PHONY: wire-build
wire-build: wire build
