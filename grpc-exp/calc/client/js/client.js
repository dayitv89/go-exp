var grpc = require('grpc');
const Client = require('./calcpb/calc_grpc_pb');
const calc_pb = require('./calcpb/calc_pb');
const { Request } = calc_pb;

var client = new Client.CalculatorClient('localhost:50050', grpc.credentials.createInsecure());

const request = new Request();
request.setN1(10);
request.setN2(20);

client.sum(request, (_, response) => {
	console.log('sum', response.getResult());
});

client.subtract(request, (_, response) => {
	console.log('subtract', response.getResult());
});
