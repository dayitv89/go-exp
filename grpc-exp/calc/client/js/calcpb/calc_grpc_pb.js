// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var calcpb_calc_pb = require('../calcpb/calc_pb.js');

function serialize_calcpb_Request(arg) {
  if (!(arg instanceof calcpb_calc_pb.Request)) {
    throw new Error('Expected argument of type calcpb.Request');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_calcpb_Request(buffer_arg) {
  return calcpb_calc_pb.Request.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_calcpb_Response(arg) {
  if (!(arg instanceof calcpb_calc_pb.Response)) {
    throw new Error('Expected argument of type calcpb.Response');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_calcpb_Response(buffer_arg) {
  return calcpb_calc_pb.Response.deserializeBinary(new Uint8Array(buffer_arg));
}


var CalculatorService = exports.CalculatorService = {
  sum: {
    path: '/calcpb.Calculator/Sum',
    requestStream: false,
    responseStream: false,
    requestType: calcpb_calc_pb.Request,
    responseType: calcpb_calc_pb.Response,
    requestSerialize: serialize_calcpb_Request,
    requestDeserialize: deserialize_calcpb_Request,
    responseSerialize: serialize_calcpb_Response,
    responseDeserialize: deserialize_calcpb_Response,
  },
  subtract: {
    path: '/calcpb.Calculator/Subtract',
    requestStream: false,
    responseStream: false,
    requestType: calcpb_calc_pb.Request,
    responseType: calcpb_calc_pb.Response,
    requestSerialize: serialize_calcpb_Request,
    requestDeserialize: deserialize_calcpb_Request,
    responseSerialize: serialize_calcpb_Response,
    responseDeserialize: deserialize_calcpb_Response,
  },
};

exports.CalculatorClient = grpc.makeGenericClientConstructor(CalculatorService);
