package objects

import (
	"time"

	tl "github.com/JoelOtter/termloop"
)

// BatMoveEvent is fired when the bat moves
type BatMoveEvent struct {
	X int
	Y int
}

// Bat defines a bat in the game
type Bat struct {
	*tl.Rectangle
	px           int
	py           int
	offsetRight  int
	isControlled bool
	speed        int
	lastPress    time.Time
	eventHandler func(e interface{})
}

// NewBat comment
func NewBat(x, y, w, h int, color tl.Attr, offsetRight int, isControlled bool, eventHandler func(e interface{})) *Bat {
	return &Bat{
		Rectangle:    tl.NewRectangle(x, y, w, h, color),
		px:           x,
		py:           y,
		offsetRight:  offsetRight,
		isControlled: isControlled,
		speed:        1,
		eventHandler: eventHandler,
	}
}

// SetPos sets the position of the bat
func (r *Bat) SetPos(x, y int) {
	r.px = x
	r.py = y
}

// GetPos returns the x and y position for the bat
func (r *Bat) GetPos() (int, int) {
	return r.px, r.py
}

// Draw something
func (r *Bat) Draw(s *tl.Screen) {
	sx, sy := s.Size()
	bx, by := r.Size()

	// is this the first draw if so set to center
	if r.py == 0 {
		r.py = (sy / 2) - (by / 2)
		return
	}

	// if the bat y is less than the bounds set to the bounds
	if minY := (0 - by/2); r.py < minY {
		r.py = minY
	}

	// if the bat is greater than the bounds set to the bounds
	if maxY := (sy - by/2); r.py > maxY {
		r.py = maxY
	}

	if r.offsetRight < 0 {
		r.px = sx - (bx - r.offsetRight)
	}

	r.Rectangle.Draw(s)
}

// Tick comment
func (r *Bat) Tick(ev tl.Event) {

	// Enable arrow key movement
	if ev.Type == tl.EventKey && r.isControlled {
		switch ev.Key {
		case tl.KeyArrowUp:
			r.py -= r.speed

			// increase the bat speed as the arrow is held down
			r.speed++
			r.lastPress = time.Now()
		case tl.KeyArrowDown:
			r.py += r.speed

			// increase the bat speed as the arrow is held down
			r.speed++
			r.lastPress = time.Now()
		}

		if r.isControlled {
			r.eventHandler(&BatMoveEvent{r.px, r.py})
		}
	} else {
		// reset the bat speed after a timeout
		if time.Now().Sub(r.lastPress) > 200*time.Millisecond {
			r.speed = 1
		}
	}

	r.SetPosition(r.px, r.py)
}

// Collide comment
func (r *Bat) Collide(p tl.Physical) {
	// Check if it's a CollRec we're colliding with
	/*
		if _, ok := p.(*CollBat); ok && r.move {
			r.SetColor(tl.ColorBlue)
			r.SetPosition(r.px, r.py)
		}
	*/
}

// IsControlled is this bat user controlled?
func (r *Bat) IsControlled() bool {
	return r.isControlled
}

func (r *Bat) handleEvent(e interface{}) {
	r.eventHandler(e)
}

// Reset the original settings of the bat
func (r *Bat) Reset() {
	r.py = 0
	r.SetPosition(r.px, r.py)
}
