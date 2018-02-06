GOCMD=go
GOPATH=~/go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
PROJECT_NAME=ProjectX
TARGET_OS=darwin
BINARY_NAME=main_$(TARGET_OS)

.PHONY: all

all: build # test

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

deps:
	$(GOGET) github.com/golang/dep/cmd/dep
	dep ensure

build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v

docker-build:
	docker run --rm -it -v $(GOPATH):/go -w /go/src/github.com/maxkrivich/$(PROJECT_NAME)/ --env GOOS=$(TARGET_OS) golang:latest go build -o $(BINARY_NAME) -v