package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"seotrang.com/rgpc-clint-client/greeterpb"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ could not connect: %v", err)
	}
	defer conn.Close()

	client := greeterpb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for i := 0; i < 15; i++ {
		start := time.Now() // Thời gian bắt đầu gửi request

		res, err := client.SayHello(ctx, &greeterpb.HelloRequest{Name: "Go Developer"})

		elapsed := time.Since(start) // Thời gian sau khi nhận response

		if err != nil {
			log.Fatalf("❌ could not greet: %v", err)
		}

		fmt.Println("👋 Response:", res.GetMessage())
		fmt.Printf("⏱️ Response time: %s\n", elapsed)
	}

}
