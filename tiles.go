package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"image"
	"os"
	"path/filepath"

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

// Get sprite from file.
func load(path string) (*pixel.Sprite, error) {
	file, err := os.Open(tmx.GetFileFullPath(path))
	if err != nil {
		return nil, fmt.Errorf("error opening car: %s", err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("error decoding car: %s", err)
	}
	pic := pixel.PictureDataFromImage(img)
	return pixel.NewSprite(pic, pic.Bounds()), nil
}

func findTileInTileset(lt *tiled.LayerTile) (*tiled.TilesetTile, error) {
	for _, t := range lt.Tileset.Tiles {
		if t.ID == lt.ID {
			return t, nil
		}
	}

	return nil, fmt.Errorf("Something is very wrong, tile ID '%d' not found in the tileset", lt.ID)
}

// Convert tile coords (x,y) to world coordinates.
func tileVec(x int, y int) pixel.Vec {
	// Some offesting due to the tiles being referenced via the centre
	ox := tmx.TileWidth / 2
	oy := tmx.TileHeight / 2
	return pixel.V(float64(x*(tmx.TileWidth)+ox), float64(y*tmx.TileHeight+oy))
}
