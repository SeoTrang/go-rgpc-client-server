package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"seotrang.com/rgpc-clint-server/greeterpb"
)

type server struct {
	greeterpb.UnimplementedGreeterServer
}

type User struct {
	id   int32
	name string
	age  int32
}

var Users = []User{
	{id: 1, name: "Alice", age: 25},
	{id: 2, name: "Bob", age: 30},
	{id: 3, name: "Charlie", age: 22},
	{id: 4, name: "David", age: 35},
	{id: 5, name: "Eve", age: 28},
}

func (s *server) SayHello(ctx context.Context, req *greeterpb.HelloRequest) (*greeterpb.HelloResponse, error) {
	name := req.GetName()
	res := &greeterpb.HelloResponse{
		Message: "Hello " + name,
	}
	return res, nil
}

func (s *server) GetUser(ctx context.Context, req *greeterpb.GetUserRequest) (*greeterpb.GetUserResponse, error) {
	var id int32 = req.GetId()
	for _, user := range Users {
		if user.id == id {
			return &greeterpb.GetUserResponse{
				Id:   user.id,
				Name: user.name,
				Age:  user.age,
			}, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "User with ID %d not found", id)

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
