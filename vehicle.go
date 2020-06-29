package main

import (
	"image/color"
	_ "image/png"

	"github.com/faiface/pixel"
)

const (
	sprDirRight = 0
	sprDirUp    = 1
	sprDirDown  = 2
	sprDirLeft  = 3
)

type vehicle struct {
	// World position
	wp pixel.Vec
	// Velocity
	v         pixel.Vec
	spriteset *spriteset
	startID   uint32
	stickyDir uint32
	colorMask color.RGBA
}

func (v *vehicle) Draw(t pixel.Target) {
	v.dirToSpr(v.v.X, v.v.Y).DrawColorMask(t, pixel.IM.Moved(v.wp), v.colorMask)
}

func (v *vehicle) GetZ() float64 {
	return v.wp.Y
}

func (v *vehicle) GetY() float64 {
	return v.wp.Y
}

func (v *vehicle) GetX() float64 {
	return v.wp.X
}

func (v *vehicle) Update(dt float64) {
	v.wp = v.wp.Add(v.v.Scaled(dt))
}

func (v *vehicle) dirToSpr(dx, dy float64) *pixel.Sprite {
	if dx > 0 {
		v.stickyDir = sprDirRight
	}
	if dx < 0 {
		v.stickyDir = sprDirLeft
	}
	if dy > 0 {
		v.stickyDir = sprDirUp
	}
	if dy < 0 {
		v.stickyDir = sprDirDown
	}
	// ... and if 0,0, then use the old stickyDir so that the car doesn't randomly
	// flip after stopping!

	return v.spriteset.sprites[v.startID+v.stickyDir]
}
