.PHONY: build test vet lint fmt clean

GO ?= go
GOFMT ?= gofmt
GOLANGCI_LINT ?= golangci-lint

BINARY ?= ksearch
GOOS ?= $(shell $(GO) env GOOS)
GOARCH ?= $(shell $(GO) env GOARCH)
CGO_ENABLED ?= 0
OUTPUT ?= $(BINARY)$(if $(filter windows,$(GOOS)),.exe,)

build:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -o $(OUTPUT) .

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

lint:
	$(GOLANGCI_LINT) run --timeout=5m ./...

fmt:
	$(GOFMT) -w .

clean:
	rm -f $(BINARY) $(BINARY).exe
