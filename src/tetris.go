package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"slices"
)

type User struct {
	ID   int
	Name string
}

type Position struct {
	X, Y int
}

type Collision struct {
	A, B     TetrisFigure
	Position Position
}

type TetrisCell struct {
	Position Position
}

type MoveDirection int

const (
	Left MoveDirection = iota
	Right
	Down
)

type TetrisFigure struct {
	ID       int
	Position Position
	Cells    []TetrisCell
}

func (t *TetrisFigure) Equal(other TetrisFigure) bool {
	if t.ID != other.ID {
		return false
	}
	//
	// if len(t.Cells) != len(other.Cells) {
	// 	return false
	// }
	//
	// for i := range t.Cells {
	// 	if t.Cells[i].Position != other.Cells[i].Position {
	// 		return false
	// 	}
	// }

	return true
}

func NewTetrisFigureI(id int) TetrisFigure {
	return TetrisFigure{
		ID:       id,
		Position: Position{X: 0, Y: 0},
		Cells: []TetrisCell{
			{Position: Position{X: 0, Y: 0}},
			{Position: Position{X: 0, Y: 1}},
			{Position: Position{X: 0, Y: 2}},
			{Position: Position{X: 0, Y: 3}},
		},
	}
}

func NewTetrisFigureJ(id int) TetrisFigure {
	return TetrisFigure{
		ID:       id,
		Position: Position{X: 0, Y: 0},
		Cells: []TetrisCell{
			{Position: Position{X: 1, Y: 0}},
			{Position: Position{X: 1, Y: 1}},
			{Position: Position{X: 1, Y: 2}},
			{Position: Position{X: 0, Y: 2}},
		},
	}
}

func NewTetrisFigureL(id int) TetrisFigure {
	return TetrisFigure{
		ID:       id,
		Position: Position{X: 0, Y: 0},
		Cells: []TetrisCell{
			{Position: Position{X: 0, Y: 0}},
			{Position: Position{X: 0, Y: 1}},
			{Position: Position{X: 0, Y: 2}},
			{Position: Position{X: 1, Y: 2}},
		},
	}
}

func NewTetrisFigureO(id int) TetrisFigure {
	return TetrisFigure{
		ID:       id,
		Position: Position{X: 0, Y: 0},
		Cells: []TetrisCell{
			{Position: Position{X: 0, Y: 0}},
			{Position: Position{X: 0, Y: 1}},
			{Position: Position{X: 1, Y: 0}},
			{Position: Position{X: 1, Y: 1}},
		},
	}
}

func NewTetrisFigureS(id int) TetrisFigure {
	return TetrisFigure{
		ID:       id,
		Position: Position{X: 0, Y: 0},
		Cells: []TetrisCell{
			{Position: Position{X: 0, Y: 1}},
			{Position: Position{X: 1, Y: 1}},
			{Position: Position{X: 1, Y: 0}},
			{Position: Position{X: 2, Y: 0}},
		},
	}
}

func NewTetrisFigureT(id int) TetrisFigure {
	return TetrisFigure{
		ID:       id,
		Position: Position{X: 0, Y: 0},
		Cells: []TetrisCell{
			{Position: Position{X: 0, Y: 0}},
			{Position: Position{X: 1, Y: 0}},
			{Position: Position{X: 2, Y: 0}},
			{Position: Position{X: 1, Y: 1}},
		},
	}
}

func NewTetrisFigureZ(id int) TetrisFigure {
	return TetrisFigure{
		ID:       id,
		Position: Position{X: 0, Y: 0},
		Cells: []TetrisCell{
			{Position: Position{X: 0, Y: 0}},
			{Position: Position{X: 1, Y: 0}},
			{Position: Position{X: 1, Y: 1}},
			{Position: Position{X: 2, Y: 1}},
		},
	}
}

// Rotate cells in 90 degrees and keep them in place relative to position
func rotate90Cells(cells []TetrisCell, position Position, cols, rows int) []TetrisCell {
	rotated := make([]TetrisCell, len(cells))

	slices.Reverse(cells)

	for id, cell := range cells {
		rotated[id] = TetrisCell{Position: Position{X: -cell.Position.Y, Y: cell.Position.X}}
	}

	return rotated
}

type ErrCollision struct {
	Collision *Collision
}

func (e *ErrCollision) Error() string {
	return fmt.Sprintf("collision: %v", e.Collision)
}

type ErrMoveFailed struct {
	Collision *Collision
	Direction MoveDirection
}

func (e *ErrMoveFailed) Error() string {
	var direction string
	switch e.Direction {
	case Left:
		direction = "Left"
	case Right:
		direction = "Right"
	case Down:
		direction = "Down"
	}
	return fmt.Sprintf("move failed: %v, %s", e.Collision, direction)
}

type ErrRotateFailed struct {
	Collision *Collision
}

func (e *ErrRotateFailed) Error() string {
	return fmt.Sprintf("rotate failed: %v", e.Collision)
}

type Tetris struct {
	figuresCreators []func(id int) TetrisFigure
	figures         []*TetrisFigure
	active          *TetrisFigure
	lastFigureId    int
	user            User
	score           int
	cols            int
	rows            int
	isFinished      bool
}

func NewTetris(cols, rows int) *Tetris {
	return &Tetris{
		cols:            cols,
		rows:            rows,
		figuresCreators: []func(id int) TetrisFigure{NewTetrisFigureI, NewTetrisFigureJ, NewTetrisFigureL, NewTetrisFigureO, NewTetrisFigureS, NewTetrisFigureT, NewTetrisFigureZ},
	}
}

func (t *Tetris) Tick() {
	if t.isFinished {
		return
	}

	t.detectCombos()

	if t.active == nil {
		t.generateFigure()
		return
	}

	err := t.MoveActive(Down)
	var collisionErr *ErrMoveFailed

	if errors.As(err, &collisionErr) {
		log.Printf("Collision: %v", collisionErr)
		if collisionErr.Collision == nil {
			return
		}

		if collisionErr.Direction == Down && collisionErr.Collision.A.Position.Y == 1 {
			t.isFinished = true
			t.active = nil
			return
		}

		if collisionErr.Collision.B.Equal(TetrisFigure{}) || collisionErr.Direction == Down {
			t.active = nil

			return
		}
	}
}

func (t *Tetris) generateFigure() {
	if t.active != nil {
		return
	}

	t.lastFigureId++
	figure := t.figuresCreators[rand.Intn(len(t.figuresCreators))](t.lastFigureId)

	t.figures = append(t.figures, &figure)
	t.active = &figure
}

func (t *Tetris) DetectColission(updatedFigures []TetrisFigure) *Collision {
	for _, figure := range updatedFigures {
		for _, cell := range figure.Cells {
			cellX := figure.Position.X + cell.Position.X
			cellY := figure.Position.Y + cell.Position.Y
			if cellX < 0 || cellX >= t.cols || cellY < 0 || cellY >= t.rows {
				return &Collision{figure, TetrisFigure{}, cell.Position}
			}
		}
	}

	for _, updatedFigure := range updatedFigures {
		for _, figure := range t.figures {
			if figure.Equal(updatedFigure) {
				continue
			}

			for _, cell := range updatedFigure.Cells {
				updatedX := cell.Position.X + updatedFigure.Position.X
				updatedY := cell.Position.Y + updatedFigure.Position.Y

				for _, c := range figure.Cells {
					x := c.Position.X + figure.Position.X
					y := c.Position.Y + figure.Position.Y

					if x == updatedX && y == updatedY {
						return &Collision{updatedFigure, *figure, cell.Position}
					}
				}
			}
		}
	}

	return nil
}

func (t *Tetris) moveFigure(candidate *TetrisFigure, direction MoveDirection) error {
	switch direction {
	case Left:
		candidate.Position.X--
	case Right:
		candidate.Position.X++
	case Down:
		candidate.Position.Y++
	}

	collision := t.DetectColission([]TetrisFigure{*candidate})
	if collision != nil {
		return &ErrMoveFailed{Collision: collision, Direction: direction}
	}

	return nil
}

func (t *Tetris) MoveActive(direction MoveDirection) error {
	if t.active == nil {
		return nil
	}
	log.Printf("Move: %v, currentFigure: %v", direction, t.active)

	candidate := *t.active

	collision := t.moveFigure(&candidate, direction)
	if collision != nil {
		return collision
	}

	t.updateCurrent(candidate)
	return nil
}

func (t *Tetris) updateCurrent(candidate TetrisFigure) {
	t.active.Position = candidate.Position
	t.active.Cells = candidate.Cells
}

func (t *Tetris) Rotate() error {
	if t.active == nil {
		return nil
	}
	log.Printf("Rotate: currentFigure: %v", t.active)

	candidate := *t.active

	candidate.Cells = rotate90Cells(candidate.Cells, candidate.Position, t.cols, t.rows)

	collision := t.DetectColission([]TetrisFigure{candidate})

	if collision != nil {
		return &ErrRotateFailed{Collision: collision}
	}

	t.updateCurrent(candidate)

	return nil
}

func (t *Tetris) Restart() {
	t.isFinished = false
	t.score = 0

	t.figures = nil
	t.active = nil
}

type lineCombo struct {
	cellsCount int
	figures    map[int]*TetrisFigure
}

func (t *Tetris) detectCombos() {
	rows := make(map[int]lineCombo)

	for _, figure := range t.figures {
		for _, cell := range figure.Cells {
			position := cell.Position.Y + figure.Position.Y
			combo := rows[position]

			if combo.figures == nil {
				combo.figures = make(map[int]*TetrisFigure)
			}
			combo.cellsCount++
			combo.figures[figure.ID] = figure
			rows[position] = combo
		}
	}
	for row, combo := range rows {
		if combo.cellsCount == t.cols {
			for _, figure := range combo.figures {
				cells := make([]TetrisCell, 0, len(figure.Cells))
				for _, cell := range figure.Cells {
					position := cell.Position.Y + figure.Position.Y

					if position == row {
						t.score += 10
						continue
					}

					cells = append(cells, cell)
				}
				cellHasChanged := len(cells) != len(figure.Cells)
				updatedFigure := *figure
				updatedFigure.Cells = cells

				if cellHasChanged {
					err := t.moveFigure(&updatedFigure, Down)
					if err == nil {
						figure.Position = updatedFigure.Position
						figure.Cells = updatedFigure.Cells
					}

					for id := range t.figures {
						if t.figures[id].Equal(*figure) {
							t.figures[id] = figure
						}
					}
				}

				if len(cells) == 0 {
					t.removeFigure(figure)
				}
			}
		}
	}
}

func (t *Tetris) removeFigure(figure *TetrisFigure) {
	updated := make([]*TetrisFigure, 0, len(t.figures))

	for _, f := range t.figures {
		if f.Equal(*figure) {
			if t.active != nil && t.active.Equal(*figure) {
				t.active = nil
			}
			continue
		}

		updated = append(updated, f)
	}

	t.figures = updated
}

func (t *Tetris) AddFigureToBottom(figure TetrisFigure) {
	figureHeight := 0
	for _, cell := range figure.Cells {
		if cell.Position.Y > figureHeight {
			figureHeight = cell.Position.Y
		}
	}
	figure.Position.Y = t.rows - figureHeight - 1

	t.figures = append(t.figures, &figure)

	if figure.ID > t.lastFigureId {
		t.lastFigureId = figure.ID
	}
}
