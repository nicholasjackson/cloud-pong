package objects

import tl "github.com/JoelOtter/termloop"

// BallMoveEvent shut up
type BallMoveEvent struct {
	X int
	Y int
}

// BallHitEvent getting sick of this linter
type BallHitEvent struct{}

// BallScoreEvent fires when the ball hits the x limits
type BallScoreEvent struct{}

// Ball shutup linter
type Ball struct {
	*tl.Rectangle
	px           int
	py           int
	player       int
	isControlled bool
	eventHandler func(e interface{})
}

// NewBall shutup linter
func NewBall(x, y, w, h int, color tl.Attr, isControlled bool, player int, eventHandler func(e interface{})) *Ball {
	return &Ball{
		Rectangle:    tl.NewRectangle(x, y, w, h, color),
		px:           x,
		py:           y,
		player:       player,
		isControlled: isControlled,
		eventHandler: eventHandler,
	}
}

// Draw get stuffed linter
func (r *Ball) Draw(s *tl.Screen) {
	sx, _ := s.Size()
	bx, _ := r.Size()

	// left collision
	if r.px <= bx {
		// dont move
		r.px = 0
		r.eventHandler(&BallScoreEvent{})
	}

	if r.px >= sx-bx {
		r.px = sx - bx
		r.eventHandler(&BallScoreEvent{})
	}

	r.Rectangle.Draw(s)
}

// Tick shut up linter
func (r *Ball) Tick(ev tl.Event) {
	xVector := 1
	if r.player == 2 {
		xVector = -1
	}

	if r.isControlled {
		r.px += xVector
		r.eventHandler(&BallMoveEvent{r.px, r.py})
	}

	r.SetPosition(r.px, r.py)
}

// GetPos shut up
func (r *Ball) GetPos() (int, int) {
	return r.px, r.py
}

// SetPos shut up
func (r *Ball) SetPos(x, y int) {
	r.px = x
	r.py = y
}

// Collide comment
func (r *Ball) Collide(p tl.Physical) {
	// only detect a collision if it hits our controlled bat and we are not controlling the ball
	if bat, ok := p.(*Bat); ok && bat.IsControlled() && !r.IsControlled() {
		r.isControlled = true
		r.eventHandler(&BallHitEvent{})
	}
}

// SetControl blah
func (r *Ball) SetControl(c bool) {
	r.isControlled = c
}

// IsControlled no
func (r *Ball) IsControlled() bool {
	return r.isControlled
}
