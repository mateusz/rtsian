package main

import(
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	scrollSpeed float64
	scrollHotZone float64
)

type player struct {
	wp pixel.Vec
}

func (p *player) Update(dt float64, win *pixelgl.Window) {
	if win.MouseInsideWindow() {
		if win.MousePosition().X < scrollHotZone {
			p.wp.X -= scrollSpeed * dt
		}
		if win.MousePosition().X > monW-scrollHotZone {
			p.wp.X += scrollSpeed * dt
		}
		if win.MousePosition().Y < scrollHotZone {
			p.wp.Y -= scrollSpeed * dt
		}
		if win.MousePosition().Y > monH-scrollHotZone {
			p.wp.Y += scrollSpeed * dt
		}
	}

}