protos:
	protoc -I api/protos/ api/protos/api.proto --go_out=plugins=grpc:api/protos/pong

start-server-1:
	(cd api && PLAYER=1 BIND_PORT=6000 UPSTREAM_ADDRESS=localhost:6001 go run main.go)
	
player-1:
	(cd cli && PLAYER=1 API_URI=localhost:6000 go run main.go)

player-2:
	(cd cli && PLAYER=2 API_URI=localhost:6001 go run main.go)