package objects

import tl "github.com/JoelOtter/termloop"

// Net defines a Net for the game
type Net struct {
	*tl.Rectangle
	blocks []*tl.Rectangle
	color  tl.Attr
	prevH  int
	prevW  int
}

// NewNet creates a new net and positions at the middle of the screen
func NewNet(color tl.Attr) *Net {
	return &Net{
		color:     color,
		Rectangle: tl.NewRectangle(0, 0, 0, 0, color),
	}
}

// Draw something
func (r *Net) Draw(s *tl.Screen) {
	sx, sy := s.Size()

	if sx == r.prevW && sy == r.prevH {
		return
	}

	r.prevW, r.prevH = sx, sy

	// remove the old blocks
	for _, b := range r.blocks {
		s.RemoveEntity(b)
	}

	// calculate the required number of blocks
	blocks := make([]*tl.Rectangle, sy/3)

	for n := 0; n < len(blocks); n++ {
		hpos := n * 4
		b := tl.NewRectangle(sx/2-3, hpos, 3, 2, r.color)
		s.AddEntity(b)
		r.blocks = append(r.blocks, b)
	}
	//r.Rectangle.Draw(s)
}
