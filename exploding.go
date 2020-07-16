package main

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/mateusz/rtsian/piksele"
)

type exploding struct {
	piksele.Sprite
	exploding      bool
	explodingSince time.Time
}

func (e *exploding) startExploding() {
	if e.exploding {
		return
	}
	e.exploding = true
	e.explodingSince = time.Now()
}

func (e *exploding) explosionFrame() uint32 {
	return uint32(time.Now().Sub(e.explodingSince) / (100 * time.Millisecond))
}

func (e *exploding) totallyExploded() bool {
	if !e.exploding {
		return false
	}
	explosionFrame := uint32(time.Now().Sub(e.explodingSince) / (100 * time.Millisecond))
	return explosionFrame >= uint32(mobsExplosionFrames)
}

func (e *exploding) drawExplosion(t pixel.Target, p pixel.Vec) {
	if !e.exploding {
		return
	}

	explosionFrame := uint32(time.Now().Sub(e.explodingSince) / (100 * time.Millisecond))
	if explosionFrame < uint32(mobsExplosionFrames) {
		e.Spriteset.Sprites[e.SpriteID+explosionFrame].Draw(t, rescueBottomPixels.Moved(p))
	}
}
