# go-lang experiments

### grpc calculator experiments

It has a demo of grpc implement of simple calculator with sum and subtract method.

Implemented:

- grpc-exp/calc/calcpb/calc.proto
- grpc server in golang
- grpc client in golang and nodejs

Additional Implemented:

- [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) for golang server
- swagger json file as available in [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway).
- grpc reverser proxy server in `runtime.ServeMux`
- **grpc reverser proxy server in [`gin gonic`](https://github.com/gin-gonic/gin).**

**\*why not default implementation?** Because, [gin gonic](https://github.com/gin-gonic/gin) supports multiple way to interact with client. My idea is to use grpc with my existing server and my existing server gives me a tons of features like http2 push, SSE, already implemented whole monolithic application on gin gonic and additionally control to pass to grpc like hijeck or middleware. So just want to plug grpc service to my http client.

BTW [default example](https://github.com/gin-gonic/examples/tree/master/grpc) for writing individual client is so much time taken and TDS job for me, so I just want to a wildcard option for that.

#### Setup protobuf:

```sh
go install \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
    github.com/golang/protobuf/protoc-gen-go

cd grpc-exp/calc
chmod +x gen_proto.sh
./gen_proto.sh

```

#### Run Server

```
cd grpc-exp/calc
go run server/server.go
```

#### Demo

- `GET /sum`

```sh
# gin gonic
curl --location --request GET 'http://localhost:3000/grpc/sum?n1=12&n2=10'

# runtime.ServeMux (default)
curl --location --request GET 'http://localhost:3001/sum?n1=12&n2=10'
```

Output:

```json
// gin gonic
{
    "code": 200,
    "grpc_res": {
        "result": 22
    }
}

// runtime.ServeMux (default)
{
    "result": 22
}
```

**As Error:**

```sh
# gin gonic
curl --location --request GET 'http://localhost:3000/grpc/this_method_does_not_exist?n1=12&n2=10'

# runtime.ServeMux (default)
curl --location --request GET 'http://localhost:3001/this_method_does_not_exist?n1=12&n2=10'
```

Output:

```json
// gin gonic
{
    "code": 404,
    "error": "cannot convert to json string, use status code or grpc_res to find the issue.",
    "error_detail": "invalid character 'N' looking for beginning of value",
    "grpc_res": "Not Found\n"
}

// runtime.ServeMux (default)
Not Found\n
```

- `POST /subtract`

```sh
curl --location --request POST 'http://localhost:3000/grpc/subtract' \
--header 'Content-Type: application/json' \
--data-raw '{
    "n1": 12,
    "n2": 10
}'
```

```json
{
	"code": 200,
	"grpc_res": {
		"result": 2
	}
}
```

### Server file

```go
//grpc-exp/calc/server/server.go
// run as $ go run grpc-exp/calc/server/server.go
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

	pb "github.com/dayitv89/go-exp/grpc-exp/calc/calcpb"
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

const grpcPort string = ":50051"

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

```
