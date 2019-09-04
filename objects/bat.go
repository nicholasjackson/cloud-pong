package objects

import (
	//	"fmt"
	tl "github.com/JoelOtter/termloop"
)

// Ball shutup linter
type Bat struct {
	*tl.Rectangle
	screenW int
	screenH int
	handler func(string)
}

// NewBall shutup linter
func NewBat(x, y, w, h int, color tl.Attr, handler func(string)) *Bat {
	return &Bat{
		Rectangle: tl.NewRectangle(x, y, w, h, color),
		handler:   handler,
	}
}

func (b *Bat) Tick(ev tl.Event) {
	if ev.Type == tl.EventKey && b.handler != nil {
		switch ev.Ch {
		case 119: // W
			b.handler("BAT_UP_1")
			return
		case 115: // S
			b.handler("BAT_DOWN_1")
			return
		case 101: // E
			b.handler("SERVE_1")
			return
		case 111: // O
			b.handler("BAT_UP_2")
			return
		case 108: // L
			b.handler("BAT_DOWN_2")
			return
		case 112: // P
			b.handler("SERVE_2")
			return
		}

		switch ev.Key {
		case tl.KeyCtrlR:
			b.handler("RESET_GAME")
		case tl.KeySpace:
			b.handler("SERVE")
		}
	}
}

func (b *Bat) Draw(s *tl.Screen) {
	b.screenW, b.screenH = s.Size()
	b.Rectangle.Draw(s)
}

func (b *Bat) SetData(bx, by, bw, bh, gw, gh int32) {
	// before setting translate coordinates to screen size
	//panic(b.translateToGameX(bw, gw))
	b.SetSize(
		translateToGameX(bw, gw, b.screenW),
		translateToGameY(bh, gh, b.screenH),
	)

	b.SetPosition(
		translateToGameX(bx, gw, b.screenW),
		translateToGameY(by, gh, b.screenH),
	)
}
