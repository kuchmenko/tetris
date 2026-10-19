package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func runApp() {
	rl.InitWindow(1000, 1000, "raylib [core] example - basic window")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	camera := rl.NewCamera2D(rl.NewVector2(0, 0), rl.NewVector2(0, 0), 0, 1)

	var stepTickAcc float32 = 0
	var stepTick float32 = 500
	var inputTick float32 = 100
	var inputTickAcc float32 = 0

	tetris := NewTetris(20, 20)

	cellSize := int32(rl.GetScreenWidth() / tetris.cols)
	renderer := NewRenderer(rl.GetScreenWidth(), rl.GetScreenHeight(), &camera)

	for !rl.WindowShouldClose() {
		stepTickAcc += rl.GetFrameTime() * 1000
		inputTickAcc += rl.GetFrameTime() * 100

		handleInput(tetris)
		for inputTickAcc >= inputTick {
			inputTickAcc -= inputTick
		}

		for stepTickAcc >= stepTick {
			stepTickAcc -= stepTick

			tetris.Tick()

			if tetris.currentFigure != nil {
				if _, ok := renderer.colors[tetris.currentFigure.ID]; !ok {
					renderer.colors[tetris.currentFigure.ID] = TetrisColors[rand.Intn(len(TetrisColors))]
				}
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		for i := 0; i < tetris.rows; i++ {
			for j := 0; j < tetris.cols; j++ {
				posX := int32(j) * cellSize
				posY := int32(i) * cellSize
				rl.DrawRectangleLines(posX, posY, cellSize, cellSize, rl.LightGray)
			}
		}
		for _, figure := range tetris.figures {
			position := rl.NewVector2(float32(figure.Position.X)*float32(cellSize), float32(figure.Position.Y)*float32(cellSize))

			for _, cell := range figure.Cells {
				cellPosition := rl.NewVector2(float32(cell.Position.X)*float32(cellSize), float32(cell.Position.Y)*float32(cellSize))
				rl.DrawRectangle(int32(position.X+cellPosition.X), int32(position.Y+cellPosition.Y), cellSize, cellSize, renderer.colors[figure.ID])
			}
		}

		rl.EndDrawing()

		// transform tetris coordinates to screen coordinates
		// draw tetris
		// free spots
		// figures with stable colors (map?)

	}
}
