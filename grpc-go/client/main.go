package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/saiddis/grpc-go/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	DefaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", DefaultName, "Name to greet")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server at localhost:50051: %v", err)
	}
	defer conn.Close()
	c := pb.NewHelloWorldServiceClient(conn)
	cAgain := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloWorldRequest{})
	if err != nil {
		log.Fatalf("error calling function SayHello: %v", err)
	}

	rAgain, err := cAgain.SayHelloAgain(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("error calling function SayHelloAgain: %v", err)
	}

	log.Printf("Response from gRPC server's SayHello function: %s", r.GetMessage())
	log.Printf("Response from gRPC server's SayHelloAgain function: %s", rAgain.GetMessage())
}
