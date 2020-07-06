package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type unit struct {
	position pixel.Vec
	target   pixel.Vec
	d        float64
	v        pixel.Vec
	sprites  *spriteset
	spriteID uint32
	selected bool
}

func UnitInput(win *pixelgl.Window, cam pixel.Matrix) {
	if win.JustPressed(pixelgl.KeyQ) {
		u := unit{
			position: cam.Unproject(win.MousePosition().Scaled(1.0 / pixSize)),
			sprites:  &mobSprites,
			spriteID: 0,
		}
		u.target = u.position
		mobs = append(mobs, &u)
	}
	if win.JustPressed(pixelgl.MouseButtonRight) {
		for _, m := range mobs {
			unit, ok := m.(*unit)
			if !ok {
				continue
			}

			if unit.selected {
				unit.target = cam.Unproject(win.MousePosition().Scaled(1.0 / pixSize))
				unit.v = unit.target.Sub(unit.position).Scaled(0.01)
				unit.d = unit.target.Sub(unit.position).Len()
			}
		}
	}
}

func (u *unit) Input(win *pixelgl.Window, cam pixel.Matrix) {
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		mp := cam.Unproject(win.MousePosition().Scaled(1.0 / pixSize))

		if mp.X > u.position.X-u.sprites.sprites[u.spriteID].Frame().W()/2.0 &&
			mp.X < u.position.X+u.sprites.sprites[u.spriteID].Frame().W()/2.0 &&
			mp.Y > u.position.Y-u.sprites.sprites[u.spriteID].Frame().H()/2.0 &&
			mp.Y < u.position.Y+u.sprites.sprites[u.spriteID].Frame().H()/2.0 {
			// Hit
			u.selected = !u.selected
		}

	}
}

func (u *unit) Update(dt float64) {
	u.position = u.position.Add(u.v)
	u.d -= u.v.Len()
	if u.d < 0.0 {
		u.v = pixel.ZV
		u.d = 0.0
		u.target = pixel.ZV
	}
}

func (u *unit) Draw(t pixel.Target) {
	u.sprites.sprites[u.spriteID].Draw(t, pixel.IM.Moved(u.position))
	if u.selected {
		cursorSprites.sprites[0].Draw(t, pixel.IM.Moved(u.position))
	}
}

func (u *unit) GetZ() float64 {
	return u.position.Y
}

func (u *unit) GetY() float64 {
	return u.position.Y
}

func (u *unit) GetX() float64 {
	return u.position.X
}
