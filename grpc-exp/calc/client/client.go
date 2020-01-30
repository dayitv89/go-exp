package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/dayitv89/go-exp/grpc-exp/calc/calcpb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50050", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCalculatorClient(conn)

	req := &pb.Request{N1: 10, N2: 15}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("could not sum: %v", err)
	}
	fmt.Printf("Sum: %d\n", res.GetResult())

	res, err = c.Subtract(context.Background(), req)
	if err != nil {
		log.Fatalf("could not Subtract: %v", err)
	}
	fmt.Printf("Subtract: %d\n", res.GetResult())
}
