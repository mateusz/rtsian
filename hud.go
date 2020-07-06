package main

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

type hud struct {
	bounds   pixel.Rect
	mobCount int
}

func (h *hud) Draw(t pixel.Target) {
	mobCount := text.New(pixel.ZV, text.NewAtlas(basicfont.Face7x13, text.ASCII))
	fmt.Fprintf(mobCount, "%d", h.mobCount)

	mobCount.Draw(t, pixel.IM.Moved(pixel.Vec{
		X: 5.0,
		Y: h.bounds.H() - 15.0,
	}))
}

func (h *hud) Update(dt float64) {
	h.mobCount = len(mobs)
}
