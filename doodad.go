package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/mateusz/rtsian/piksele"
)

type doodad struct {
	position pixel.Vec
	sprites  *piksele.Spriteset
	spriteID uint32
}

func (d *doodad) Input(win *pixelgl.Window, cam pixel.Matrix) {

}

func (d *doodad) Update(dt float64) {

}

func (d *doodad) Draw(t pixel.Target) {
	d.sprites.Sprites[d.spriteID].Draw(t, pixel.IM.Moved(d.position))
}

func (d *doodad) GetZ() float64 {
	return d.position.Y
}

func (d *doodad) GetY() float64 {
	return d.position.Y
}

func (d *doodad) GetX() float64 {
	return d.position.X
}
