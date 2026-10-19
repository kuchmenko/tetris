package main

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func handleInput(tetris *Tetris) {
	if tetris.isFinished {
		return
	}

	if rl.IsKeyPressed(rl.KeyR) {
		err := tetris.Rotate()
		if err != nil {
			log.Printf("handle rotate move: %s", err)
		}
	}

	var move MoveDirection
	var hasMoved bool

	if rl.IsKeyDown(rl.KeyDown) {
		move = Down
		hasMoved = true
	}

	if rl.IsKeyPressed(rl.KeyLeft) {
		move = Left
		hasMoved = true
	}

	if rl.IsKeyPressed(rl.KeyRight) {
		move = Right
		hasMoved = true
	}

	if !hasMoved {
		return
	}

	err := tetris.Move(move)
	if err != nil {
		log.Printf("handle input move: %s", err)
	}
}
