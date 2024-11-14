# syntax=docker/dockerfile:1

FROM golang:1.23 AS builder

RUN apt update && apt upgrade -y &&\
	apt install -y git\
	make openssh-client

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN go mod tidy

RUN go build -o /main

CMD ["/main"]
