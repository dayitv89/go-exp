#!/bin/bash 

protoc calcpb/calc.proto --go_out=plugins=grpc:.