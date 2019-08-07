package server

import (
	"io"

	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/pong/api/client"
	pb "github.com/nicholasjackson/pong/api/protos/pong"
)

// PongServer comment
type PongServer struct {
	logger       hclog.Logger
	apiClient    *client.Client
	serverClient pb.PongService_ClientStreamServer
}

// New creates a new server
func New(logger hclog.Logger, apiClient *client.Client) *PongServer {
	return &PongServer{logger: logger, apiClient: apiClient}
}

// ClientStream handles connections to and from the game
func (s *PongServer) ClientStream(stream pb.PongService_ClientStreamServer) error {
	s.logger.Info("Started client stream")
	s.serverClient = stream

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		s.logger.Info(
			"Got event", "handler", "client",
			"bat-x", in.Bat.X,
			"bat-y", in.Bat.Y,
			"ball-x", in.Ball.X,
			"ball-y", in.Ball.Y)

		// forward the message to the other server
		s.apiClient.SendServer(int(in.Bat.X), int(in.Bat.Y), int(in.Ball.X), int(in.Ball.Y), in.Hit) // send data back
	}
}

// ServerStream handles connections to and from the server
func (s *PongServer) ServerStream(stream pb.PongService_ServerStreamServer) error {
	s.logger.Info("Started server stream")

	for {
		in, err := stream.Recv()
		s.logger.Error("Error reading stream", "error", err)
		if err == io.EOF {
			return nil
		}

		if err != nil {
			s.logger.Error("Error reading stream", "error", err)
			return err
		}

		s.logger.Info(
			"Got event", "handler", "server",
			"bat-x", in.Bat.X,
			"bat-y", in.Bat.Y,
			"ball-x", in.Ball.X,
			"ball-y", in.Ball.Y)

		// forward the message to the client
		s.serverClient.Send(in)
	}
}
