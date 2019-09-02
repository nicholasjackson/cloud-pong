package objects

import (
	tl "github.com/JoelOtter/termloop"
)

var scores = [10][6][6]int{
	{
		{1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1},
	},
	{
		{1, 1, 1, 1, 0, 0},
		{0, 0, 0, 1, 0, 0},
		{0, 0, 0, 1, 0, 0},
		{0, 0, 0, 1, 0, 0},
		{0, 0, 0, 1, 0, 0},
		{1, 1, 1, 1, 1, 1},
	},
	{
		{1, 1, 1, 1, 1, 1},
		{0, 0, 0, 0, 0, 1},
		{0, 0, 0, 1, 1, 0},
		{0, 1, 1, 0, 0, 0},
		{1, 0, 0, 0, 0, 0},
		{1, 1, 1, 1, 1, 1},
	},
	{
		{1, 1, 1, 1, 1, 1},
		{0, 0, 0, 0, 1, 1},
		{0, 0, 1, 1, 0, 0},
		{0, 0, 1, 1, 0, 0},
		{0, 0, 0, 0, 1, 1},
		{1, 1, 1, 1, 1, 1},
	},
	{
		{1, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 0, 0},
		{1, 0, 1, 1, 0, 0},
		{1, 1, 1, 1, 1, 1},
		{0, 0, 1, 1, 0, 0},
	},
	{
		{1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0},
		{1, 1, 1, 1, 1, 1},
		{0, 0, 0, 0, 0, 1},
		{0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1},
	},
}

// Score defines a Score for the game
type Score struct {
	*tl.Rectangle
	blocks []*tl.Rectangle
	color  tl.Attr
	value  int
	x      int
	y      int
}

// NewScore creates a new score
func NewScore(x, y int, color tl.Attr) *Score {
	return &Score{
		color:     color,
		Rectangle: tl.NewRectangle(0, 0, 0, 0, color),
		x:         x,
		y:         y,
	}
}

// Draw something
func (s *Score) Draw(sc *tl.Screen) {
	sx, _ := sc.Size()

	// remove the old blocks
	for _, b := range s.blocks {
		sc.RemoveEntity(b)
	}

	width := 1
	height := 1

	x := sx/2 + s.x
	y := s.y

	blocks := make([]*tl.Rectangle, 0)

	for r, row := range scores[s.value] {
		for c, col := range row {
			if col == 1 {
				rect := tl.NewRectangle(x+c+width, y+r+height, width, height, s.color)
				blocks = append(blocks, rect)
				sc.AddEntity(rect)
			}
		}
	}

	s.blocks = blocks
}

// UpdateScore updates the score
func (s *Score) UpdateScore(v int32) {
	s.value = int(v)
}

// IncrementScore updates the score by 1
func (s *Score) IncrementScore() {
	s.value++

	if s.value > len(scores) {
		panic("Score out of range")
	}
}
