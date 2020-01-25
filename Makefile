# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=ectoplasma

all: test build
build_server: 
	$(GOBUILD) -o bin/$(BINARY_NAME)_server cmd/server/main.go
build_scraper: 
	$(GOBUILD) -o bin/$(BINARY_NAME)_scraper cmd/scraper/main.go
build_processor: 
	$(GOBUILD) -o bin/$(BINARY_NAME)_processor cmd/processor/main.go
build: build_server build_scraper build_processor
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -rf bin/
