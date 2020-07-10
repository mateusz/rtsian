package main

import (
	"image/color"
	"math"

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

type unit struct {
	mobile
	sprite
	mouseTarget
	exploding
	selected bool
	army     int
	hp       float64
}

func NewUnit(position pixel.Vec, army int) unit {
	u := unit{
		mobile: mobile{
			position:  position,
			baseSpeed: 2.0,
		},
		sprite: sprite{
			spriteset: &mobSprites,
			spriteID:  MOBS_TANK_START_ID,
		},
		exploding: exploding{
			sprite: sprite{
				spriteset: &mobSprites,
				spriteID:  MOBS_EXPLOSION_START_ID,
			},
		},
		army: army,
	}
	return u
}

func UnitInput(win *pixelgl.Window, cam pixel.Matrix) {
	mp := cam.Unproject(win.MousePosition().Scaled(1.0 / pixSize))
	alignedMp := gameWorld.alignToTile(mp)

	if win.JustPressed(pixelgl.KeyW) {
		u := NewUnit(alignedMp, 1)
		u.target = u.position
		gameEntities.Add(&u)
		gamePositionables.Add(&u)
		gameDrawables.Add(&u)
	}
	if win.JustPressed(pixelgl.KeyE) {
		u := NewUnit(alignedMp, 2)
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
				if len(gameMouseHits) > 0 {
					// Exclude self-shots
					if !u.hitRight {
						m := NewMissile(u.position, alignedMp)
						gameEntities.Add(&m)
						gamePositionables.Add(&m)
						gameDrawables.Add(&m)
					}
				} else {
					u.pathingTarget = mp
					u.pathing = FindPath(u, u.pathingTarget)
				}
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
				u.pathing = nil
				u.startExploding()
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

			if !selectConsumed && u.mouseTarget.hitLeft {
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
	if u.totallyExploded() {
		gameEntities.Remove(u)
		gamePositionables.Remove(u)
		gameDrawables.Remove(u)
	}

	u.mobile.Update(dt)
}

func (u *unit) Draw(t pixel.Target) {
	explosionFrame := u.explosionFrame()
	if !u.exploding.exploding || explosionFrame < uint32(math.Ceil(float64(mobsExplosionFrames)/2.0)) {
		u.spriteset.sprites[u.spriteID+u.stickyDirOffset].DrawColorMask(
			t,
			rescueBottomPixels.Moved(u.position),
			colorMap[u.army],
		)
	}
	u.drawExplosion(t, u.position)
	if u.selected {
		cursorSprites.sprites[CURSOR_UNIT_MARKER].Draw(t, rescueBottomPixels.Moved(u.position))
	}
}
