package main

import (
	"container/list"
	"log"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type unit struct {
	position      pixel.Vec
	target        pixel.Vec // position of current movement target
	d             float64   // distance to current movement target
	v             pixel.Vec // velocity vector
	pathingTarget pixel.Vec
	pathing       *list.List
	sprites       *spriteset
	spriteID      uint32
	selected      bool
}

func UnitInput(win *pixelgl.Window, cam pixel.Matrix) {
	mp := cam.Unproject(win.MousePosition().Scaled(1.0 / pixSize))

	if win.JustPressed(pixelgl.KeyW) {
		mpa := gameWorld.alignToTile(mp)
		u := unit{
			position: mpa,
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
				u.pathingTarget = mp
				u.pathing = FindPath(u, u.pathingTarget)
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
				if !win.Pressed(pixelgl.KeyLeftShift) && !win.Pressed(pixelgl.KeyRightShift) {
					u.selected = false
				}
			}
		}
	}
}

func (u *unit) Input(win *pixelgl.Window, cam pixel.Matrix) {
}

func (u *unit) Update(dt float64) {
	if u.d > 0.0 {
		u.position = u.position.Add(u.v)
		u.d -= u.v.Len()
	} else {
		u.applyPath()
	}
}

func (u *unit) applyPath() {
	if u.pathing == nil || u.pathing.Len() == 0 {
		u.target = u.position
		u.d = 0.0
		u.v = pixel.ZV
		u.pathing = nil
		return
	}

	u.pathing = FindPath(u, u.pathingTarget)
	if u.pathing.Len() == 0 {
		return
	}

	// Next path step
	n, ok := u.pathing.Remove(u.pathing.Front()).(*patherNode)
	if !ok {
		log.Panic("Fatal: pathing list contained non-patherNode!")
	}

	u.target = gameWorld.tileToVec(n.X, n.Y)
	mv := u.target.Sub(u.position)
	u.d = mv.Len()
	u.v = mv.Unit().Scaled(1.0)
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
