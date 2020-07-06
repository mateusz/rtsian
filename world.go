package main

import (
	"fmt"
	"os"

	"github.com/faiface/pixel"
	"github.com/lafriks/go-tiled"
)

type world struct {
	tiles   *tiled.Map
	sprites spriteset
}

func (w *world) load() {
	var err error
	w.tiles, err = tiled.LoadFromFile(fmt.Sprintf("%s/../assets/world.tmx", workDir))
	if err != nil {
		fmt.Printf("Error parsing map: %s\n", err)
		os.Exit(2)
	}
	w.sprites, err = fillMissingMapPieces(w.tiles)
	if err != nil {
		fmt.Printf("Error loading aux tilesets: %s\n", err)
		os.Exit(2)
	}
}

func (w *world) pixelWidth() int {
	return w.tiles.Width * w.tiles.TileWidth
}

func (w *world) pixelHeight() int {
	return w.tiles.Height * w.tiles.TileHeight
}

func (w *world) Draw(t pixel.Target) {
	l := w.tiles.Layers[0]
	for y := 0; y < w.tiles.Height; y++ {
		for x := 0; x < w.tiles.Width; x++ {
			lt := l.Tiles[y*w.tiles.Width+x]
			w.sprites.sprites[lt.ID].Draw(t, pixel.IM.Moved(w.tileVec(x, w.tiles.Height-y-1)))
		}
	}
}

// Convert tile coords (x,y) to world coordinates.
func (w *world) tileVec(x int, y int) pixel.Vec {
	// Some offesting due to the tiles being referenced via the centre
	ox := w.tiles.TileWidth / 2
	oy := w.tiles.TileHeight / 2
	return pixel.V(float64(x*(w.tiles.TileWidth)+ox), float64(y*w.tiles.TileHeight+oy))
}