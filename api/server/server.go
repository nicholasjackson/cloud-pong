package server

import (
	"context"
	"io"

	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/pong/api/game"
	pb "github.com/nicholasjackson/pong/api/protos/pong"
	"google.golang.org/grpc"
)

// PongServer comment
type PongServer struct {
	logger       hclog.Logger
	playerClient pb.PongService_ClientStreamServer

	serverURI      string
	upstreamClient pb.PongService_ServerStreamClient
	upstreamServer pb.PongService_ServerStreamServer
	player         int
	g              *game.Game
}

// New creates a new server
func New(logger hclog.Logger, player int, serverURI string) *PongServer {
	return &PongServer{logger: logger, player: player, serverURI: serverURI}
}

// ClientStream handles connections to and from the game
func (s *PongServer) ClientStream(stream pb.PongService_ClientStreamServer) error {
	s.logger.Info("Started client stream")
	s.playerClient = stream

	if s.player == 1 {
		s.g = game.NewGame() // start a new game
	}

	// open a connection to the other API
	conn, err := grpc.Dial(s.serverURI, grpc.WithInsecure())
	if err != nil {
		return err
	}

	defer conn.Close()
	c := pb.NewPongServiceClient(conn)
	s.upstreamClient, err = c.ServerStream(context.Background())
	if err != nil {
		s.logger.Error("Unable to connect to upstream server")
		return err
	}

	// handle messages from the upstream server
	go func() {
		for {
			d, err := s.upstreamClient.Recv()
			if err == io.EOF {
				s.logger.Error("Upstream disconnected", "error", err)
				break
			}
			if err != nil {
				s.logger.Error("Upstream disconnected", "error", err)
				break
			}

			// got message forward it
			s.logger.Info("Received from Server", "Message", d)
			stream.Send(d)
		}
	}()

	// handle messages from the pong client
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			s.logger.Error("Close received from player", "error", err)
			break
		}
		if err != nil {
			s.logger.Error("Close received from player", "error", err)
			break
		}

		s.logger.Info(
			"Got event", "handler", "client",
			"name", in.Name,
			"bat-x", in.X,
			"ball-y", in.Y)

		if s.player == 1 {
			s.handlePlayer1(in.Name)
		} else {
			s.handlePlayer2(in.Name)
		}
	}

	s.playerClient = nil

	if s.player == 1 {
		s.logger.Info("Player went away")
		// quit game
		s.g.ResetGame()
	}

	return nil
}

func (s *PongServer) handlePlayer1(m string) {
	if m == "RESET_GAME" {
		s.g.ResetGame()
	}

	if m == "SERVE" {
		s.g.StartGame()
		go s.gameTick()
	}

	if m == "BAT_UP" {
		s.g.MoveBatUp(1)
	}

	if m == "BAT_DOWN" {
		s.g.MoveBatDown(1)
	}

	data := s.g.DataAsProto()

	// send data to the other player
	if s.upstreamServer != nil {
		s.upstreamServer.Send(data)
	} else {
	}

	// send data back to the client
	s.playerClient.Send(data)
}

func (s *PongServer) handlePlayer2(m string) {
	// forward the data to the other server
	s.upstreamClient.Send(&pb.Event{Name: m})
}

func (s *PongServer) gameTick() {
	for _ = range s.g.Tick() {
		dp := s.g.DataAsProto()
		s.logger.Info("Send Data", "ball-x", dp.Ball.X, "ball-y", dp.Ball.Y)

		if s.upstreamServer == nil {
			s.logger.Info("Upstream server not connected")
		} else {
			s.upstreamServer.Send(dp)
		}

		s.playerClient.Send(dp)
	}
}

// ServerStream handles connections to and from the server
func (s *PongServer) ServerStream(stream pb.PongService_ServerStreamServer) error {
	s.logger.Info("Started server stream")
	s.upstreamServer = stream

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		s.logger.Info(
			"Got event", "name", in.Name,
			"bat-x", in.X,
			"bat-y", in.Y,
		)

		/*
			// forward the message to the client
			if s.serverClient == nil {
				s.logger.Info("No client connected")
				continue
			}

			err = s.serverClient.Send(in)
			if err != nil {
				s.logger.Error("Unable to send message to client: %s", err)
			}
		*/
	}

	s.logger.Info("Disconnected client")
	return nil
}
