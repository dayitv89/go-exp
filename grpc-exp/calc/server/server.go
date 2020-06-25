package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	pb "github.com/dayitv89/go-exp/grpc-exp/calc/gen/calcpb"
	"github.com/gin-gonic/gin"
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

const grpcPort string = ":50050"

func grpcServer() {
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

func main() {
	go grpcServer()
	go webServerMux()

	if err := webServerGin(); err != nil {
		fmt.Println("webServerMux error", err)
	}
}

//GRPCResponseHandler hijeck the response body
type GRPCResponseHandler struct {
	gin.ResponseWriter
	body       []byte
	statusCode int
}

//NewGRPCResponseHandler get new instance for each request
func NewGRPCResponseHandler(w gin.ResponseWriter) *GRPCResponseHandler {
	return &GRPCResponseHandler{w, []byte{}, http.StatusOK}
}

//Write stop the data flush and pushing to the client as we need to override the body at the end from grpc response.
func (grh *GRPCResponseHandler) Write(data []byte) (int, error) {
	grh.body = append(grh.body, data...)
	return 0, nil
}

//WriteHeader Holds the statusCode
func (grh *GRPCResponseHandler) WriteHeader(code int) {
	grh.statusCode = code
	grh.ResponseWriter.WriteHeader(code)
}

func webServerGin() error {
	r := gin.Default()
	http.DefaultClient.Timeout = 30 * time.Second

	// setup global middlewares
	r.Use(func(c *gin.Context) {
		c.Next()
		fmt.Println("CUSTOM LOG:\nRequest:", c.Request, "\nResponse:", c.Writer.Status())
	})

	r.GET("/ping", func(c *gin.Context) {
		msg := fmt.Sprintf("ping at the server: %d", time.Now().Unix())
		c.JSON(http.StatusOK, gin.H{"message": msg})
	})

	mux, cancel, err := setupMux()
	defer cancel()
	if err != nil {
		return err
	}
	r.Any("/grpc/*path", func(c *gin.Context) {
		path := c.Param("path")
		c.Status(http.StatusOK)

		c.Request.URL.Path = path
		fmt.Println("\ngin mux grpc request", c.Request.URL)
		grpcResponse := NewGRPCResponseHandler(c.Writer)
		mux.ServeHTTP(grpcResponse, c.Request)

		var jsonS interface{}
		c.Header("Content-Type", "application/json")
		if err := json.Unmarshal(grpcResponse.body, &jsonS); err != nil {
			c.JSON(grpcResponse.statusCode, gin.H{
				"code":         grpcResponse.statusCode,
				"error":        "cannot convert to json string, use status code or grpc_res to find the issue.",
				"error_detail": err.Error(),
				"grpc_res":     string(grpcResponse.body),
			})
		} else {
			c.JSON(grpcResponse.statusCode, gin.H{"code": grpcResponse.statusCode, "grpc_res": jsonS})
		}
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "error": "route not available"})
	})

	os.Setenv("PORT", "3000")
	fmt.Printf("\n\nRunning SERVER on port :%s and GIN_MODE=%s\n\n", os.Getenv("PORT"), os.Getenv("GIN_MODE"))

	return r.Run(":" + os.Getenv("PORT"))
}

func setupMux() (*runtime.ServeMux, context.CancelFunc, error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterCalculatorHandlerFromEndpoint(ctx, mux, "127.0.0.1"+grpcPort, opts)
	return mux, cancel, err
}

func webServerMux() error {
	mux, cancel, err := setupMux()
	defer cancel()
	if err != nil {
		return err
	}
	return http.ListenAndServe(":3001", mux)
}
