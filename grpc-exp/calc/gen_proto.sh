#!/bin/bash 

# rm -rf gen && mkdir -p gen
# cp -f proto/calc.proto gen/calc.proto

## go grpc-server
# protoc proto/calc.proto -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway --go_out=plugins=grpc,paths=source_relative:./gen

## go api-gateway
# protoc proto/calc.proto -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway --grpc-gateway_out=logtostderr=true,paths=source_relative:./gen

## swagger-json
# protoc proto/calc.proto -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway --swagger_out=logtostderr=true:./gen

## go all 3 (grpc-server, api-gateway & swagger-json)
# protoc proto/calc.proto -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway --go_out=plugins=grpc,paths=source_relative:./gen --grpc-gateway_out=logtostderr=true,paths=source_relative:./gen --swagger_out=logtostderr=true:./gen

# mv gen/proto gen/calcpb
# mv gen/calc.proto gen/calcpb/calc.proto


##### Docker way (recommended)
rm -rf gen/calcpb && mkdir -p gen/calcpb
cp -f proto/calc.proto gen/calcpb/calc.proto
docker run --rm -v `pwd`/gen/calcpb/:/defs namely/protoc-all:1.29_2 -f ./calc.proto -l go --with-gateway -o . --go-source-relative .

# npm install -g grpc-tools
# client
# grpc_tools_node_protoc calcpb/calc.proto --js_out=import_style=commonjs,binary:./client/js --grpc_out=./client/js --plugin=protoc-gen-grpc=`which grpc_tools_node_protoc_plugin` 
# server
# grpc_tools_node_protoc calcpb/calc.proto --js_out=import_style=commonjs,binary:./server/js --grpc_out=./server/js --plugin=protoc-gen-grpc=`which grpc_tools_node_protoc_plugin` 
