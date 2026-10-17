package main

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Grid struct {
	Size int
}

func (g *Grid) Draw() {
	for i := 0; i < g.Size; i++ {
		for j := 0; j < g.Size; j++ {
			rl.DrawRectangle(int32(i*g.Size), int32(j*g.Size), int32(g.Size), int32(g.Size), rl.LightGray)
			rl.DrawRectangleLines(int32(i*g.Size), int32(j*g.Size), int32(g.Size), int32(g.Size), rl.Black)
		}
	}
}

func (g *Grid) IsInGrid(position rl.Vector2) bool {
	x := int(position.X)
	y := int(position.Y)
	max := g.Size * g.Size

	if x < 0 || x >= max {
		return false
	}

	if y < 0 || y >= max {
		return false
	}

	return true
}

func (g *Grid) ToGridPosition(x, y int) rl.Vector2 {
	return rl.NewVector2(float32(x*g.Size), float32(y*g.Size))
}

func (g *Grid) ViewportBounds(x, y float32) rl.Vector2 {
	return rl.NewVector2(x*float32(g.Size), y*float32(g.Size))
}

type Cell struct {
	Position rl.Vector2
	Size     int
	Color    rl.Color
}

func (c *Cell) Move(position rl.Vector2) {
	c.Position.X += position.X
	c.Position.Y += position.Y
}

func NewCell(position rl.Vector2, size int, color rl.Color) Cell {
	position.X = position.X * float32(size)
	position.Y = position.Y * float32(size)
	return Cell{
		Position: position,
		Size:     size,
		Color:    color,
	}
}

func (c *Cell) Draw() {
	rl.DrawRectangle(int32(c.Position.X), int32(c.Position.Y), int32(c.Size), int32(c.Size), c.Color)
}

type Figure struct {
	Cells         []Cell
	Position      rl.Vector2
	Color         rl.Color
	InactiveColor rl.Color
	IsActive      bool
}

func LocalToGlobalPosition(position rl.Vector2, localPosition rl.Vector2) rl.Vector2 {
	return rl.NewVector2(position.X+localPosition.X, position.Y+localPosition.Y)
}

func NewFigure(cells []Cell, position rl.Vector2, color rl.Color, inactiveColor rl.Color) Figure {
	for i := range cells {
		cells[i].Position = LocalToGlobalPosition(position, cells[i].Position)
	}

	return Figure{
		Cells:         cells,
		Position:      position,
		Color:         color,
		InactiveColor: inactiveColor,
		IsActive:      true,
	}
}

func (f *Figure) Move(position rl.Vector2) {
	f.Position = position
	for i := range f.Cells {
		f.Cells[i].Move(position)
	}
}

func (f *Figure) Draw() {
	for _, cell := range f.Cells {
		if cell.Color != f.Color {
			cell.Color = f.Color
		}

		cell.Draw()
	}
}

func (f *Figure) IsOutideGrid(grid *Grid) bool {
	for _, cell := range f.Cells {
		if !grid.IsInGrid(cell.Position) {
			return true
		}
	}

	return false
}

func (f *Figure) IsCollide(figure *Figure) bool {
	for _, cell := range figure.Cells {
		for _, cell2 := range f.Cells {
			if cell.Position == cell2.Position {
				return true
			}
		}
	}

	return false
}

func (f *Figure) Clone() Figure {
	figure := NewFigure([]Cell{}, f.Position, f.Color, f.InactiveColor)
	for _, cell := range f.Cells {
		figure.Cells = append(figure.Cells, cell)
	}

	return figure
}

func NewLShapehape(grid *Grid, position rl.Vector2, color rl.Color, inactiveColor rl.Color) Figure {
	gridSize := grid.Size

	return NewFigure([]Cell{
		NewCell(rl.NewVector2(0, 0), gridSize, color),
		NewCell(rl.NewVector2(0, 1), gridSize, color),
		NewCell(rl.NewVector2(0, 2), gridSize, color),
		NewCell(rl.NewVector2(1, 2), gridSize, color),
	}, position, color, inactiveColor)
}

func isFigureCollide(figure *Figure, figures []Figure) bool {
	for _, fig := range figures {
		if fig.IsCollide(figure) {
			return true
		}
	}

	return false
}

func main() {
	// rl.SetConfigFlags(rl.FlagWindowHighdpi)
	rl.InitWindow(1000, 1000, "raylib [core] example - basic window")
	defer rl.CloseWindow()
	camera := rl.Camera2D{}

	camera.Offset = rl.NewVector2(0, 0)
	camera.Target = rl.NewVector2(0, 0)
	camera.Zoom = 1

	grid := Grid{}
	grid.Size = 100

	var activeFigure *Figure
	figures := []Figure{}
	moveStep := 1000
	lastTime := time.Now()
	bottomRight := rl.NewVector2(float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight()))

	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		if activeFigure == nil {
			figure := NewLShapehape(&grid, grid.ToGridPosition(0, 0), rl.Red, rl.LightGray)
			activeFigure = &figure
			figures = append(figures, figure)
		}
		// if delay passed
		if time.Since(lastTime) > time.Duration(moveStep)*time.Millisecond {
			lastTime = time.Now()

			if activeFigure != nil {
				figure := activeFigure.Clone()
				figure.Move(grid.ToGridPosition(0, 1))
				isCollide := false

				// if figure + move outside of viewport
				for _, cell := range figure.Cells {
					if cell.Position.X < 0 || cell.Position.X >= bottomRight.X {
						isCollide = true
						break
					}

					if cell.Position.Y < 0 || cell.Position.Y >= bottomRight.Y {
						isCollide = true
						break
					}
				}

				if !isCollide {
					isCollide = isFigureCollide(&figure, figures[:len(figures)-1])
				}

				if isCollide {
					activeFigure.IsActive = false
					activeFigure = nil
				} else {
					activeFigure.Move(grid.ToGridPosition(0, 1))
				}
			}

		}

		if activeFigure != nil {
			var move rl.Vector2
			if rl.IsKeyPressed(rl.KeyLeft) {
				move = grid.ToGridPosition(-1, 0)
			}

			if rl.IsKeyPressed(rl.KeyRight) {
				move = grid.ToGridPosition(1, 0)
			}

			if rl.IsKeyPressed(rl.KeyDown) {
				move = grid.ToGridPosition(0, 1)
			}

			if (move != rl.Vector2{}) {
				fig := activeFigure.Clone()
				fig.Move(move)

				if !isFigureCollide(&fig, figures[:len(figures)-1]) {
					activeFigure.Move(move)
				}
			}
		}

		rl.BeginDrawing()

		rl.BeginMode2D(camera)

		rl.ClearBackground(rl.RayWhite)
		rl.DrawText("Congrats! You created your first window!", 500, 500, 20, rl.LightGray)

		grid.Draw()

		for _, figure := range figures {
			figure.Draw()
		}

		rl.EndMode2D()

		rl.EndDrawing()
	}
}
