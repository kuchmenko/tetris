package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
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
	for i, cell := range cells {
		rotated[i] = TetrisCell{
			Position: Position{
				// FIXME
			},
		}
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
	currentFigure   *TetrisFigure
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

	if t.currentFigure == nil {
		t.generateFigure()
		return
	}

	err := t.Move(Down)
	var collisionErr *ErrMoveFailed

	if errors.As(err, &collisionErr) {
		log.Printf("Collision: %v", collisionErr)
		if collisionErr.Collision == nil {
			return
		}

		if collisionErr.Direction == Down && collisionErr.Collision.A.Position.Y == 1 {
			t.isFinished = true
			t.currentFigure = nil
			return
		}

		if collisionErr.Collision.B.Equal(TetrisFigure{}) || collisionErr.Direction == Down {
			t.currentFigure = nil

			return
		}
	}
}

func (t *Tetris) generateFigure() {
	if t.currentFigure != nil {
		return
	}

	t.lastFigureId++
	figure := t.figuresCreators[rand.Intn(len(t.figuresCreators))](t.lastFigureId)

	t.figures = append(t.figures, &figure)
	t.currentFigure = &figure
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

func (t *Tetris) Move(direction MoveDirection) error {
	if t.currentFigure == nil {
		return nil
	}
	log.Printf("Move: %v, currentFigure: %v", direction, t.currentFigure)

	candidate := *t.currentFigure

	switch direction {
	case Left:
		candidate.Position.X--
	case Right:
		candidate.Position.X++
	case Down:
		candidate.Position.Y++
	}

	err := t.updateCurrent(candidate)
	if err != nil {
		return &ErrMoveFailed{Collision: t.DetectColission([]TetrisFigure{candidate}), Direction: direction}
	}

	return nil
}

func (t *Tetris) updateCurrent(candidate TetrisFigure) error {
	collision := t.DetectColission([]TetrisFigure{candidate})
	if collision != nil {
		log.Printf("collision: %v", collision)
		return &ErrCollision{collision}
	}

	log.Printf("Saving currentFigure: %v", candidate)
	t.currentFigure.Position = candidate.Position
	t.currentFigure.Cells = candidate.Cells

	return nil
}

func (t *Tetris) Rotate() error {
	if t.currentFigure == nil {
		return nil
	}
	log.Printf("Rotate: currentFigure: %v", t.currentFigure)

	candidate := *t.currentFigure

	candidate.Cells = rotate90Cells(candidate.Cells, candidate.Position, t.cols, t.rows)

	err := t.updateCurrent(candidate)
	if err != nil {
		return &ErrRotateFailed{Collision: t.DetectColission([]TetrisFigure{candidate})}
	}

	return t.updateCurrent(candidate)
}
