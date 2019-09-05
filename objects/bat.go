package objects

import (
	//	"fmt"
	"fmt"

	tl "github.com/JoelOtter/termloop"
)

// Ball shutup linter
type Bat struct {
	*tl.Rectangle
	player         int
	singleKeyboard bool
	screenW        int
	screenH        int
	handler        func(string)
}

// NewBall shutup linter
func NewBat(x, y, w, h int, color tl.Attr, player int, singleKeyboard bool, handler func(string)) *Bat {
	return &Bat{
		Rectangle:      tl.NewRectangle(x, y, w, h, color),
		player:         player,
		singleKeyboard: singleKeyboard,
		handler:        handler,
	}
}

func (b *Bat) Tick(ev tl.Event) {
	if ev.Type == tl.EventKey && b.handler != nil {

		if b.singleKeyboard {
			// Control both bats from a single keyboard
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
		} else {
			// Normal keys
			switch ev.Key {
			case tl.KeyArrowUp:
				b.handler(fmt.Sprintf("BAT_UP_%d", b.player))
				return
			case tl.KeyArrowDown:
				b.handler(fmt.Sprintf("BAT_DOWN_%d", b.player))
				return
			}
		}

		switch ev.Key {
		case tl.KeyCtrlR:
			b.handler("RESET_GAME")
		case tl.KeySpace:
			b.handler(fmt.Sprintf("SERVE_%d", b.player))
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
