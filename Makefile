GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

BINARY_NAME=candump
VERSION=$(shell git describe --exact-match --tags 2>/dev/null)
BUILD_DIR=build
BUILD_SRC=cmd/candump.go

all: test build
build:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -i $(BUILD_SRC)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -rf $(BINARY_NAME)

package-rpi: build-rpi
	tar -cvzf $(BINARY_NAME)-$(VERSION)_linux_armhf.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)

build-rpi:
	GOOS=linux GOARCH=arm GOARM=6 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -i $(BUILD_SRC)