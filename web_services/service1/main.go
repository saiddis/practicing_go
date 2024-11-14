package main

import (
	"github.com/saiddis/practicing_go/web_services/service1/pkg/server"
	"log"
)

func main() {
	server, err := server.New("localhost",
		server.WithPort(8080),
		server.WithNatsAddr("nats://nats:4222"))

	if err != nil {
		log.Fatalf("error creating server: %v", err)
	}
	err = server.Run()

	//nc, _ := nats.Connect("nats://nats:4222")
	//defer nc.Drain()

	//nc.Publish("greet.joe", []byte("Hi!"))

	//sub, _ := nc.SubscribeSync("greet.*")

	//msg, _ := sub.NextMsg(10 * time.Millisecond)
	//log.Println("subscribed after a publish...")
	//log.Printf("msg is nil? %v\n", msg == nil)

	//nc.Publish("greet.joe", []byte("Hi!"))
	//nc.Publish("greet.pam", []byte("Hi!"))

	//msg, _ = sub.NextMsg(10 * time.Millisecond)
	//log.Printf("msg data: %q on subject %q\n", string(msg.Data), msg.Subject)

	//msg, _ = sub.NextMsg(10 * time.Millisecond)
	//log.Printf("msg data: %q on subject %q\n", string(msg.Data), msg.Subject)

	//nc.Publish("greet.bob", []byte("hello"))

	//msg, _ = sub.NextMsg(10 * time.Millisecond)
	//log.Printf("msg data: %q on subject %q\n", string(msg.Data), msg.Subject)
}
