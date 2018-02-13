GOCMD=go
GOPATH=~/go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
PROJECT_NAME=ProjectX
TARGET_OS=linux
BINARY_NAME=main_$(TARGET_OS)

NO_COLOR=\x1b[0m
OK_COLOR=\x1b[32;01m
ERROR_COLOR=\x1b[31;01m
WARN_COLOR=\x1b[33;01m

OK_STRING=$(OK_COLOR)[OK]$(NO_COLOR)
ERROR_STRING=$(ERROR_COLOR)[ERRORS]$(NO_COLOR)
WARN_STRING=$(WARN_COLOR)[WARNINGS]$(NO_COLOR)

.PHONY: all

all: build # test

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN) -x
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

compose:
	@echo "${OK_COLOR}Building binary file${NO_COLOR}" 
	@make docker-build
	@echo "${OK_COLOR}Start docker containers${NO_COLOR}"
	docker-compose up --build