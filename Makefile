SOURCES := $(wildcard *.go) $(wildcard **/*.go)

ifeq ($(shell uname),Darwin)
    GOOS = darwin
    GOARCH = amd64
    EXEEXT =
else ifeq ($(shell uname),Linux)
    GOOS = linux
    GOARCH = $(shell arch)
    EXEEXT =
else ifeq ($(shell uname),Windows_NT)
    GOOS = windows
    GOARCH = amd64
    EXEEXT = .exe
endif

TARGET := ./dist/flashbiter_$(GOOS)_$(GOARCH)_v1/flashbiter

flashbiter: $(TARGET)
	cp $< $@

$(TARGET): $(SOURCES)
	gofumpt -w $(SOURCES)
	goreleaser build --single-target --snapshot --clean
	go vet ./...

.PHONY: clean
clean:
	rm -f flashbiter
	rm -f $(TARGET)
	rm -rf dist
