package main

import (
	"github.com/faiface/pixel"
)

type mobile interface {
	Update(dt float64)
	Draw(pixel.Target)
	GetZ() float64
	GetX() float64
	GetY() float64
	Colliding(mobile, float64) bool
}
