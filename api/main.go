package main

import (
	"log"
	"net"

	pb "github.com/nicholasjackson/pong/api/protos/pong"
	"github.com/nicholasjackson/pong/api/server"
	"google.golang.org/grpc"
)

func main() {
	grpcServer := grpc.NewServer()
	pb.RegisterPongServiceServer(grpcServer, &server.PongServer{})

	l, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Listening on tcp://localhost:6000")
	grpcServer.Serve(l)
}
