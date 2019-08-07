package main

import (
	"fmt"
	"log"
	"os"

	tl "github.com/JoelOtter/termloop"
	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"
	"github.com/nicholasjackson/pong/api/client"
	"github.com/nicholasjackson/pong/objects"
)

var c *client.Client

var player = env.Int("PLAYER", false, 1, "Player number")
var apiURI = env.String("API_URI", false, "localhost:6000", "URI for the api server")

var bat1 *objects.Bat
var bat2 *objects.Bat
var ball *objects.Ball

var logger hclog.Logger

func main() {
	env.Parse()

	f, err := os.Create(fmt.Sprintf("%d_out.log", *player))
	if err != nil {
		log.Fatal(err)
	}
	opt := &hclog.LoggerOptions{Output: f}
	logger = hclog.New(opt)

	logger.Info("Starting client", "player", *player, "uri", *apiURI)

	c = client.New(*apiURI)
	c.Dial()

	// setup monitoring for inbound events
	go streamReceive()

	if *player == 1 {
		bat1 = objects.NewBat(3, 3, 3, 6, tl.ColorRed, 0, true, batEventHandler)
		bat2 = objects.NewBat(3, 3, 3, 6, tl.ColorGreen, -3, false, nil)
		ball = objects.NewBall(6, 5, 3, 2, tl.ColorBlack, true, *player, ballEventHandler)
	} else {
		bat1 = objects.NewBat(3, 3, 3, 6, tl.ColorRed, 0, false, nil)
		bat2 = objects.NewBat(3, 3, 3, 6, tl.ColorGreen, -3, true, batEventHandler)
		ball = objects.NewBall(6, 5, 3, 2, tl.ColorBlack, false, *player, ballEventHandler)
	}

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

	g.Screen().SetLevel(l)
	g.Screen().AddEntity(tl.NewFpsText(0, 0, tl.ColorRed, tl.ColorDefault, 0.5))

	g.Start()
}

func batEventHandler(e interface{}) {
	logger.Info("Send bat pos to server")
	batPos := e.(*objects.BatMoveEvent)

	c.SendClient(batPos.X, batPos.Y, 0, 0, false)
}

func ballEventHandler(e interface{}) {
	switch e.(type) {
	case *objects.BallMoveEvent:
		logger.Info("Send ball pos to server")
		ballPos := e.(*objects.BallMoveEvent)
		c.SendClient(0, 0, ballPos.X, ballPos.Y, false)
	case *objects.BallHitEvent:
		logger.Info("Collided")
		c.SendClient(0, 0, 0, 0, true)
	case *objects.BallScoreEvent:
		resetGame()
	}
}

func resetGame() {
	panic("Reset")
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

		// figure out if we have lost control
		if (d.BallX != 0 || d.BallY != 0) && !ball.IsControlled() {
			ball.SetPos(d.BallX, d.BallY)
		}

		if d.Hit {
			// remove control
			ball.SetControl(false)
		}
	}
}
