package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"seotrang.com/rgpc-clint-server/greeterpb"
)

type server struct {
	greeterpb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *greeterpb.HelloRequest) (*greeterpb.HelloResponse, error) {
	name := req.GetName()
	res := &greeterpb.HelloResponse{
		Message: "Hello " + name,
	}
	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greeterpb.RegisterGreeterServer(s, &server{})

	fmt.Println("🚀 Server is running at :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
