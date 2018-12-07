GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=candump
BUILD_DIR=build
BUILD_SRC=cmd/candump.go
PACKAGE_RPI=$(BINARY_NAME)_linux_armhf

all: test build
build:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -i $(BUILD_SRC)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -rf $(BINARY_NAME)

package-rpi: build-rpi
	tar -cvzf $(PACKAGE_RPI).tar.gz -C $(BUILD_DIR) $(BINARY_NAME)

build-rpi:
	GOOS=linux GOARCH=arm GOARM=6 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -i $(BUILD_SRC)