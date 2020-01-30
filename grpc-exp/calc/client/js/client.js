var grpc = require('grpc');
const Service = require('./calcpb/calc_grpc_pb');
const messages = require('./calcpb/calc_pb');

var client = new Service.CalculatorClient('localhost:50050', grpc.credentials.createInsecure());

const request = new messages.Request();
request.setN1(10);
request.setN2(20);

client.sum(request, (err, response) => {
	console.log('sum', 'error:', err, 'res:', response && response.getResult());
});

client.subtract(request, (err, response) => {
	console.log('subtract', 'error:', err, 'res:', response && response.getResult());
});
