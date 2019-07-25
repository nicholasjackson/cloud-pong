package client

import (
	"log"

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
}

// Client nn
type Client struct {
	clientConn *grpc.ClientConn
	client     pb.PongServiceClient
	stream     pb.PongService_StreamClient
	out        chan GameData
}

// New sss
func New(uri string) *Client {
	c := &Client{}
	c.out = make(chan GameData)

	var err error
	c.clientConn, err = grpc.Dial(uri, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}

	c.client = pb.NewPongServiceClient(c.clientConn)
	c.stream, err = c.client.Stream(context.Background())
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	return c
}

// Send something
func (c *Client) Send(batX, batY, ballX, ballY int) {
	c.stream.Send(&pb.PongData{
		Bat:  &pb.Bat{X: int32(batX), Y: int32(batY)},
		Ball: &pb.Ball{X: int32(ballX), Y: int32(ballY)},
	})
}

// Recieve something
func (c *Client) Recieve() chan GameData {
	go func() {
		for {
			resp, err := c.stream.Recv()
			if err != nil {
				log.Fatal(err)
			}

			c.out <- GameData{BatX: int(resp.Bat.X), BatY: int(resp.Bat.Y)}
		}
	}()

	return c.out
}
