# syntax=docker/dockerfile:1

FROM golang:1.23

RUN apt update && apt upgrade -y &&\
	apt install -y git\
	make openssh-client

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy

RUN go build -o /main

CMD ["/main"]
