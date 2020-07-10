package main

import (
	"image/color"
	"log"
	"math"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	MOBS_TANK_START_ID      = 20
	MOBS_EXPLOSION_START_ID = 24
	MOBS_EXPLOSION_END_ID   = 30
	CURSOR_UNIT_MARKER      = 0
)

var (
	mobsExplosionFrames = MOBS_EXPLOSION_END_ID - MOBS_EXPLOSION_START_ID + 1
	rescueBottomPixels  = pixel.IM.Scaled(pixel.Vec{X: 8.0, Y: 8.0}, 1.01)
	colorMap            = map[int]color.Color{
		0: colornames.Pink,
		1: colornames.Green,
		2: colornames.Red,
	}
)

type missile struct {
}

type unit struct {
	mobile
	sprites        *spriteset
	spriteID       uint32
	selected       bool
	exploding      bool
	explodingSince time.Time
	missile        *missile
	missileSince   time.Time
	army           int
	hp             float64
}

func NewUnit(position pixel.Vec, army int) unit {
	u := unit{
		mobile: mobile{
			position:  position,
			baseSpeed: 2.0,
		},
		sprites:  &mobSprites,
		spriteID: MOBS_TANK_START_ID,
		army:     army,
		hp:       100.0,
	}
	return u
}

func UnitInput(win *pixelgl.Window, cam pixel.Matrix) {
	mp := cam.Unproject(win.MousePosition().Scaled(1.0 / pixSize))

	if win.JustPressed(pixelgl.KeyW) {
		mpa := gameWorld.alignToTile(mp)
		u := NewUnit(mpa, 1)
		u.target = u.position
		gameEntities.Add(&u)
		gamePositionables.Add(&u)
		gameDrawables.Add(&u)
	}

	if win.JustPressed(pixelgl.MouseButtonRight) {
		for _, ent := range gameEntities.List {
			u, ok := ent.(*unit)
			if !ok {
				continue
			}

			if u.selected {
				u.pathingTarget = mp
				u.pathing = FindPath(u, u.pathingTarget)
			}
		}
	}

	if win.JustPressed(pixelgl.KeyQ) {
		for _, ent := range gameEntities.List {
			u, ok := ent.(*unit)
			if !ok {
				continue
			}

			if u.selected {
				u.selected = false
				u.exploding = true
				u.explodingSince = time.Now()
				u.pathing = nil
			}
		}
	}

	if win.JustPressed(pixelgl.MouseButtonLeft) {
		selectConsumed := false
		for _, m := range gamePositionables.ByReverseZ() {
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
	if u.exploding {
		explosionFrame := uint32(time.Now().Sub(u.explodingSince) / (100 * time.Millisecond))
		if explosionFrame >= uint32(mobsExplosionFrames) {
			// Totally exploded.
			gameEntities.Remove(u)
			gamePositionables.Remove(u)
			gameDrawables.Remove(u)
		}
	}

	u.mobile.Update(dt)
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
	u.v = mv.Unit().Scaled(u.baseSpeed / n.Cost)
}

func (u *unit) updateDirOffset() {
	if math.Abs(u.v.X) > math.Abs(u.v.Y) {
		if u.v.X < 0 {
			u.stickyDirOffset = 0
		}
		if u.v.X > 0 {
			u.stickyDirOffset = 2
		}
	} else {
		if u.v.Y > 0 {
			u.stickyDirOffset = 1
		}
		if u.v.Y < 0 {
			u.stickyDirOffset = 3
		}
	}
}

func (u *unit) Draw(t pixel.Target) {
	explosionFrame := uint32(time.Now().Sub(u.explodingSince) / (100 * time.Millisecond))
	if !u.exploding || explosionFrame < uint32(math.Ceil(float64(mobsExplosionFrames)/2.0)) {
		u.sprites.sprites[u.spriteID+u.stickyDirOffset].DrawColorMask(
			t,
			rescueBottomPixels.Moved(u.position),
			colorMap[u.army],
		)
	}
	if u.exploding {
		if explosionFrame < uint32(mobsExplosionFrames) {
			u.sprites.sprites[MOBS_EXPLOSION_START_ID+explosionFrame].Draw(t, rescueBottomPixels.Moved(u.position))
		}
	}
	if u.selected {
		cursorSprites.sprites[CURSOR_UNIT_MARKER].Draw(t, rescueBottomPixels.Moved(u.position))
	}
}
