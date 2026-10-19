package main

import (
	"fmt"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Drawable interface {
	Draw()
	Move(position rl.Vector2)
}

type Grid struct {
	CellSize int
	Cols     int
	Rows     int
	Position rl.Vector2
}

func (g *Grid) Draw() {
	for i := 0; i < g.Cols; i++ {
		for j := 0; j < g.Rows; j++ {
			rl.DrawRectangle(int32(i*g.CellSize), int32(j*g.CellSize), int32(g.CellSize), int32(g.CellSize), rl.LightGray)
			rl.DrawRectangleLines(int32(i*g.CellSize), int32(j*g.CellSize), int32(g.CellSize), int32(g.CellSize), rl.Black)
		}
	}
}

func (g *Grid) Move(position rl.Vector2) {
	_ = position
}

func (g *Grid) IsInGrid(position rl.Vector2) bool {
	x := int(position.X)
	y := int(position.Y)
	max := g.CellSize * g.CellSize

	if x < 0 || x >= max {
		return false
	}

	if y < 0 || y >= max {
		return false
	}

	return true
}

func (g *Grid) ToGridPosition(x, y int) rl.Vector2 {
	return rl.NewVector2(float32(x*g.CellSize), float32(y*g.CellSize))
}

func (g *Grid) ViewportBounds(x, y float32) rl.Vector2 {
	return rl.NewVector2(x*float32(g.CellSize), y*float32(g.CellSize))
}

type Cell struct {
	Position rl.Vector2
	Size     int
	Color    rl.Color
}

func (c *Cell) Move(position rl.Vector2) {
	c.Position = c.Position.Add(position)
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
	Cells    []Cell
	Position rl.Vector2
	Velocity rl.Vector2
	Color    rl.Color
	IsActive bool
}

func LocalToGlobalPosition(position rl.Vector2, localPosition rl.Vector2) rl.Vector2 {
	return rl.NewVector2(position.X+localPosition.X, position.Y+localPosition.Y)
}

func NewFigure(cells []Cell, position rl.Vector2, color rl.Color) Figure {
	for i := range cells {
		cells[i].Position = LocalToGlobalPosition(position, cells[i].Position)
	}

	return Figure{
		Cells:    cells,
		Position: position,
		Velocity: rl.NewVector2(0, 0),
		Color:    color,
		IsActive: true,
	}
}

func (f *Figure) Move(vector rl.Vector2) {
	f.Position = f.Position.Add(vector)
	for i := range f.Cells {
		f.Cells[i].Move(vector)
	}
}

func (f *Figure) Draw() {
	for i := range f.Cells {
		f.Cells[i].Color = f.Color
		f.Cells[i].Draw()
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
	figure := NewFigure([]Cell{}, f.Position, f.Color)
	for _, cell := range f.Cells {
		figure.Cells = append(figure.Cells, cell)
	}

	return figure
}

func NewLShape(grid *Grid, position rl.Vector2, color rl.Color) Figure {
	gridSize := grid.CellSize

	return NewFigure([]Cell{
		NewCell(rl.NewVector2(0, 0), gridSize, color),
		NewCell(rl.NewVector2(0, 1), gridSize, color),
		NewCell(rl.NewVector2(0, 2), gridSize, color),
		NewCell(rl.NewVector2(1, 2), gridSize, color),
	}, position, color)
}

func NewSquareShape(grid *Grid, position rl.Vector2, color rl.Color) Figure {
	gridSize := grid.CellSize

	return NewFigure([]Cell{
		NewCell(rl.NewVector2(0, 0), gridSize, color),
		NewCell(rl.NewVector2(1, 0), gridSize, color),
		NewCell(rl.NewVector2(1, 1), gridSize, color),
		NewCell(rl.NewVector2(0, 1), gridSize, color),
	}, position, color)
}

func NewLongShapehape(grid *Grid, position rl.Vector2, color rl.Color) Figure {
	gridSize := grid.CellSize

	return NewFigure([]Cell{
		NewCell(rl.NewVector2(0, 0), gridSize, color),
		NewCell(rl.NewVector2(0, 1), gridSize, color),
		NewCell(rl.NewVector2(0, 2), gridSize, color),
		NewCell(rl.NewVector2(0, 3), gridSize, color),
	}, position, color)
}

func isFigureCollide(figure *Figure, figures []Figure) bool {
	for _, fig := range figures {
		if fig.IsCollide(figure) {
			return true
		}
	}

	return false
}

const MoveTick = 500

type World struct {
	Grid           Grid
	Active         *Figure
	Figures        []Figure
	ActiveMoveTick float32
	BottomRight    rl.Vector2
}

func (w *World) PrepareFigure() {
	if w.Active != nil {
		return
	}

	figures := []Figure{
		NewLShape(&w.Grid, w.Grid.ToGridPosition(0, 0), rl.Red),
		NewSquareShape(&w.Grid, w.Grid.ToGridPosition(0, 0), rl.Blue),
		NewLongShapehape(&w.Grid, w.Grid.ToGridPosition(0, 0), rl.Green),
	}

	figure := figures[rand.Intn(len(figures))]
	w.Active = &figure
	w.Active.Velocity.Y = 100
	w.Figures = append(w.Figures, figure)
}

func (w *World) ProcessInput(dt float32) {
	if w.Active == nil {
		return
	}

	var move rl.Vector2
	if rl.IsKeyPressed(rl.KeyLeft) {
		move = w.Grid.ToGridPosition(-1, 0)
	}

	if rl.IsKeyPressed(rl.KeyRight) {
		move = w.Grid.ToGridPosition(1, 0)
	}

	if rl.IsKeyPressed(rl.KeyDown) {
		move = w.Grid.ToGridPosition(0, 1)
	}

	if (move != rl.Vector2{}) {
		fig := w.Active.Clone()
		fig.Move(move)

		if !isFigureCollide(&fig, w.Figures[:len(w.Figures)-1]) {
			w.Active.Move(move)
		}
	}

	if rl.IsKeyPressed(rl.KeyR) {
		w.Active.Rotate()
	}
}

func (f *Figure) Rotate() {
	if len(f.Cells) == 0 {
		return
	}

	// Rotate cells 90 degrees clockwise around the figure's origin
	for i := range f.Cells {
		local := rl.NewVector2(f.Cells[i].Position.X-f.Position.X, f.Cells[i].Position.Y-f.Position.Y)
		f.Cells[i].Position = rl.NewVector2(f.Position.X-local.Y, f.Position.Y+local.X)
	}
}

func (w *World) Update(dt float32) {
	w.ProcessInput(dt)
	w.PrepareFigure()
	w.ActiveMoveTick += dt * 1000

	fmt.Printf("ActiveMoveTick: %f, dt: %f\n", w.ActiveMoveTick, dt)

	if w.ActiveMoveTick >= MoveTick {
		w.ActiveMoveTick -= MoveTick
		if w.Active != nil {
			figure := w.Active.Clone()
			move := rl.NewVector2(0, float32(w.Grid.CellSize))
			figure.Move(move)
			isCollide := false

			// if figure + move outside of viewport
			for _, cell := range figure.Cells {
				if cell.Position.X < 0 || cell.Position.X >= w.BottomRight.X {
					isCollide = true
					break
				}

				if cell.Position.Y < 0 || cell.Position.Y >= w.BottomRight.Y {
					isCollide = true
					break
				}
			}

			if !isCollide {
				isCollide = isFigureCollide(&figure, w.Figures[:len(w.Figures)-1])
			}

			if isCollide {
				w.Active.IsActive = false
				w.Active = nil
			} else {
				w.Active.Move(move)
			}
		}
	}
	w.DetectCombosCollisions()
}

func (w *World) Render(camera rl.Camera2D) {
	rl.BeginDrawing()
	rl.BeginMode2D(camera)
	rl.ClearBackground(rl.RayWhite)
	rl.DrawText("Congrats! You created your first window!", 500, 500, 20, rl.LightGray)
	w.Grid.Draw()
	for _, figure := range w.Figures {
		figure.Draw()
	}
	rl.EndMode2D()
	rl.EndDrawing()
}

func (w *World) DetectCombosCollisions() {
	// for i := 0; i < len(w.Figures); i++ {
	// 	for j := i + 1; j < len(w.Figures); j++ {
	// 		if w.Figures[i].IsCollide(&w.Figures[j]) {
	// 			w.Figures[i].Color = w.Figures[i].InactiveColor
	// 			w.Figures[j].Color = w.Figures[j].InactiveColor
	// 		}
	// 	}
	// }
}

func main() {
	// rl.SetConfigFlags(rl.FlagWindowHighdpi)
	runApp()
}
