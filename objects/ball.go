package objects

import tl "github.com/JoelOtter/termloop"

// BallMoveEvent shut up
type BallMoveEvent struct {
	X int
	Y int
}

// BallHitEvent getting sick of this linter
type BallHitEvent struct{}

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

// Tick shut up linter
func (r *Ball) Tick(ev tl.Event) {
	// check to see if we have exceeded the bounds

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
