package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	pb "github.com/dayitv89/go-exp/grpc-exp/calc/calcpb"
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

func grpcServer() {
	lis, err := net.Listen("tcp", ":50050")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCalculatorServer(s, &grpcController{})
	fmt.Println("register for evans the reflection")
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	go grpcServer()
	if err := webServerMux(); err != nil {
		fmt.Println("webServerMux error", err)
	}
}

// type AppController struct{}

// func (ctrl *AppController) Ping(c *gin.Context) {
// 	msg := fmt.Sprintf("ping at the server: %d", time.Now().Unix())
// 	c.JSON(http.StatusOK, gin.H{"message": msg})
// }

// func (ctrl *AppController) Sum(c *gin.Context) {
// 	msg := fmt.Sprintf("ping at the server: %d", time.Now().Unix())
// 	c.JSON(http.StatusOK, gin.H{"message": msg})
// }

// func setupRoutes(r *gin.Engine) {
// 	// setup global middlewares
// 	r.Use(APIAuth())

// 	appCtrl := new(AppController)

// 	r.GET("/ping", appCtrl.Ping)
// 	r.GET("/sum", appCtrl.Sum)
// 	r.POST("/sum", appCtrl.Sum)
// }

// //IsJSONContent ...
// func IsJSONContent(c *gin.Context) bool {
// 	val, ok := c.Request.Header["Content-Type"]
// 	return ok && strings.Contains(val[0], "application/json")
// }

// //IsJSONRespond ...
// func IsJSONRespond(c *gin.Context) bool {
// 	if val, ok := c.Request.Header["Accept"]; ok {
// 		return strings.Contains(val[0], "application/json")
// 	}
// 	return IsJSONContent(c)
// }

// //APIAuth ...
// func APIAuth() func(c *gin.Context) {
// 	return func(c *gin.Context) {
// 		if IsJSONRespond(c) {
// 			c.Next()
// 		} else {
// 			c.JSON(http.StatusForbidden, gin.H{"message": "Invalid access to perform this action"})
// 			c.Abort()
// 		}
// 	}
// }

// func webServer() {
// 	r := gin.Default()
// 	http.DefaultClient.Timeout = 30 * time.Second
// 	setupRoutes(r)

// 	os.Setenv("PORT", "3002")
// 	fmt.Printf("\n\nRunning SERVER on port :%s and GIN_MODE=%s\n\n", os.Getenv("PORT"), os.Getenv("GIN_MODE"))
// 	r.Run(":3002" + os.Getenv("PORT"))
// 	return
// }

func webServerMux() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterCalculatorHandlerFromEndpoint(ctx, mux, "localhost:50050", opts)
	if err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(":3002", mux)
}
