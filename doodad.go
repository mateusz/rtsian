package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type doodad struct {
	position pixel.Vec
	sprites  *spriteset
	spriteID uint32
}

func DoodadInput(win *pixelgl.Window, cam pixel.Matrix) {
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		d := doodad{
			position: cam.Unproject(win.MousePosition().Scaled(1.0 / pixSize)),
			sprites:  &mobSprites,
			spriteID: 0,
		}
		mobs = append(mobs, &d)
	}
}

func (d *doodad) Update(dt float64) {
}

func (d *doodad) Draw(t pixel.Target) {
	d.sprites.sprites[d.spriteID].Draw(t, pixel.IM.Moved(d.position))
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
