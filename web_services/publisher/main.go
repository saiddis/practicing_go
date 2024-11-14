package main

import (
	"github.com/saiddis/practicing_go/web_services/publisher/pkg/server"
	"log"
)

func main() {
	server, err := server.New("0.0.0.0",
		server.WithPort(8080),
		server.WithNatsAddr("nats://nats:4222"))

	if err != nil {
		log.Fatalf("error creating server: %v", err)
	}
	err = server.Run()
}
