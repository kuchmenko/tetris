package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Renderer struct {
	width  int
	height int
	camera *rl.Camera2D

	colors map[int]rl.Color
}

func NewRenderer(width int, height int, camera *rl.Camera2D) *Renderer {
	return &Renderer{
		width:  width,
		height: height,
		camera: camera,
		colors: make(map[int]rl.Color),
	}
}

var TetrisColors = []rl.Color{
	rl.Red,
	rl.Green,
	rl.Blue,
	rl.Yellow,
	rl.Purple,
	rl.Orange,
	rl.Beige,
}
