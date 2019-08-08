protos:
	protoc -I api/protos/ api/protos/api.proto --go_out=plugins=grpc:api/protos/pong

start-servers:
	#(cd api && PLAYER=1 BIND_PORT=6000 UPSTREAM_ADDRESS=localhost:6001 go run main.go)
	(cd api && PLAYER=2 BIND_PORT=6001 UPSTREAM_ADDRESS=localhost:6000 go run main.go &)
	
start-client:
	(cd cli && PLAYER=2 API_URI=localhost:6001 go run main.go)