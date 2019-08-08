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

var p1s *objects.Score
var p2s *objects.Score

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
		bat1 = objects.NewBat(3, 0, 3, 6, tl.ColorRed, 0, true, batEventHandler)
		bat2 = objects.NewBat(3, 0, 3, 6, tl.ColorGreen, -3, false, nil)
		ball = objects.NewBall(6, 0, 3, 2, tl.ColorBlack, true, *player, ballEventHandler)
	} else {
		bat1 = objects.NewBat(3, 0, 3, 6, tl.ColorRed, 0, false, nil)
		bat2 = objects.NewBat(3, 0, 3, 6, tl.ColorGreen, -3, true, batEventHandler)
		ball = objects.NewBall(6, 0, 3, 2, tl.ColorBlack, false, *player, ballEventHandler)
	}

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
		resetGame()

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

func resetGame() {
	// reset the bat and ball position
	bat1.Reset()
	bat2.Reset()
	ball.Reset()
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

		if d.Score > 0 {
			scoreGame(d.Score)
			resetGame()
		}
	}
}
