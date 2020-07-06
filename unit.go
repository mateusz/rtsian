package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type unit struct {
	position pixel.Vec
	target   pixel.Vec // position of target
	d        float64   // distance to target
	v        pixel.Vec // velocity vector
	sprites  *spriteset
	spriteID uint32
	selected bool
}

func UnitInput(win *pixelgl.Window, cam pixel.Matrix) {
	mp := cam.Unproject(win.MousePosition().Scaled(1.0 / pixSize))

	if win.JustPressed(pixelgl.KeyQ) {
		u := unit{
			position: cam.Unproject(win.MousePosition().Scaled(1.0 / pixSize)),
			sprites:  &mobSprites,
			spriteID: 0,
		}
		u.target = u.position
		gameMobiles.Add(&u)
	}

	if win.JustPressed(pixelgl.MouseButtonRight) {
		for _, m := range gameMobiles.List {
			u, ok := m.(*unit)
			if !ok {
				continue
			}

			if u.selected {
				u.target = mp
				mv := u.target.Sub(u.position)
				u.v = mv.Unit().Scaled(1.0)
				u.d = mv.Len()
			}
		}
	}

	if win.JustPressed(pixelgl.MouseButtonLeft) {
		selectConsumed := false
		for _, m := range gameMobiles.ByReverseZ() {
			u, ok := m.(*unit)
			if !ok {
				continue
			}
			if !selectConsumed &&
				mp.X > u.position.X-u.sprites.sprites[u.spriteID].Frame().W()/2.0 &&
				mp.X < u.position.X+u.sprites.sprites[u.spriteID].Frame().W()/2.0 &&
				mp.Y > u.position.Y-u.sprites.sprites[u.spriteID].Frame().H()/2.0 &&
				mp.Y < u.position.Y+u.sprites.sprites[u.spriteID].Frame().H()/2.0 {
				// Hit
				u.selected = !u.selected
				selectConsumed = true
			} else {
				u.selected = false
			}
		}
	}
}

func (u *unit) Input(win *pixelgl.Window, cam pixel.Matrix) {
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
	return -u.position.Y
}

func (u *unit) GetY() float64 {
	return u.position.Y
}

func (u *unit) GetX() float64 {
	return u.position.X
}
