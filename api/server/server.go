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

		s.logger.Info("Got bat", "handler", "client", "x", in.Bat.X, "y", in.Bat.Y)

		// forward the message to the other server
		s.apiClient.SendServer(int(in.Bat.X), int(in.Bat.Y), 0, 0) // send data back
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

		s.logger.Info("Got bat", "handler", "server", "x", in.Bat.X, "y", in.Bat.Y)

		// forward the message to the client
		//s.apiClient.SendServer(int(in.Bat.X), int(in.Bat.Y), 0, 0) // send data back
		s.serverClient.Send(&pb.PongData{Bat: &pb.Bat{X: in.Bat.X, Y: in.Bat.Y}}) // send data back
	}
}
