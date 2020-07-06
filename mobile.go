package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type mobile interface {
	Input(win *pixelgl.Window, cam pixel.Matrix)
	Update(dt float64)
	Draw(pixel.Target)
	GetZ() float64
	GetX() float64
	GetY() float64
}
