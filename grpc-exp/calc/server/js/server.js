var grpc = require('grpc');
const services = require('./calcpb/calc_grpc_pb');
const messages = require('./calcpb/calc_pb');

function sum(call, callback) {
	const response = new messages.Response();
	response.setResult(call.request.getN1() + call.request.getN2());
	callback(null, response);
}
function subtract(call, callback) {
	const response = new messages.Response();
	response.setResult(call.request.getN1() - call.request.getN2());
	callback(null, response);
}

var server = new grpc.Server();
server.addService(services.CalculatorService, { sum, subtract });
server.bind('0.0.0.0:50050', grpc.ServerCredentials.createInsecure());
server.start();
