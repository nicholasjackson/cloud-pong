package objects

import (
	//	"fmt"
	tl "github.com/JoelOtter/termloop"
)

// Ball shutup linter
type Ball struct {
	*tl.Rectangle
	screenW int
	screenH int
}

// NewBall shutup linter
func NewBall(x, y, w, h int, color tl.Attr) *Ball {
	return &Ball{
		Rectangle: tl.NewRectangle(x, y, w, h, color),
	}
}

func (b *Ball) Draw(s *tl.Screen) {
	b.screenW, b.screenH = s.Size()
	b.Rectangle.Draw(s)
}

func (b *Ball) SetData(bx, by, bw, bh, gw, gh int32) {
	// before setting translate coordinates to screen size
	//panic(translateToGameX(bw, gw, b.screenW))
	b.SetSize(
		translateToGameX(bw, gw, b.screenW),
		translateToGameY(bh, gh, b.screenH),
	)

	b.SetPosition(
		translateToGameX(bx, gw, b.screenW),
		translateToGameY(by, gh, b.screenH),
	)
}
