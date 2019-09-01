package game

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setup() *Game {
	return &Game{
		ball:              &Object{0, 0, 3, 3},
		bat1:              &Object{0, 0, 3, 6},
		bat2:              &Object{0, 0, 3, 6},
		controllingPlayer: 1,
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
	assert.True(t, g.ball.px < 1015)
}

func TestGameScoreWhenMissingPlayer2Bat(t *testing.T) {
	g := setup()

	g.ResetGame()
	g.StartGame()

	g.ball.px = 1017 // move the ball to be incontact with the bat
	g.bat2.py = 34
	g.controllingPlayer = 1
	g.tick()

	// ball position should be moving in the opposite direction
	assert.Equal(t, 1, g.player1Score)
	// ball should be reset
	assert.Equal(t, 6.0, g.ball.px)
}

func TestGameScoreWhenMissingPlayer1Bat(t *testing.T) {
	g := setup()

	g.ResetGame()
	g.StartGame()

	g.ball.px = 3 // move the ball to be incontact with the bat
	g.bat1.py = 34
	g.controllingPlayer = 2
	g.speedX = -initialSpeed
	g.tick()

	// ball position should be moving in the opposite direction
	assert.Equal(t, 1, g.player2Score)
	// ball should be reset
	assert.Equal(t, 6.0, g.ball.px)
}

func TestGameBat1MovesUpToCorrectPossition(t *testing.T) {
	g := setup()

	g.ResetGame()
	g.MoveBatUp(1)

	assert.Equal(t, 371.0, g.bat1.py)
}

func TestGameBat2MovesUpToCorrectPossition(t *testing.T) {
	g := setup()

	g.ResetGame()
	g.MoveBatUp(2)

	assert.Equal(t, 371.0, g.bat2.py)
}

func TestGameBat1MovesDownToCorrectPossition(t *testing.T) {
	g := setup()

	g.ResetGame()
	g.MoveBatDown(1)

	assert.Equal(t, 391.0, g.bat1.py)
}

func TestGameBat2MovesDownToCorrectPossition(t *testing.T) {
	g := setup()

	g.ResetGame()
	g.MoveBatDown(2)

	assert.Equal(t, 391.0, g.bat2.py)
}

func TestGameTickerCanBeCancelledWhenReset(t *testing.T) {
	g := setup()
	done := make(chan struct{})

	g.ResetGame()
	go func() {
		for _ = range g.Tick() {
			fmt.Println("tick")
		}
		done <- struct{}{} // signal the ticker has been cancelled
	}()

	time.Sleep(100 * time.Millisecond)
	g.ResetGame() // resetting game cancels ticker

	select {
	case <-done:
		return
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for ticker to cancel")
	}
}
