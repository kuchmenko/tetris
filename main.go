package main

import (
	"math/rand"
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

func (f *Figure) Move(dt float32) {
	f.Position = f.Position.Add(rl.NewVector2(f.Velocity.X*dt, f.Velocity.Y*dt))
	position := rl.NewVector2(f.Velocity.X*dt, f.Velocity.Y*dt)
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
	figure := NewFigure([]Cell{}, f.Position, f.Color)
	for _, cell := range f.Cells {
		figure.Cells = append(figure.Cells, cell)
	}

	return figure
}

func NewLShape(grid *Grid, position rl.Vector2, color rl.Color) Figure {
	gridSize := grid.Size

	return NewFigure([]Cell{
		NewCell(rl.NewVector2(0, 0), gridSize, color),
		NewCell(rl.NewVector2(0, 1), gridSize, color),
		NewCell(rl.NewVector2(0, 2), gridSize, color),
		NewCell(rl.NewVector2(1, 2), gridSize, color),
	}, position, color)
}

func NewSquareShape(grid *Grid, position rl.Vector2, color rl.Color) Figure {
	gridSize := grid.Size

	return NewFigure([]Cell{
		NewCell(rl.NewVector2(0, 0), gridSize, color),
		NewCell(rl.NewVector2(1, 0), gridSize, color),
		NewCell(rl.NewVector2(1, 1), gridSize, color),
		NewCell(rl.NewVector2(0, 1), gridSize, color),
	}, position, color)
}

func NewLongShapehape(grid *Grid, position rl.Vector2, color rl.Color) Figure {
	gridSize := grid.Size

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

type World struct {
	Grid        Grid
	Active      *Figure
	Figures     []Figure
	MoveStep    int
	LastTime    time.Time
	BottomRight rl.Vector2
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
		fig.Move(dt)

		if !isFigureCollide(&fig, w.Figures[:len(w.Figures)-1]) {
			w.Active.Move(dt)
		}
	}
}

func (w *World) Update(dt float32) {
	w.ProcessInput(dt)
	w.PrepareFigure()

	// for _, figure := range w.Figures {
	// 	figure.Move(dt)
	// }

	// if time.Since(w.LastTime) > time.Duration(w.MoveStep)*time.Millisecond {
	// 	w.LastTime = time.Now()
	// 	if w.Active != nil {
	// 		figure := w.Active.Clone()
	// 		figure.Move(dt)
	// 	}

	if w.Active != nil {
		isCollide := false

		// if figure + move outside of viewport
		for _, cell := range w.Active.Cells {
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
			isCollide = isFigureCollide(w.Active, w.Figures[:len(w.Figures)-1])
		}

		if isCollide {
			w.Active.IsActive = false
			w.Active = nil
		} else {
			w.Active.Move(dt)
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
	rl.InitWindow(1000, 1000, "raylib [core] example - basic window")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

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
	world := World{
		Grid:        grid,
		Active:      activeFigure,
		Figures:     figures,
		MoveStep:    moveStep,
		LastTime:    lastTime,
		BottomRight: bottomRight,
	}
	const tick = 1.0 / 60.0
	var accumulator float32

	for !rl.WindowShouldClose() {
		accumulator += rl.GetFrameTime()
		for accumulator >= tick {
			world.Update(tick)
			accumulator -= tick
		}

		world.Render(camera)

	}
}
