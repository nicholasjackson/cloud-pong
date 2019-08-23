package client

import (
	"io"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/hashicorp/go-hclog"
	pb "github.com/nicholasjackson/pong/api/protos/pong"
)

// GameData something
type GameData struct {
	BatX  int
	BatY  int
	BallX int
	BallY int
	Hit   bool
	Score int
}

// Client nn
type Client struct {
	uri          string
	clientConn   *grpc.ClientConn
	client       pb.PongServiceClient
	clientStream pb.PongService_ClientStreamClient
	serverStream pb.PongService_ServerStreamClient
	outClient    chan GameData
	outServer    chan GameData
	Ready        bool
	logger       hclog.Logger
}

// New sss
func New(uri string, logger hclog.Logger) *Client {
	c := &Client{uri: uri, logger: logger}

	return c
}

// Dial syncronously dials the server
func (c *Client) Dial(isServer bool) error {
	c.outClient = make(chan GameData)
	c.outServer = make(chan GameData)

	c.logger.Info("Dialing connection", "server", isServer)

	var err error
	c.clientConn, err = grpc.Dial(c.uri, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		c.logger.Error("failed to connect", "error", err)
		return err
	}

	c.client = pb.NewPongServiceClient(c.clientConn)

	if !isServer {
		c.clientStream, err = c.client.ClientStream(context.Background())
		if err != nil {
			c.logger.Error("Error connecing to client", "error", err)
			return err
		}
	}

	if isServer {
		c.serverStream, err = c.client.ServerStream(context.Background())
		if err != nil {
			c.logger.Error("Error connecting to service", "errror", err)
			return err
		}
	}

	c.logger.Info("Client succesfully connected to the server")
	c.Ready = true
	return nil
}

// DialAsync dials the client, if an error occurs when dialing retries
func (c *Client) DialAsync(isServer bool) {
	go func() {
		c.Dial(isServer)
	}()
}

// SendClient something
func (c *Client) SendClient(batX, batY, ballX, ballY int, hit bool, score int) error {
	err := c.clientStream.Send(&pb.PongData{
		Bat:   &pb.Bat{X: int32(batX), Y: int32(batY)},
		Ball:  &pb.Ball{X: int32(ballX), Y: int32(ballY)},
		Hit:   hit,
		Score: int32(score),
	})

	// reconnect if the other server went away
	if err == io.EOF {
		c.logger.Error("Client disconnected, reconnecting")

		clientStream, err := c.client.ClientStream(context.Background())
		if err == nil {
			c.clientStream = clientStream
			return nil
		}

		c.logger.Error("Error connecting to service", "errror", err)
		return err
	}

	return nil
}

// RecieveClient something
func (c *Client) RecieveClient() chan GameData {
	go func() {
		for {
			resp, err := c.clientStream.Recv()
			if err != nil {
				c.logger.Error("Client error receiving stream", "error", err)
				close(c.outClient)
				return
			}

			c.outClient <- GameData{
				BatX:  int(resp.Bat.X),
				BatY:  int(resp.Bat.Y),
				BallX: int(resp.Ball.X),
				BallY: int(resp.Ball.Y),
				Hit:   resp.Hit,
				Score: int(resp.Score),
			}

		}

		c.logger.Error("Server stream terminated")
	}()

	return c.outClient
}

// SendServer something
func (c *Client) SendServer(batX, batY, ballX, ballY int, hit bool, score int) error {
	err := c.serverStream.Send(&pb.PongData{
		Bat:   &pb.Bat{X: int32(batX), Y: int32(batY)},
		Ball:  &pb.Ball{X: int32(ballX), Y: int32(ballY)},
		Hit:   hit,
		Score: int32(score),
	})

	// reconnect if the other server went away
	if err == io.EOF {
		c.logger.Error("Service disconnected, reconnecting")

		serverStream, err := c.client.ServerStream(context.Background())
		if err == nil {
			c.serverStream = serverStream
			return nil
		}

		c.logger.Error("Error connecting to service", "errror", err)
	}

	return err
}

// RecieveServer something
func (c *Client) RecieveServer() chan GameData {
	go func() {
		for {
			resp, err := c.serverStream.Recv()
			if err != nil {
				log.Fatal(err)
			}

			c.outServer <- GameData{BatX: int(resp.Bat.X), BatY: int(resp.Bat.Y)}
		}
	}()

	return c.outServer
}
