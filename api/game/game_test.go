package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup() *Game {
	return &Game{
		ball: &Object{0, 0, 3, 3},
		bat1: &Object{0, 0, 3, 6},
		bat2: &Object{0, 0, 3, 6},
	}
}

func TestGameDrawsBatsToCorrectLocationOnStart(t *testing.T) {
	g := setup()

	g.ResetGame()

	// bat position should be width of the bat from the 0
	assert.Equal(t, 3.0, g.bat1.px)

	// bat position should be 2x width of the bat from the screenWidth
	assert.Equal(t, 1018.0, g.bat2.px)

	// bat position should be centered on the screen
	assert.Equal(t, 381.0, g.bat1.py)
	assert.Equal(t, 381.0, g.bat2.py)
}

func TestGameDrawsBallToCorrectLocationOnStart(t *testing.T) {
	g := setup()

	g.ResetGame()

	// ball position should be 2x width of the bat from the screenWidth
	assert.Equal(t, 6.0, g.ball.px)

	// ball position should be centered on the screen
	assert.Equal(t, 382.5, g.ball.py)
}

func TestGameMovesBallToCorrectLocationOnTick(t *testing.T) {
	g := setup()

	g.ResetGame()
	g.StartGame()

	g.tick()
	g.tick()

	// ball position should be 2x width of the bat from the screenWidth
	assert.Equal(t, 16.0, g.ball.px)
}

func TestGameBallDoesNotMoveOnTickWhenNotStarted(t *testing.T) {
	g := setup()

	g.ResetGame()

	g.tick()
	g.tick()

	// ball position should be 2x width of the bat from the screenWidth
	assert.Equal(t, 6.0, g.ball.px)
}

func TestGameBallChangesDirectionAndSpeedWhenHittingPlayer2Bat(t *testing.T) {
	g := setup()

	g.ResetGame()
	g.StartGame()

	g.ball.px = 1015 // move the ball to be incontact with the bat
	g.tick()
	g.tick()
	g.tick()

	// ball position should be moving in the opposite direction
	assert.Equal(t, 1008.0, g.ball.px)
}

func TestGameBallChangesDirectionWhenHittingPlayer1Bat(t *testing.T) {
	g := setup()

	g.ResetGame()
	g.StartGame()

	g.ball.px = 6 // move the ball to be incontact with the bat
	g.controllingPlayer = 2
	g.speedX = -initialSpeed
	g.tick()
	g.tick()
	g.tick()

	// ball position should be moving in the opposite direction
	assert.Equal(t, 13.0, g.ball.px)
}
