package srcgo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

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

func WebServerGin(grpcPort string) error {
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

	mux, cancel, err := setupMux(grpcPort)
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
