package main

import (
	"log"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Bounded struct {
	MinX, MaxX, MinY, MaxY int
}

func BoundedFromFigure(figure TetrisFigure) *Bounded {
	b := &Bounded{}
	for _, row := range figure.Cells {
		if row.Position.X < b.MinX {
			b.MinX = row.Position.X
		}
		if row.Position.X > b.MaxX {
			b.MaxX = row.Position.X
		}
		if row.Position.Y < b.MinY {
			b.MinY = row.Position.Y
		}
		if row.Position.Y > b.MaxY {
			b.MaxY = row.Position.Y
		}
	}
	return b
}

func (b *Bounded) Height() int {
	return b.MaxY - b.MinY + 1
}

func (b *Bounded) Width() int {
	return b.MaxX - b.MinX + 1
}

func runApp() {
	rl.InitWindow(1000, 1000, "raylib [core] example - basic window")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	camera := rl.NewCamera2D(rl.NewVector2(0, 0), rl.NewVector2(0, 0), 0, 1)

	var stepTickAcc float32 = 0
	var stepTick float32 = 500
	rows := 10
	cols := 10

	tetris := NewTetris(cols, rows)

	cellSize := int32(rl.GetScreenWidth() / tetris.cols)
	renderer := NewRenderer(rl.GetScreenWidth(), rl.GetScreenHeight(), &camera)
	figures := []TetrisFigure{
		NewTetrisFigureO(1),
		NewTetrisFigureO(2),
		NewTetrisFigureL(3),
		NewTetrisFigureJ(4),
		NewTetrisFigureO(5),
		NewTetrisFigureO(6),
	}

	lastX := 0
	posY := tetris.rows

	boundedFigures := make([]*Bounded, 0, len(figures))

	for id, figure := range figures {
		if id > 4 && id%5 == 0 {
			log.Printf("id: %d, posY: %d\n", id, posY)
			posY -= boundedFigures[0].Height()
			lastX = 0
		}
		bounded := BoundedFromFigure(figure)
		figure.Position.X = lastX
		figure.Position.Y = posY - bounded.Height()
		boundedFigures = append(boundedFigures, bounded)
		tetris.AddFigure(figure)
		renderer.colors[figure.ID] = TetrisColors[rand.Intn(len(TetrisColors))]
		lastX += bounded.Width()
	}

	for !rl.WindowShouldClose() {
		stepTickAcc += rl.GetFrameTime() * 1000
		inputTickAcc += rl.GetFrameTime() * 100

		handleInput(tetris)

		for stepTickAcc >= stepTick {
			stepTickAcc -= stepTick

			tetris.Tick()

			if tetris.active != nil {
				if _, ok := renderer.colors[tetris.active.ID]; !ok {
					renderer.colors[tetris.active.ID] = TetrisColors[rand.Intn(len(TetrisColors))]
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
	}
}
