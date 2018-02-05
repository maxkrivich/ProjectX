FROM golang:1.9.3
MAINTAINER "Maxim Krivich"
ENV GOBIN /go/bin

RUN mkdir -p /go/src/github.com/maxkrivich/ProjectX
ADD . /go/src/github.com/maxkrivich/ProjectX
WORKDIR /go/src/github.com/maxkrivich/ProjectX

RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure

EXPOSE 8080