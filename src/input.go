package main

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var inputTick float32 = 100
var inputTickAcc float32 = 0

func handleInput(tetris *Tetris) {
	inputTickAcc += rl.GetFrameTime() * 1000

	if tetris.isFinished && rl.IsKeyPressed(rl.KeyEnter) {
		tetris.Restart()
		return
	}

	if tetris.isFinished {
		return
	}

	var move MoveDirection
	hasMoved := false
	hasRotated := false

	if rl.IsKeyDown(rl.KeyR) {
		if rl.IsKeyPressed(rl.KeyR) {
			hasRotated = true
			inputTickAcc = inputTick
		}
	}

	if rl.IsKeyDown(rl.KeyDown) {
		hasMoved = true
		move = Down

		if rl.IsKeyPressed(rl.KeyDown) {
			inputTickAcc = inputTick
		}
	}

	if rl.IsKeyDown(rl.KeyLeft) {
		hasMoved = true
		move = Left

		if rl.IsKeyPressed(rl.KeyLeft) {
			inputTickAcc = inputTick
		}
	}

	if rl.IsKeyDown(rl.KeyRight) {
		hasMoved = true
		move = Right

		if rl.IsKeyPressed(rl.KeyRight) {
			inputTickAcc = inputTick
		}
	}

	if inputTickAcc >= inputTick {
		if hasRotated {
			err := tetris.Rotate()
			if err != nil {
				log.Printf("handle rotate move: %s", err)
			}
		}

		if hasMoved {
			err := tetris.MoveActive(move)
			if err != nil {
				log.Printf("handle input move: %s", err)
			}
		}

		inputTickAcc = -inputTick
	}
}
