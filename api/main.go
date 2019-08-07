package main

import (
	"fmt"
	"net"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"
	"github.com/nicholasjackson/pong/api/client"
	pb "github.com/nicholasjackson/pong/api/protos/pong"
	"github.com/nicholasjackson/pong/api/server"
	"google.golang.org/grpc"
)

var player = env.Int("PLAYER", false, 1, "Player number")
var port = env.Int("BIND_PORT", false, 6000, "Bind port for server")
var upstream = env.String("UPSTREAM_ADDRESS", false, "localhost:6001", "Upstream address for other server")

var logger hclog.Logger

func main() {
	env.Parse()

	logger = hclog.Default()
	apiClient := client.New(*upstream)
	apiClient.DialAsync()

	grpcServer := grpc.NewServer()
	server := server.New(logger, apiClient)
	pb.RegisterPongServiceServer(grpcServer, server)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		logger.Error("Failed to listen", "error", err)
		os.Exit(1)
	}

	logger.Info("Listening on port", "port", *port, "player", *player)
	grpcServer.Serve(l)
}
