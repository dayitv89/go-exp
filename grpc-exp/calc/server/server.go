package main

import (
	"context"
	"log"
	"net"

	pb "github.com/dayitv89/go-exp/grpc-exp/calc/calcpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	return &pb.Response{Result: req.GetN1() + req.GetN2()}, nil
}

func (*server) Subtract(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	return &pb.Response{Result: req.GetN1() - req.GetN2()}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50050")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCalculatorServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
