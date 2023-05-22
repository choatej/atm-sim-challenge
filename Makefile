# Phony targets (do not correspond to actual files)
.PHONY: build test lint

# Binary output name
BINARY := atm-sim

# install path
ifeq ($(shell uname -s),Windows_NT)
	INSTALL_PATH := C:\your\install\path
else ifeq ($(shell uname -s),Darwin)
	INSTALL_PATH := /usr/local/bin
else
	INSTALL_PATH := /usr/local/bin
endif

all: clean fmt lint build build-all-arch test install docker

# Build target
build-all-arch: bindata
	@echo "Building for Windows (amd64)..."
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o $(BINARY)_windows_amd64 data.go atm-sim.go

	@echo "Building for macOS (amd64)..."
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o $(BINARY)_darwin_amd64 data.go atm-sim.go

	@echo "Building for macOS (arm64)..."
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o $(BINARY)_darwin_arm64 data.go atm-sim.go

	@echo "Building for Linux (amd64)..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BINARY)_linux_amd64 data.go atm-sim.go

	@echo "Building for Linux (arm64)..."
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o $(BINARY)_linux_arm64 data.go atm-sim.go

bindata:
	go-bindata -o data.go ./data

build: bindata
	@echo "building $(BINARY) for $(shell go env GOOS)_$(shell go env GOARCH)"
	go build -o $(BINARY)_$(shell go env GOOS)_$(shell go env GOARCH) data.go atm-sim.go

# Run unit tests
test: bindata
	go test -covermode=count -coverprofile=coverage.out ./... && \
   	go tool cover -html=coverage.out

# Run goimports to format Go code
lint: bindata
	golangci-lint run

fmt:
	go fmt $$(go list ./... | grep -v /vendor/)

install: build
	cp $(BINARY)_$(shell go env GOOS)_$(shell go env GOARCH) $(INSTALL_PATH)/$(BINARY)
	@echo "Installation complete."

docker: build-all-arch
	docker build -t atm-sim:latest .

# Clean target
clean:
	rm -f $(BINARY) coverage.out logfile.log atm-sim_* data.go


