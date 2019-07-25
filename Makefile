protos:
	protoc -I api/protos/ api/protos/api.proto --go_out=plugins=grpc:api/protos/pong
