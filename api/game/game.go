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
	maxSpeed                = 20
	speedMultiplier         = 1.2
	batSpeed                = 10
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
	bat1LastMove      time.Time
	bat1Speed         float64
	bat2Speed         float64
	bat2LastMove      time.Time
	started           bool
	speedX            float64
	speedY            float64
	controllingPlayer int
	player1Score      int
	player2Score      int
	cancelTick        chan struct{}
}

func NewGame() *Game {
	return &Game{
		ball:              &Object{0, 0, 15, 20},
		bat1:              &Object{0, 0, 15, 90},
		bat2:              &Object{0, 0, 15, 90},
		controllingPlayer: 1,
	}
}

func (r *Game) DataAsProto() *pb.PongData {
	fmt.Println("Score1", r.player1Score, "Score2", r.player2Score)
	return &pb.PongData{
		Bat1:         &pb.Bat{X: int32(r.bat1.px), Y: int32(r.bat1.py), W: int32(r.bat1.w), H: int32(r.bat1.h)},
		Bat2:         &pb.Bat{X: int32(r.bat2.px), Y: int32(r.bat2.py), W: int32(r.bat2.w), H: int32(r.bat2.h)},
		Ball:         &pb.Ball{X: int32(r.ball.px), Y: int32(r.ball.py), W: int32(r.ball.w), H: int32(r.ball.h)},
		Game:         &pb.Game{W: int32(gameWidth), H: int32(gameHeight)},
		Player1Score: int32(r.player1Score),
		Player2Score: int32(r.player2Score),
	}
}

func (r *Game) StartGame(player int) {
	if player == r.controllingPlayer {
		r.started = true
	}
}

// HardReset the game to the intial state
func (r *Game) HardReset() {
	r.controllingPlayer = 1
	r.player1Score = 0
	r.player2Score = 0
	r.ResetGame()
}

func (r *Game) ResetGame() {
	fmt.Println("Reset Game")

	// ball position is based on who has the serve
	r.ball.px = r.bat1.w * 2

	if r.controllingPlayer == 2 {
		r.ball.px = gameWidth - (r.bat1.w * 2)
	}

	r.ball.py = (gameHeight / 2) - (r.ball.h / 2)

	r.bat1.px = r.bat1.w
	r.bat1.py = (gameHeight / 2) - (r.bat1.h / 2)

	r.bat2.px = gameWidth - r.bat2.w*2
	r.bat2.py = (gameHeight / 2) - (r.bat2.h / 2)

	// reset parameters
	r.started = false

	r.speedY = 0
	r.speedX = -initialSpeed

	if r.controllingPlayer == 2 {
		r.speedX = initialSpeed
	}

	r.bat1Speed = batSpeed
	r.bat2Speed = batSpeed
}

func (r *Game) MoveBatUp(player int) {
	r.CalculateBatSpeed(player)

	minY := 0 - (r.bat1.h / 2)
	if player == 1 {
		r.bat1.py -= r.bat1Speed

		if r.bat1.py < minY {
			r.bat1.py = minY
		}
		return
	}

	r.bat2.py -= r.bat2Speed
	if r.bat2.py < minY {
		r.bat2.py = minY
	}
}

func (r *Game) MoveBatDown(player int) {
	r.CalculateBatSpeed(player)

	maxY := gameHeight - r.bat1.h/2

	if player == 1 {
		r.bat1.py += r.bat1Speed

		if r.bat1.py > maxY {
			r.bat1.py = maxY
		}

		return
	}

	r.bat2.py += r.bat2Speed
	if r.bat2.py > maxY {
		r.bat2.py = maxY
	}
}

func (r *Game) CalculateBatSpeed(player int) {
	// if the user is holding the up button set the multiplier
	if player == 1 {
		if !r.bat1LastMove.IsZero() && time.Now().Sub(r.bat1LastMove) < (200*time.Millisecond) {
			r.bat1Speed += batSpeed
		} else {
			r.bat1Speed = batSpeed
		}

		fmt.Println("Bat speed", r.bat1Speed)

		r.bat1LastMove = time.Now()
	}

	if player == 2 {
		if !r.bat2LastMove.IsZero() && time.Now().Sub(r.bat2LastMove) < (200*time.Millisecond) {
			r.bat2Speed += batSpeed
		} else {
			r.bat2Speed = batSpeed
		}

		r.bat2LastMove = time.Now()
	}
}

func (r *Game) Tick() (tick chan struct{}, cancel chan struct{}) {
	tick = make(chan struct{})
	cancel = make(chan struct{})

	ticker := time.NewTicker(tickInterval)

	go func() {
		for r.started {
			<-ticker.C
			r.tick()
			tick <- struct{}{}
		}
		fmt.Println("Cancel game")
		ticker.Stop()
		cancel <- struct{}{}
	}()

	return tick, cancel
}

func (r *Game) tick() {
	if !r.started == true {
		return
	}

	// check to see if the ball would hit a bat and flip x speed
	if b := r.hitBat(); b != nil {
		// increase the ball speed if it less than the max
		if r.speedX < maxSpeed && r.speedX > -maxSpeed {
			r.speedX = r.speedX * speedMultiplier
		}

		// calculate the y speed
		// angle is based on the distance from the center
		cbat := b.h/2 + b.py            // center of bat
		cball := r.ball.h/2 + r.ball.py //center of ball

		r.speedY = (float64(cball-cbat) * initialSpeed) / 6

		r.speedX = -r.speedX
		r.controllingPlayer = -r.controllingPlayer // flip the controlling player
	}

	if r.score1() {
		r.player2Score++
		r.ResetGame()
	}

	if r.score2() {
		r.player1Score++
		r.ResetGame()
	}

	// if the ball hits the bounds flip direction
	if r.ball.py <= 0 || r.ball.py >= gameHeight {
		r.speedY = -r.speedY
	}

	r.moveBall()
}

// is the ball hitting the player bat?
func (r *Game) hitBat() *Object {
	if r.ball.px < r.bat1.px &&
		(r.ball.py+r.ball.h >= r.bat1.py && r.ball.py <= r.bat1.py+r.bat1.h) {
		return r.bat1
	}

	if r.ball.px+r.ball.w > r.bat2.px &&
		(r.ball.py+r.ball.h >= r.bat2.py && r.ball.py <= r.bat2.py+r.bat2.h) {
		return r.bat2
	}

	return nil
}

// did the ball go out of bounds in the player 1 area?
func (r *Game) score1() bool {
	if r.ball.px <= 0 {
		// give player 1 the serve
		r.controllingPlayer = 1
		return true
	}

	return false
}

// did the ball go out of bounds in the player 2 areas?
func (r *Game) score2() bool {
	if r.ball.px+r.ball.w >= gameWidth {
		// give player 2 the serve
		r.controllingPlayer = 2
		return true
	}

	return false
}

func (r *Game) moveBall() {
	r.ball.py += r.speedY

	// if we are hitting the player 1 bat set the px to the bat position
	b := r.hitBat()

	if b == r.bat1 {
		r.ball.px = r.bat1.px + r.bat1.w
		return
	}

	// if we are hitting the player 2 bat set the px to the bat position
	if b == r.bat2 {
		r.ball.px = r.bat2.px - r.ball.w
		return
	}

	r.ball.px += r.speedX
}
