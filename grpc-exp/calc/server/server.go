package main

import (
	"fmt"

	srcgo "github.com/dayitv89/go-exp/grpc-exp/calc/server/srcgo"
)

const grpcPort string = ":50050"

func main() {
	go srcgo.GRPCServer(grpcPort)
	go srcgo.WebServerMux(grpcPort)

	if err := srcgo.WebServerGin(grpcPort); err != nil {
		fmt.Println("webServerMux error", err)
	}
}
