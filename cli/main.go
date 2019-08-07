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
		bat2 = objects.NewBat(3, 4, 3, 6, tl.ColorGreen, -3, false, nil)
	} else {
		bat1 = objects.NewBat(3, 3, 3, 6, tl.ColorRed, 0, false, nil)
		bat2 = objects.NewBat(3, 4, 3, 6, tl.ColorGreen, -3, true, batEventHandler)
	}

	g := tl.NewGame()
	g.Screen().SetFps(60)
	l := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorWhite,
	})

	l.AddEntity(bat1)
	l.AddEntity(bat2)

	g.Screen().SetLevel(l)
	g.Screen().AddEntity(tl.NewFpsText(0, 0, tl.ColorRed, tl.ColorDefault, 0.5))

	g.Start()
}

func batEventHandler(e interface{}) {
	logger.Info("Send pos to server")
	batPos := e.(*objects.BatMoveEvent)

	c.SendClient(batPos.X, batPos.Y, 0, 0)
}

func streamReceive() {
	for d := range c.RecieveClient() {

		if *player == 1 {
			logger.Info("move bat 2")
			bat2.SetPos(-3, d.BatY)
		} else {
			logger.Info("move bat 1")
			bat1.SetPos(3, d.BatY)
		}

		logger.Info("something", "y", d.BatY)
	}
}
