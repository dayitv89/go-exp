package srcgo

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/dayitv89/go-exp/grpc-exp/calc/gen/calcpb"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcController struct{}

func (*grpcController) Sum(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	return &pb.Response{Result: req.GetN1() + req.GetN2()}, nil
}

func (*grpcController) Subtract(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	return &pb.Response{Result: req.GetN1() - req.GetN2()}, nil
}

// GRPCServer ...
func GRPCServer(grpcPort string) {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCalculatorServer(s, &grpcController{})
	reflection.Register(s)
	fmt.Println("grpc server is running on", grpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	fmt.Println("grpc server closed from", grpcPort)
}

// grpc gateway to call this method
func setupMux(grpcPort string) (*runtime.ServeMux, context.CancelFunc, error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterCalculatorHandlerFromEndpoint(ctx, mux, "127.0.0.1"+grpcPort, opts)
	return mux, cancel, err
}
