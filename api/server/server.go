package server

import (
	"io"

	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/pong/api/game"
	pb "github.com/nicholasjackson/pong/api/protos/pong"
)

// PongServer comment
type PongServer struct {
	logger       hclog.Logger
	serverClient pb.PongService_ClientStreamServer
}

// New creates a new server
func New(logger hclog.Logger) *PongServer {
	return &PongServer{logger: logger}
}

// ClientStream handles connections to and from the game
func (s *PongServer) ClientStream(stream pb.PongService_ClientStreamServer) error {
	s.logger.Info("Started client stream")
	s.serverClient = stream

	// start the game
	g := game.NewGame()

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
			"name", in.Name,
			"bat-x", in.X,
			"ball-y", in.Y)

		if in.Name == "RESET_GAME" {
			g.ResetGame()
		}

		if in.Name == "SERVE" {
			g.StartGame()
			go s.gameTick(g, stream)
		}

		if in.Name == "BAT_UP" {
			g.MoveBatUp(1)
		}

		if in.Name == "BAT_DOWN" {
			g.MoveBatDown(1)
		}

		stream.Send(g.DataAsProto())

		/*
			// forward the message to the other server
			err = s.apiClient.SendServer(int(in.Bat.X), int(in.Bat.Y), int(in.Ball.X), int(in.Ball.Y), in.Hit, int(in.Score)) // send data back
			if err != nil {
				s.logger.Error("Unable to send data to server", "error", err)
			}
		*/
	}
}

func (s *PongServer) gameTick(g *game.Game, stream pb.PongService_ClientStreamServer) {
	for _ = range g.Tick() {
		dp := g.DataAsProto()
		s.logger.Info("Send Data", "ball-x", dp.Ball.X, "ball-y", dp.Ball.Y)

		stream.Send(dp)
	}
}

// ServerStream handles connections to and from the server
func (s *PongServer) ServerStream(stream pb.PongService_ServerStreamServer) error {
	s.logger.Info("Started server stream")
	/*

		for {
			in, err := stream.Recv()
			if err == io.EOF {
				return nil
			}

			if err != nil {
				return err
			}

			s.logger.Info(
				"Got event", "handler", "server",
				"bat-x", in.Bat.X,
				"bat-y", in.Bat.Y,
				"ball-x", in.Ball.X,
				"ball-y", in.Ball.Y)

			// forward the message to the client
			if s.serverClient == nil {
				s.logger.Info("No client connected")
				continue
			}

			err = s.serverClient.Send(in)
			if err != nil {
				s.logger.Error("Unable to send message to client: %s", err)
			}
		}
	*/

	s.logger.Info("Disconnected client")
	return nil
}
