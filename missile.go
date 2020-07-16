package main

import (
	"github.com/faiface/pixel"
	"github.com/mateusz/rtsian/piksele"
)

const (
	MOBS_MISSILE_START_ID = 31
)

type missile struct {
	mobile
	piksele.Sprite
}

func NewMissile(position pixel.Vec, target pixel.Vec) missile {
	mv := target.Sub(position)
	m := missile{
		mobile: mobile{
			position: position,
			target:   target,
			d:        mv.Len(),
			v:        mv.Unit().Scaled(10.0),
		},
		Sprite: piksele.Sprite{
			Spriteset: &mobSprites,
			SpriteID:  MOBS_MISSILE_START_ID,
		},
	}
	return m
}

func (m *missile) Update(dt float64) {
	if m.d <= 0.0 {
		for _, ent := range gameEntities.List {
			p, okp := ent.(positionable)
			e, oke := ent.(explodable)
			if okp && oke {
				if m.position.X > p.GetX()-8 &&
					m.position.X < p.GetX()+8 &&
					m.position.Y > p.GetY()-8 &&
					m.position.Y < p.GetY()+8 {
					e.startExploding()
				}
			}
		}

		gameEntities.Remove(m)
		gamePositionables.Remove(m)
		gameDrawables.Remove(m)
	}
	m.mobile.Update(dt)
}

func (m *missile) Draw(t pixel.Target) {
	m.Spriteset.Sprites[m.SpriteID].Draw(t, rescueBottomPixels.Moved(m.position))
}
