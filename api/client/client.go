package client

import (
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

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
}

// New sss
func New(uri string) *Client {
	c := &Client{uri: uri}
	c.outClient = make(chan GameData)
	c.outServer = make(chan GameData)

	return c
}

// Dial syncronously dials the server
func (c *Client) Dial() {
	for {
		var err error
		c.clientConn, err = grpc.Dial(c.uri, grpc.WithInsecure())
		if err != nil {
			log.Printf("failed to connect: %s\n", err)
		}

		c.client = pb.NewPongServiceClient(c.clientConn)
		c.clientStream, err = c.client.ClientStream(context.Background())
		if err != nil {
			log.Printf("Error: %s\n", err)
		}

		c.serverStream, err = c.client.ServerStream(context.Background())
		if err != nil {
			log.Printf("Error: %s\n", err)
		}

		if err != nil {
			// an error has occured dialing wait then retry
			time.Sleep(1 * time.Second)
			continue
		}

		log.Println("Client succesfully connected to the server")

		break
	}
}

// DialAsync dials the client, if an error occurs when dialing retries
func (c *Client) DialAsync() {
	go func() {
		c.Dial()
	}()
}

// SendClient something
func (c *Client) SendClient(batX, batY, ballX, ballY int, hit bool, score int) {
	c.clientStream.Send(&pb.PongData{
		Bat:   &pb.Bat{X: int32(batX), Y: int32(batY)},
		Ball:  &pb.Ball{X: int32(ballX), Y: int32(ballY)},
		Hit:   hit,
		Score: int32(score),
	})
}

// RecieveClient something
func (c *Client) RecieveClient() chan GameData {
	go func() {
		for {
			resp, err := c.clientStream.Recv()
			if err != nil {
				log.Fatal(err)
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
	}()

	return c.outClient
}

// SendServer something
func (c *Client) SendServer(batX, batY, ballX, ballY int, hit bool, score int) {
	c.serverStream.Send(&pb.PongData{
		Bat:   &pb.Bat{X: int32(batX), Y: int32(batY)},
		Ball:  &pb.Ball{X: int32(ballX), Y: int32(ballY)},
		Hit:   hit,
		Score: int32(score),
	})
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
