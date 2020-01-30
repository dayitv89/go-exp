#!/bin/bash 

protoc calcpb/calc.proto --go_out=plugins=grpc:.

# npm install -g grpc-tools
grpc_tools_node_protoc calcpb/calc.proto --js_out=import_style=commonjs,binary:./client/js --grpc_out=./client/js --plugin=protoc-gen-grpc=`which grpc_tools_node_protoc_plugin` 