package main

import "github.com/faiface/pixel/pixelgl"

type mouseTarget struct {
	hitLeft  bool
	hitRight bool
}

func (mh *mouseTarget) MouseClear() {
	mh.hitLeft = false
	mh.hitRight = false
}

func (mh *mouseTarget) MouseHit(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		mh.hitLeft = true
	}
	if win.JustPressed(pixelgl.MouseButtonRight) {
		mh.hitRight = true
	}
}
