SOURCES := $(shell find . -name '*.go')
TARGET := ./dist/flashbiter_darwin_amd64_v1/flashbiter

run: flashbiter
	./flashbiter

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
