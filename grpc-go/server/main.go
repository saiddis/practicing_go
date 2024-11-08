package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/saiddis/grpc-go/helloworld"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type serverHelloWorld struct {
	pb.UnimplementedHelloWorldServiceServer
}

type serverGreeter struct {
	pb.UnimplementedGreeterServer
}

func (s *serverHelloWorld) SayHello(ctx context.Context, in *pb.HelloWorldRequest) (*pb.HelloWorldResponse, error) {
	return &pb.HelloWorldResponse{Message: "Hello, World!"}, nil
}

func (s *serverGreeter) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello again " + in.GetName()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen on port 50051: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterHelloWorldServiceServer(s, &serverHelloWorld{})
	pb.RegisterGreeterServer(s, &serverGreeter{})
	log.Printf("gRPC server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
