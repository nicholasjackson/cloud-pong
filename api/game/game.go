package game

import (
	"fmt"
	"time"

	pb "github.com/nicholasjackson/pong/api/protos/pong"
)

const (
	gameWidth       float64 = 1024
	gameHeight      float64 = 768
	tickInterval            = 50 * time.Millisecond
	initialSpeed            = 5
	speedMultiplier         = 1.2
)

type Object struct {
	px float64
	py float64
	w  float64
	h  float64
}

type Game struct {
	ball              *Object
	bat1              *Object
	bat2              *Object
	started           bool
	speedX            float64
	speedY            float64
	controllingPlayer int
}

func NewGame() *Game {
	return &Game{
		ball: &Object{0, 0, 15, 20},
		bat1: &Object{0, 0, 15, 90},
		bat2: &Object{0, 0, 15, 90},
	}
}

func (r *Game) DataAsProto() *pb.PongData {
	return &pb.PongData{
		Bat1:  &pb.Bat{X: int32(r.bat1.px), Y: int32(r.bat1.py), W: int32(r.bat1.w), H: int32(r.bat1.h)},
		Bat2:  &pb.Bat{X: int32(r.bat2.px), Y: int32(r.bat2.py), W: int32(r.bat2.w), H: int32(r.bat2.h)},
		Ball:  &pb.Ball{X: int32(r.ball.px), Y: int32(r.ball.py), W: int32(r.ball.w), H: int32(r.ball.h)},
		Game:  &pb.Game{W: int32(gameWidth), H: int32(gameHeight)},
		Score: 0,
	}
}

func (r *Game) StartGame() {
	r.started = true
}

func (r *Game) ResetGame() {
	r.ball.px = r.bat1.w * 2
	r.ball.py = (gameHeight / 2) - (r.ball.h / 2)

	r.bat1.px = r.bat1.w
	r.bat1.py = (gameHeight / 2) - (r.bat1.h / 2)

	r.bat2.px = gameWidth - r.bat2.w*2
	r.bat2.py = (gameHeight / 2) - (r.bat2.h / 2)

	// reset parameters
	r.started = false
	r.speedX = initialSpeed
	r.speedY = 0
	r.controllingPlayer = 1
}

func (r *Game) Tick() <-chan struct{} {
	tick := make(chan struct{})
	ticker := time.NewTicker(tickInterval)

	go func() {
		for {
			select {
			case <-ticker.C:
				r.tick()
				tick <- struct{}{}
			}
		}
	}()

	return tick
}

func (r *Game) tick() {

	if r.started == true {
		r.moveBall()
	}

	fmt.Println(r.ball.px)

	// check to see if the ball would hit a bat and flip x speed
	if r.ball.px+r.ball.w >= r.bat2.px && r.controllingPlayer == 1 {
		r.speedX = r.speedX * speedMultiplier
		r.speedX = -r.speedX
		r.controllingPlayer = 2
	}

	if r.ball.px <= r.bat1.px && r.controllingPlayer == 2 {
		r.speedX = r.speedX * speedMultiplier
		r.speedX = -r.speedX
		r.controllingPlayer = 1
	}
}

func (r *Game) moveBall() {
	r.ball.px += r.speedX
	r.ball.py += r.speedY
}
