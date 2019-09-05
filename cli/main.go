package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	tl "github.com/JoelOtter/termloop"
	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"
	pb "github.com/nicholasjackson/pong/api/protos/pong"
	"github.com/nicholasjackson/pong/objects"
	"google.golang.org/grpc"
)

var player = env.Int("PLAYER", false, 1, "Player number")
var apiURI = env.String("API_URI", false, "localhost:6000", "URI for the api server")
var singleKeyboard = env.Bool("SINGLE_KEYBOARD", false, false, "Control both bats from a single terminal")

var bat1 *objects.Bat
var bat2 *objects.Bat
var ball *objects.Ball

var p1s *objects.Score
var p2s *objects.Score

var logger hclog.Logger

var client pb.PongService_ClientStreamClient

func main() {
	env.Parse()

	f, err := os.Create(fmt.Sprintf("%d_out.log", *player))
	if err != nil {
		log.Fatal(err)
	}
	opt := &hclog.LoggerOptions{Output: f}
	logger = hclog.New(opt)

	logger.Info("Starting client", "player", *player, "uri", *apiURI)

	// setup monitoring for inbound events
	go streamReceive()

	if *player == 1 {
		bat1 = objects.NewBat(0, 0, 0, 0, tl.ColorRed, 1, *singleKeyboard, handler)
		bat2 = objects.NewBat(0, 0, 0, 0, tl.ColorGreen, 2, false, nil)
	}
	if *player == 2 {
		bat1 = objects.NewBat(0, 0, 0, 0, tl.ColorRed, 1, false, nil)
		bat2 = objects.NewBat(0, 0, 0, 0, tl.ColorGreen, 2, false, handler)
	}

	ball = objects.NewBall(0, 0, 0, 0, tl.ColorBlack)

	// create the net
	net := objects.NewNet(tl.ColorBlack)

	// create player1 score
	p1s = objects.NewScore(-14, 3, tl.ColorBlack)
	p2s = objects.NewScore(3, 3, tl.ColorBlack)

	g := tl.NewGame()
	g.Screen().SetFps(60)
	l := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorWhite,
	})

	// add the bats
	l.AddEntity(bat1)
	l.AddEntity(bat2)

	// add the ball
	l.AddEntity(ball)

	// add the net
	l.AddEntity(net)

	// add the scores
	l.AddEntity(p1s)
	l.AddEntity(p2s)

	g.Screen().SetLevel(l)
	g.Screen().AddEntity(tl.NewFpsText(0, 0, tl.ColorRed, tl.ColorDefault, 0.5))

	g.Start()
}

func streamReceive() {
	conn, err := grpc.Dial(*apiURI, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := pb.NewPongServiceClient(conn)

	client, err = c.ClientStream(context.Background())
	if err != nil {
		panic(err)
	}

	// as soon as connected send a reset game command
	if *player == 1 {
		client.Send(&pb.Event{Name: "RESET_GAME"})
	}

	for {
		d, err := client.Recv()
		if err == io.EOF {
			log.Fatal("Stream closed by server")
		}
		if err != nil {
			log.Fatal(err)
		}

		// draw the objects
		bat1.SetData(d.Bat1.X, d.Bat1.Y, d.Bat1.W, d.Bat1.H, d.Game.W, d.Game.H)
		bat2.SetData(d.Bat2.X, d.Bat2.Y, d.Bat2.W, d.Bat2.H, d.Game.W, d.Game.H)
		ball.SetData(d.Ball.X, d.Ball.Y, d.Ball.W, d.Ball.H, d.Game.W, d.Game.H)
		p1s.UpdateScore(d.Player1Score)
		p2s.UpdateScore(d.Player2Score)
	}
}

func handler(message string) {
	logger.Info("sending", "message", message)
	client.Send(&pb.Event{Name: message})
}

/*
func batEventHandler(e interface{}) {
	logger.Info("Send bat pos to server")
	batPos := e.(*objects.BatMoveEvent)

	c.SendClient(batPos.X, batPos.Y, 0, 0, false, 0)
}

func ballEventHandler(e interface{}) {
	switch ev := e.(type) {
	case *objects.BallMoveEvent:
		logger.Info("Send ball pos to server")
		ballPos := e.(*objects.BallMoveEvent)
		c.SendClient(0, 0, ballPos.X, ballPos.Y, false, 0)
	case *objects.BallHitEvent:
		logger.Info("Collided")
		c.SendClient(0, 0, 0, 0, true, 0)
	case *objects.BallScoreEvent:
		scoreGame(ev.Player)
		resetGame(ev.Player)

		c.SendClient(0, 0, 0, 0, false, ev.Player)
	}
}

func scoreGame(player int) {
	if player == 1 {
		p1s.IncrementScore()
		return
	}

	p2s.IncrementScore()
}

func resetGame(player int) {
	// reset the bat and ball position
	bat1.Reset()
	bat2.Reset()
	ball.Reset(player)
}

func streamReceive() {
	for d := range c.RecieveClient() {
		logger.Info("move", "event", d)

		if d.BatX != 0 || d.BatY != 0 {
			if *player == 1 {
				bat2.SetPos(-3, d.BatY)
			} else {
				if d.BatY != 0 {
					bat1.SetPos(3, d.BatY)
				}
			}
		}

		if (d.BallX != 0 || d.BallY != 0) && !ball.IsControlled() {
			ball.SetPos(d.BallX, d.BallY)
		}

		// figure out if we have lost control
		if d.Hit {
			// remove control
			ball.SetControl(false)
		}

		// if the score has been updated reset the game
		if d.Score > 0 {
			scoreGame(d.Score)
			resetGame(d.Score)
		}
	}

	logger.Error("Server disconnected, restarting")
	err := c.Dial(false)
	if err != nil {
		panic(err)
	}

	logger.Error("Server reconnected")
	streamReceive()
}
*/
