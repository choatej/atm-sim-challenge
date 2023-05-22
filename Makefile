# Phony targets (do not correspond to actual files)
.PHONY: build bindata test lint clean

# Go compiler
GO := go

# Binary output name
BINARY := bin/atm-sim

# Directories containing Go source files
SOURCES := main.go cmd internal

# Source file dependencies
SOURCES_DEPS := $(shell find $(SOURCES) -name "*.go")

# Path to the accounts.csv file
ACCOUNTS_CSV := data/accounts.csv

# Build target
build: create-bin
	$(GO) build -o $(BINARY) $(SOURCES_DEPS)

# Include accounts.csv into the binary using go-bindata
bindata: create-bin
	go get -u github.com/go-bindata/go-bindata/...
	go-bindata -o bindata.go -pkg main -prefix $(dir $(ACCOUNTS_CSV)) $(ACCOUNTS_CSV)

# Run unit tests
test:
	$(GO) test ./...

# Run goimports to format Go code
lint:
	go get golang.org/x/tools/cmd/goimports
	goimports -w $(SOURCES)

# Create bin directory
create-bin:
	mkdir -p bin

# Clean target
clean:
	rm -f $(BINARY) bindata.go


