version=v0.2.0

protos:
	protoc -I api/protos/ api/protos/api.proto --go_out=plugins=grpc:api/protos/pong

start-server-1:
	(cd api && PLAYER=1 BIND_PORT=6000 UPSTREAM_ADDRESS=localhost:6001 go run main.go)

start-server-2:
	(cd api && PLAYER=2 BIND_PORT=6001 UPSTREAM_ADDRESS=localhost:6000 go run main.go)
	
player-1:
	(cd cli && PLAYER=1 API_URI=localhost:6000 go run main.go)

player-1-tf:
	(cd cli && PLAYER=1 API_URI=$(shell cd terraform && terraform output aks_pong_addr):6000 go run main.go)

player-2:
	(cd cli && PLAYER=2 API_URI=localhost:6001 go run main.go)

player-2-tf:
	(cd cli && PLAYER=2 API_URI=$(shell cd terraform && terraform output vms_pong_addr):6000 go run main.go)

pong-api-linux:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./bin/pong-api ./api/main.go

pong-cli-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./bin/pong-cli-linux ./cli/main.go

pong-cli-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o ./bin/pong-cli-mac ./cli/main.go

pong-cli-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o ./bin/pong-cli-windows.exe ./cli/main.go

pong-cli: pong-cli-linux pong-cli-mac pong-cli-windows

docker-java:
	cd java && docker build -t nicholasjackson/cloud-pong-api:java-${version} .
	docker push nicholasjackson/cloud-pong-api:java-${version}

docker-go: pong-api-linux
	docker build -t nicholasjackson/cloud-pong-api:go-${version} -f api/Dockerfile .
	docker push nicholasjackson/cloud-pong-api:go-${version}

docker-all: docker-go docker-java
