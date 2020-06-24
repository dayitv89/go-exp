#!/bin/bash 

## go grpc server
# protoc -I. --go_out=plugins=grpc,paths=source_relative:. calcpb/calc.proto

## go web Server
# protoc -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true,paths=source_relative:. calcpb/calc.proto

## go both grpc and web Server
protoc -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/ --go_out=plugins=grpc,paths=source_relative:. --grpc-gateway_out=logtostderr=true,paths=source_relative:. --swagger_out=logtostderr=true:. calcpb/calc.proto


# npm install -g grpc-tools
# client
# grpc_tools_node_protoc calcpb/calc.proto --js_out=import_style=commonjs,binary:./client/js --grpc_out=./client/js --plugin=protoc-gen-grpc=`which grpc_tools_node_protoc_plugin` 
# server
# grpc_tools_node_protoc calcpb/calc.proto --js_out=import_style=commonjs,binary:./server/js --grpc_out=./server/js --plugin=protoc-gen-grpc=`which grpc_tools_node_protoc_plugin` 