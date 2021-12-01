# These will be provided to the target
TARGET := "benford"
VERSION := 1.0.0
BUILD := `git rev-parse --short HEAD`

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.BuildCommitShort=$(BUILD)"

.PHONY: build clean install uninstall

build:
	@go build $(LDFLAGS) -o $(TARGET)

clean:
	@rm -f $(TARGET)

install:
	@go install $(LDFLAGS)

uninstall: clean
	@rm -f $$(which ${TARGET})
