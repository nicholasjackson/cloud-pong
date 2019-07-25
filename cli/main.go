package main

import (
	"log"
	"os"

	tl "github.com/JoelOtter/termloop"
	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"
	"github.com/nicholasjackson/pong/api/client"
)

var c *client.Client

var player = env.Int("PLAYER", false, 1, "Player number")

var bat1 *CollBat
var bat2 *CollBat

var logger hclog.Logger

// CollBat defines a bat in the game
type CollBat struct {
	*tl.Rectangle
	move        bool
	px          int
	py          int
	offsetRight int
}

// NewCollBat comment
func NewCollBat(x, y, w, h int, color tl.Attr, move bool, offsetRight int) *CollBat {
	return &CollBat{
		Rectangle:   tl.NewRectangle(x, y, w, h, color),
		move:        move,
		px:          x,
		py:          y,
		offsetRight: offsetRight,
	}
}

// Draw something
func (r *CollBat) Draw(s *tl.Screen) {
	if r.offsetRight < 0 {
		sx, _ := s.Size()
		bx, _ := r.Size()

		r.px = sx - (bx - r.offsetRight)
	}

	r.Rectangle.Draw(s)
}

// Tick comment
func (r *CollBat) Tick(ev tl.Event) {
	// Enable arrow key movement
	if ev.Type == tl.EventKey && r.move {
		switch ev.Key {
		case tl.KeyArrowUp:
			r.py--
		case tl.KeyArrowDown:
			r.py++
		}

		// send the coordinates to the server
		c.Send(r.px, r.py, 0, 0)
	}

	r.SetPosition(r.px, r.py)
}

// Collide comment
func (r *CollBat) Collide(p tl.Physical) {
	// Check if it's a CollRec we're colliding with
	/*
		if _, ok := p.(*CollBat); ok && r.move {
			r.SetColor(tl.ColorBlue)
			r.SetPosition(r.px, r.py)
		}
	*/
}

func main() {
	env.Parse()

	f, err := os.Create("out.log")
	if err != nil {
		log.Fatal(err)
	}
	opt := &hclog.LoggerOptions{Output: f}
	logger = hclog.New(opt)

	c = client.New("localhost:6000")
	go func() {
		for d := range c.Recieve() {
			bat2.py = d.BatY
			logger.Info("something", "y", d.BatY)
		}
	}()

	oneEnabled := true
	twoEnabled := false

	if *player != 1 {
		oneEnabled = false
		twoEnabled = true
	}

	bat1 = NewCollBat(3, 3, 3, 6, tl.ColorRed, oneEnabled, 0)
	bat2 = NewCollBat(3, 4, 3, 6, tl.ColorGreen, twoEnabled, -3)

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
