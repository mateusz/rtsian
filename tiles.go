package main

import (
	"encoding/xml"
	"errors"
	"image"
	"log"
	"os"
	"path/filepath"
	"strconv"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/lafriks/go-tiled"
)

type spriteset struct {
	sprites  map[uint32]*pixel.Sprite
	tileset  *tiled.Tileset
	basePath string
}

func newSpriteset() spriteset {
	t := spriteset{}
	t.sprites = make(map[uint32]*pixel.Sprite)
	return t
}

// TMX library does not load images. Help it out.
func fillMissingMapPieces(m *tiled.Map) (spriteset, error) {
	spr := spriteset{}
	var err error
	for _, ts := range m.Tilesets {
		if ts.Source == "" {
			return spr, errors.New("Tileset has no source")
		}
		spr, err = newSpritesetFromTileset(m.GetFileFullPath(""), ts)
		if err != nil {
			return spr, err
		}

		// Only one permitted at the moment.
		break
	}

	return spr, nil
}

func loadObjects(m *tiled.Map) {
	for _, o := range m.ObjectGroups[0].Objects {
		lt, err := m.TileGIDToTile(o.GID)
		if err != nil {
			log.Fatal(err)
		}
		army, err := strconv.Atoi(o.Properties.GetString("army"))
		if err != nil {
			army = 0
		}
		if lt.ID >= MOBS_TANK_START_ID && lt.ID < MOBS_TANK_START_ID+4 {
			p := gameWorld.alignToTile(pixel.Vec{X: o.X + 10.0, Y: tiledFlipY(m, o.Y) + 10.0})
			u := NewUnit(p, army)
			u.target = u.position
			gameEntities.Add(&u)
			gamePositionables.Add(&u)
			gameDrawables.Add(&u)
		}
	}
}

func newSpritesetFromTsx(basePath, path string) (spriteset, error) {
	spr := spriteset{}
	ts := &tiled.Tileset{Source: path}

	f, err := os.Open(filepath.Join(basePath, ts.Source))
	if err != nil {
		return spr, err
	}
	defer f.Close()

	d := xml.NewDecoder(f)
	if err := d.Decode(ts); err != nil {
		return spr, err
	}

	spr, err = newSpritesetFromTileset(basePath, ts)
	if err != nil {
		return spr, err
	}

	spr.tileset = ts
	return spr, nil
}

// Load actual sprite files and associate with tileset.
func newSpritesetFromTileset(basePath string, ts *tiled.Tileset) (spriteset, error) {
	spr := newSpriteset()
	spr.tileset = ts
	spr.basePath = basePath

	f, err := os.Open(filepath.Join(basePath, ts.Source))

	if err != nil {
		return spr, err
	}
	defer f.Close()

	d := xml.NewDecoder(f)

	if err := d.Decode(ts); err != nil {
		return spr, err
	}

	for _, t := range ts.Tiles {
		if t.Image.Source == "" {
			continue
		}

		file, err := os.Open(filepath.Join(basePath, t.Image.Source))
		if err != nil {
			return spr, err
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			return spr, err
		}

		pic := pixel.PictureDataFromImage(img)
		spr.sprites[t.ID] = pixel.NewSprite(pic, pic.Bounds())
	}

	return spr, nil
}

func tiledFlipY(t *tiled.Map, y float64) float64 {
	return float64(t.Height*t.TileHeight) - y
}
