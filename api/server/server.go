package server

import (
	"io"
	"log"

	pb "github.com/nicholasjackson/pong/api/protos/pong"
)

// PongServer comment
type PongServer struct {
}

// Stream comment
func (s *PongServer) Stream(stream pb.PongService_StreamServer) error {
	log.Println("Started stream")

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("Got bat %d,%d", in.Bat.X, in.Bat.Y)

		stream.Send(&pb.PongData{Bat: &pb.Bat{X: in.Bat.X, Y: in.Bat.Y}}) // send data back
	}
}
