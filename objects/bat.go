package objects

import (
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
	if r.offsetRight < 0 {
		sx, _ := s.Size()
		bx, _ := r.Size()

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
			r.py--
		case tl.KeyArrowDown:
			r.py++
		}

		if r.isControlled {
			r.eventHandler(&BatMoveEvent{r.px, r.py})
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

func (r *Bat) handleEvent(e interface{}) {
	r.eventHandler(e)
}
