# Go parameters
BINARY_NAME=elementarysite
PROJECTNAME=$(shell basename "$(PWD)")

GO_BIN ?= go
GOBUILD=$(GO_BIN) build
GOCLEAN=$(GO_BIN) clean
GOTEST=$(GO_BIN) test

hello:
	echo "Hello from "$(PROJECTNAME);
    
all: test build

tidy:
	if eq ($(GO111MODULE),on) \
    	$(GO_BIN) mod tidy;\
    else \
    	echo skipping go mod tidy;\
    endif;

build:
	$(GOBUILD) -o build/x86_64/$(BINARY_NAME) -v cmd/webapp/main.go;

build32:
	GOOS=linux GOARCH=386 $(GOBUILD) -o build/i686/$(BINARY_NAME) -v cmd/webapp/main.go;

test:
	$(GOTEST) -v ./...;

clean:
	$(GOCLEAN);
	rm -f $(BINARY_NAME);

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...;\
./$(BINARY_NAME)


# Cross compilation
build-win:\
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)"_win64.exe" -v
